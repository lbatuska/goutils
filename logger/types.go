package Logger

import (
	"os"
	"sync"
)

// A logger without logging functionality
type NullLoggerImpl struct{}

// A logger that logs to sdtout
type ConsoleLoggerImpl struct {
	messages chan string
}

type FileLoggerImpl struct {
	messages     chan string
	mutex        *sync.Mutex
	logFile      *os.File
	filepath     string
	initfilepath string
}

type SlogLoggerImpl struct{}
