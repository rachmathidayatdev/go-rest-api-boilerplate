package util

type (
	//Daemon interface
	Daemon interface {
		Start() error
		Stop() error
	}
)
