package Logger

import (
	"fmt"
	"time"
)

func (lgr *ConsoleLoggerImpl) init() {
	lgr.messages = make(chan string, Logbuffersize)
}

func (logger *ConsoleLoggerImpl) StartLogger() {
	fmt.Println("Starting Logger")
	loggerlogonce.Do(func() {
		go func() {
			for msg := range logger.messages {
				fmt.Print(msg)
			}
		}()
	})
}

func (logger *ConsoleLoggerImpl) StopLogger() {
	close(logger.messages)
}

func (logger *ConsoleLoggerImpl) Write(message string) {
	logger.messages <- time.Now().Format(time.UnixDate) + " : " + message + "\n"
}

func (logger *ConsoleLoggerImpl) WriteRequest(message string, uuid string) {
	logger.Write(uuid + " : " + message)
}

func (logger *ConsoleLoggerImpl) WriteErr(err error) (errnum int) {
	if err != nil {
		logger.Write("Error: " + err.Error())
		errnum = 1
	}
	return errnum
}

func (logger *ConsoleLoggerImpl) WriteErrRequest(err error, uuid string) (errnum int) {
	if err != nil {
		logger.Write(uuid + " : Error: " + err.Error())
		errnum = 1
	}
	return errnum
}

func (logger *ConsoleLoggerImpl) WriteDebug(message string) {
	if DEBUG {
		logger.Write(message)
	}
}

func (logger *ConsoleLoggerImpl) WriteRequestDebug(message string, uuid string) {
	if DEBUG {
		logger.WriteRequest(message, uuid)
	}
}

func (logger *ConsoleLoggerImpl) WriteErrDebug(err error) (errnum int) {
	if err != nil {
		logger.WriteDebug("Error: " + err.Error())
		errnum = 1
	}
	return errnum
}

func (logger *ConsoleLoggerImpl) WriteErrRequestDebug(err error, uuid string) (errnum int) {
	if err != nil {
		logger.WriteDebug(uuid + " : Error: " + err.Error())
		errnum = 1
	}
	return errnum
}
