package logger

import (
	"fmt"
	"time"
)

func (lgr *ConsoleLoggerimpl) init() {
	lgr.messages = make(chan string, logbuffersize)
}

func (logger *ConsoleLoggerimpl) StartLogger() {
	fmt.Println("Starting Logger")
	loggerlogonce.Do(func() {
		for msg := range logger.messages {
			fmt.Print(msg)
		}
	})
}

func (logger *ConsoleLoggerimpl) Write(message string) {
	logger.messages <- time.Now().Format(time.UnixDate) + " : " + message + "\n"
}

func (logger *ConsoleLoggerimpl) Write_Request(message string, uuid string) {
	logger.Write(uuid + " : " + message)
}

func (logger *ConsoleLoggerimpl) WriteErr(err error) (errnum int) {
	if err != nil {
		logger.Write("Error: " + err.Error())
		errnum = 1
	}
	return errnum
}

func (logger *ConsoleLoggerimpl) WriteErr_Request(err error, uuid string) (errnum int) {
	if err != nil {
		logger.Write(uuid + " : Error: " + err.Error())
		errnum = 1
	}
	return errnum
}

func (logger *ConsoleLoggerimpl) Write_DEBUG(message string) {
	if DEBUG {
		logger.Write(message)
	}
}

func (logger *ConsoleLoggerimpl) Write_Request_DEBUG(message string, uuid string) {
	if DEBUG {
		logger.Write_Request(message, uuid)
	}
}

func (logger *ConsoleLoggerimpl) WriteErr_DEBUG(err error) (errnum int) {
	if err != nil {
		logger.Write_DEBUG("Error: " + err.Error())
		errnum = 1
	}
	return errnum
}

func (logger *ConsoleLoggerimpl) WriteErr_Request_DEBUG(err error, uuid string) (errnum int) {
	if err != nil {
		logger.Write_DEBUG(uuid + " : Error: " + err.Error())
		errnum = 1
	}
	return errnum
}
