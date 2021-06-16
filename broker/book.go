package broker

import "database/sql"

type (
	//BookBroker interface
	BookBroker interface {
		Start() (*sql.DB, error)
		Stop() error
	}
)
