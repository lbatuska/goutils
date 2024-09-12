package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func (lgr *FileLoggerimpl) init() {
	lgr.filepath = "./log"
	lgr.messages = make(chan string, logbuffersize)
	envfp, envexist := os.LookupEnv("LOGFILE_GO_LOGGER")
	if envexist {
		if len(envfp) > 0 {
			lgr.filepath = envfp
		} else {
			Logger().Write_DEBUG(fmt.Sprintf("LOGFILE_GO_LOGGER env exist but has an empty value using default value: %s !\n", lgr.filepath))
		}
	} else {
		Logger().Write_DEBUG(fmt.Sprintf("LOGFILE_GO_LOGGER env doesn't exist using default value: %s !\n", lgr.filepath))
	}
	f, err := os.OpenFile(lgr.filepath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		// We probably really don't want to continue execution without file backed logging
		panic(fmt.Sprintf("Error opening or creating file: %s", err.Error()))
	}
	lgr.logFile = f
	lgr.mutex = &sync.Mutex{}
}

func (logger *FileLoggerimpl) StartLogger() {
	fmt.Println("Starting FileLogger")
	loggerlogonce.Do(func() {
		for msg := range logger.messages {

			logger.mutex.Lock()
			_, err := logger.logFile.WriteString(msg)
			if err != nil {
				fmt.Println(err.Error())
				logger.logFile.Close()
				logger.mutex.Unlock()
				panic("Failed to write to file")
			}
			err = logger.logFile.Sync()
			if err != nil {
				logger.logFile.Close()
				logger.mutex.Unlock()

				panic("Failed to write to file")
			}
			logger.mutex.Unlock()
		}
	})
	// Technically we should do this but this will never run
	// logger.mutex.Lock()
	// logger.logFile.Close()
	// logger.mutex.Unlock()
}

func (logger *FileLoggerimpl) Write(message string) {
	logger.messages <- time.Now().Format(time.UnixDate) + " : " + message + "\n"
}

func (logger *FileLoggerimpl) Write_Request(message string, request string) {
	logger.Write(request + " : " + message)
}

func (logger *FileLoggerimpl) WriteErr(err error) (errnum int) {
	if err != nil {
		logger.Write("Error: " + err.Error())
		errnum = 1
	}
	return errnum
}

func (logger *FileLoggerimpl) WriteErr_Request(err error, uuid string) (errnum int) {
	if err != nil {
		logger.Write(uuid + " : Error: " + err.Error())
		errnum = 1
	}
	return errnum
}

func (logger *FileLoggerimpl) Write_DEBUG(message string) {
	if DEBUG {
		logger.Write(message)
	}
}

func (logger *FileLoggerimpl) Write_Request_DEBUG(message string, uuid string) {
	if DEBUG {
		logger.Write_Request(message, uuid)
	}
}

func (logger *FileLoggerimpl) WriteErr_DEBUG(err error) (errnum int) {
	if err != nil {
		logger.Write_DEBUG("Error: " + err.Error())
		errnum = 1
	}
	return errnum
}

func (logger *FileLoggerimpl) WriteErr_Request_DEBUG(err error, uuid string) (errnum int) {
	if err != nil {
		logger.Write_DEBUG(uuid + " : Error: " + err.Error())
		errnum = 1
	}
	return errnum
}
