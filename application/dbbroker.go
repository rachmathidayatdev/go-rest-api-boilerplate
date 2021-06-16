package application

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-rest-api-boilerplate/util"
	_ "github.com/lib/pq"
)

type (
	//DBConfig struct
	DBConfig struct {
		Dialect  string
		Username string
		Password string
		Database string
		Host     string
		Port     string
		Charset  string
	}

	//DBBroker struct
	DBBroker struct {
		sync.Mutex
		Conn      *sql.DB
		watchStop chan bool
		errors    chan error
	}
)

func getDBConnectionString() *DBConfig {
	return &DBConfig{
		Dialect:  "postgres",
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASS"),
		Database: os.Getenv("DB_NAME"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Charset:  "utf8",
	}
}

//NewDBBroker func
func NewDBBroker(dbConfig *DBConfig) *DBBroker {
	broker := DBBroker{
		sync.Mutex{},
		nil,
		make(chan bool),
		nil,
	}

	return &broker
}

//Start func
func (b *DBBroker) Start() (*sql.DB, error) {
	conn, err := b.setup()
	go b.watch()
	return conn, err
}

//Stop func
func (b *DBBroker) Stop() error {
	b.stopWatch()
	defer close(b.watchStop)
	return b.Conn.Close()
}

func (b *DBBroker) setup() (*sql.DB, error) {
	util.Log.Info("Setup PostgreSql Connection")
	if b.Conn != nil {
		return b.Conn, nil
	}

	conn, err := b.connect()
	if err != nil {
		util.Log.Errorf("connection failed: %v", err)
		return conn, err
	}

	b.Lock()
	b.Conn = conn
	b.Unlock()

	errors := make(chan error, 1)
	go func() {
		errors <- err
	}()

	b.errors = errors

	return conn, nil
}

func (b *DBBroker) connect() (*sql.DB, error) {
	configDB := getDBConnectionString()

	cfg := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		configDB.Username,
		configDB.Password,
		configDB.Host,
		configDB.Port,
		configDB.Database,
	)

	conn, err := sql.Open(configDB.Dialect, cfg)

	return conn, err
}

func (b *DBBroker) watch() {
	for {
		select {
		case err := <-b.errors:
			util.Log.Warnf("Connection lost: %v\n", err)
			b.Lock()
			b.Conn = nil
			b.Unlock()
			time.Sleep(10 * time.Second)
			b.setup()
		case stop := <-b.watchStop:
			if stop {
				return
			}
		}
	}
}

func (b *DBBroker) stopWatch() {
	b.watchStop <- true
}
