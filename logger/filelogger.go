package Logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func (lgr *FileLoggerImpl) SetLogFilePath(path string) {
	lgr.initfilepath = path
}

func (lgr *FileLoggerImpl) init() {
	if lgr.initfilepath == "" {
		lgr.filepath = "./log"
	} else {
		lgr.filepath = lgr.initfilepath
	}
	lgr.messages = make(chan string, Logbuffersize)
	envfp, envexist := os.LookupEnv("LOGFILE_GO_LOGGER")
	if envexist {
		if len(envfp) > 0 {
			lgr.filepath = envfp
		} else {
			LoggerInstance().WriteDebug(fmt.Sprintf("LOGFILE_GO_LOGGER env exist but has an empty value using default value: %s !\n", lgr.filepath))
		}
	} else {
		LoggerInstance().WriteDebug(fmt.Sprintf("LOGFILE_GO_LOGGER env doesn't exist using default value: %s !\n", lgr.filepath))
	}
	f, err := os.OpenFile(lgr.filepath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		// We probably really don't want to continue execution without file backed logging
		panic(fmt.Sprintf("Error opening or creating file: %s", err.Error()))
	}
	lgr.logFile = f
	lgr.mutex = &sync.Mutex{}
}

func (logger *FileLoggerImpl) StartLogger() {
	fmt.Println("Starting FileLogger")
	loggerlogonce.Do(func() {
		go func() {
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
		}()
	})
	// Technically we should do this but this will never run
	// logger.mutex.Lock()
	// logger.logFile.Close()
	// logger.mutex.Unlock()
}

func (logger *FileLoggerImpl) StopLogger() {
	close(logger.messages)
}

func (logger *FileLoggerImpl) Write(message string) {
	logger.messages <- time.Now().Format(time.UnixDate) + " : " + message + "\n"
}

func (logger *FileLoggerImpl) WriteRequest(message string, uuid string) {
	logger.Write(uuid + " : " + message)
}

func (logger *FileLoggerImpl) WriteErr(err error) (errnum int) {
	if err != nil {
		logger.Write("Error: " + err.Error())
		errnum = 1
	}
	return errnum
}

func (logger *FileLoggerImpl) WriteErrRequest(err error, uuid string) (errnum int) {
	if err != nil {
		logger.Write(uuid + " : Error: " + err.Error())
		errnum = 1
	}
	return errnum
}

func (logger *FileLoggerImpl) WriteDebug(message string) {
	if DEBUG {
		logger.Write(message)
	}
}

func (logger *FileLoggerImpl) WriteRequestDebug(message string, uuid string) {
	if DEBUG {
		logger.WriteRequest(message, uuid)
	}
}

func (logger *FileLoggerImpl) WriteErrDebug(err error) (errnum int) {
	if err != nil {
		logger.WriteDebug("Error: " + err.Error())
		errnum = 1
	}
	return errnum
}

func (logger *FileLoggerImpl) WriteErrRequestDebug(err error, uuid string) (errnum int) {
	if err != nil {
		logger.WriteDebug(uuid + " : Error: " + err.Error())
		errnum = 1
	}
	return errnum
}
