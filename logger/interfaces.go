package logger

// Message format(s)
//
// Unixdate : message\n
//
// Unixdate : Error: error\n
//
// Unixdate : uuid : message\n
//
// Unixdate : uuid : Error: error\n
type LGRImpl interface {
	LoggerI
	DebugLoggerI
}

type (
	LoggerI interface {
		// Private, use it for member initialization etc
		init()
		// Start an infinite loop to write out messages from the channel
		StartLogger()
		Write(message string)
		Write_Request(message string, uuid string)
		// If an error that is not nill passed in it logs the error and returns 1, otherwise 0
		WriteErr(error) int
		WriteErr_Request(err error, uuid string) int
	}
	// Use _DEBUG prints to strip them out of release builds
	DebugLoggerI interface {
		// Private, use it for member initialization etc
		init()
		// Start an infinite loop to write out messages from the channel
		StartLogger()
		Write_DEBUG(message string)
		Write_Request_DEBUG(message string, uuid string)
		WriteErr_DEBUG(err error) (errnum int)
		WriteErr_Request_DEBUG(err error, uuid string) int
	}
)

// Ensure all methods from LGRImpl are implemented ccompile time
var (
	_ LGRImpl = (*NullLoggerimpl)(nil)
	_ LGRImpl = (*ConsoleLoggerimpl)(nil)
	_ LGRImpl = (*FileLoggerimpl)(nil)
)
