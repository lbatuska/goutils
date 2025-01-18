package Logger

import "log/slog"

func (lgr *SlogLoggerImpl) init() {}

func (logger *SlogLoggerImpl) StartLogger() {}

func (logger *SlogLoggerImpl) StopLogger() {}

func (logger *SlogLoggerImpl) Write(message string) {
	slog.Info(message)
}

func (logger *SlogLoggerImpl) WriteRequest(message string, uuid string) {
	slog.Info(message, "UUID", uuid)
}

func (logger *SlogLoggerImpl) WriteErr(err error) (errnum int) {
	if err != nil {
		slog.Error(err.Error())
		errnum = 1
	}
	return errnum
}

func (logger *SlogLoggerImpl) WriteErrRequest(err error, uuid string) (errnum int) {
	if err != nil {
		slog.Error(err.Error(), "UUID", uuid)
		errnum = 1
	}
	return errnum
}

func (logger *SlogLoggerImpl) WriteDebug(message string) {
	if DEBUG {
		logger.Write(message)
	}
}

func (logger *SlogLoggerImpl) WriteRequestDebug(message string, uuid string) {
	if DEBUG {
		logger.WriteRequest(message, uuid)
	}
}

func (logger *SlogLoggerImpl) WriteErrDebug(err error) (errnum int) {
	if err != nil {
		if DEBUG {
			logger.WriteErr(err)
		}
		errnum = 1
	}
	return errnum
}

func (logger *SlogLoggerImpl) WriteErrRequestDebug(err error, uuid string) (errnum int) {
	if err != nil {
		if DEBUG {
			logger.WriteErrRequest(err, uuid)
		}
		errnum = 1
	}
	return errnum
}
