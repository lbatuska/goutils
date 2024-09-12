package logger

import (
	"os"
	"sync"
)

// A logger without logging functionality
type NullLoggerimpl struct{}

// A logger that logs to sdtout
type ConsoleLoggerimpl struct {
	messages chan string
}

type FileLoggerimpl struct {
	messages chan string
	mutex    *sync.Mutex
	logFile  *os.File
	filepath string
}
