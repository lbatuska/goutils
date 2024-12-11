package Logger

func (lgr *NullLoggerImpl) init() {}

func (logger *NullLoggerImpl) StartLogger() {}

func (logger *NullLoggerImpl) StopLogger() {}

func (logger *NullLoggerImpl) Write(message string) {}

func (logger *NullLoggerImpl) WriteRequest(message string, request string) {}

func (logger *NullLoggerImpl) WriteErr(err error) (errnum int) {
	if err != nil {
		errnum = 1
	}
	return errnum
}

func (logger *NullLoggerImpl) WriteErrRequest(err error, request string) (errnum int) {
	if err != nil {
		errnum = 1
	}
	return errnum
}

func (logger *NullLoggerImpl) WriteDebug(message string) {}

func (logger *NullLoggerImpl) WriteRequestDebug(message string, uuid string) {}

func (logger *NullLoggerImpl) WriteErrDebug(err error) (errnum int) {
	if err != nil {
		errnum = 1
	}
	return errnum
}

func (logger *NullLoggerImpl) WriteErrRequestDebug(err error, uuid string) (errnum int) {
	if err != nil {
		errnum = 1
	}
	return errnum
}
