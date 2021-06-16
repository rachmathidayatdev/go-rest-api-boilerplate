package application

import (
	"time"

	"github.com/go-rest-api-boilerplate/broker"
	bookController "github.com/go-rest-api-boilerplate/server/book/controller"
	bookRepository "github.com/go-rest-api-boilerplate/server/book/repository"
	bookService "github.com/go-rest-api-boilerplate/server/book/service"
	"github.com/go-rest-api-boilerplate/util"
)

type (
	bookApp struct {
		broker   broker.BookBroker
		stopChan chan bool
	}
)

//NewBookDaemon func
func (app *Application) NewBookDaemon() util.Daemon {
	broker := NewDBBroker(app.postgresql)

	return &bookApp{
		broker,
		make(chan bool),
	}
}

func (d *bookApp) Start() error {
	conn, err := d.broker.Start()
	if err != nil {
		return err
	}

	repository := bookRepository.NewBookRepository(conn)
	usecase := bookService.NewBookService(repository)
	bookController.NewBookServer(usecase)

	go d.runLoop()

	return nil
}

func (d *bookApp) runLoop() {
	logger := util.Log.WithField("contex", "bookApp")
	for {
		select {
		default:
			logger.Debug("bookApp started")
			time.Sleep(time.Second * 1)
		case stop := <-d.stopChan:
			if stop {
				return
			}
		}
	}
}

func (d *bookApp) Stop() error {
	d.stopChan <- true
	return d.broker.Stop()
}
