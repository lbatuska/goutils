package logger

func (lgr *NullLoggerimpl) init() {}

func (logger *NullLoggerimpl) StartLogger() {}

func (logger *NullLoggerimpl) Write(message string) {}

func (logger *NullLoggerimpl) Write_Request(message string, request string) {}

func (logger *NullLoggerimpl) WriteErr(err error) (errnum int) {
	if err != nil {
		errnum = 1
	}
	return errnum
}

func (logger *NullLoggerimpl) WriteErr_Request(err error, request string) (errnum int) {
	if err != nil {
		errnum = 1
	}
	return errnum
}

func (logger *NullLoggerimpl) Write_DEBUG(message string) {}

func (logger *NullLoggerimpl) Write_Request_DEBUG(message string, uuid string) {}

func (logger *NullLoggerimpl) WriteErr_DEBUG(err error) (errnum int) {
	if err != nil {
		errnum = 1
	}
	return errnum
}

func (logger *NullLoggerimpl) WriteErr_Request_DEBUG(err error, uuid string) (errnum int) {
	if err != nil {
		errnum = 1
	}
	return errnum
}
