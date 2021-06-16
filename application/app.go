package application

import (
	"os"

	"github.com/go-rest-api-boilerplate/util"
)

type (
	//Application struct
	Application struct {
		postgresql *DBConfig
	}

	//Logger struct
	Logger struct {
		// Stdout is true if the output needs to goto standard out
		Stdout bool `yaml:"stdout"`
		// Level is the desired log level
		Level string `yaml:"level"`
		// OutputFile is the path to the log output file
		OutputFile string `yaml:"outputFile"`
	}
)

//SetupApp func
func SetupApp() *Application {
	dbConfig := DBConfig{
		Dialect:  "postgres",
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASS"),
		Database: os.Getenv("DB_NAME"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Charset:  "utf8",
	}

	logger := Logger{Stdout: true, Level: "DEBUG"}
	util.Log = logger.NewLogger()

	return &Application{
		&dbConfig,
	}
}
