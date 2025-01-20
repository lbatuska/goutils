package Logger

// Message format(s)
//
// Unixdate : message\n
//
// Unixdate : Error: error\n
//
// Unixdate : uuid : message\n
//
// Unixdate : uuid : Error: error\n
type Logger interface {
	ReleaseLogger
	DebugLogger
}

type (
	ReleaseLogger interface {
		// Private, use it for member initialization etc
		init()
		// Start an infinite loop to write out messages from the channel
		StartLogger()
		StopLogger()
		Write(message string)
		WriteRequest(message string, uuid string)
		// If an error that is not nill passed in it logs the error and returns 1, otherwise 0
		WriteErr(error) int
		WriteErrRequest(err error, uuid string) int

		WriteErrMsgRequest(err error, message string, uuid string) int
	}
	// Use _DEBUG prints to strip them out of release builds
	DebugLogger interface {
		// Private, use it for member initialization etc
		init()
		// Start an infinite loop to write out messages from the channel
		StartLogger()
		StopLogger()
		WriteDebug(message string)
		WriteRequestDebug(message string, uuid string)
		WriteErrDebug(err error) (errnum int)
		WriteErrRequestDebug(err error, uuid string) int

		WriteErrMsgRequestDebug(err error, message string, uuid string) int
	}
)

// Ensure all methods from LGRImpl are implemented ccompile time
var (
	_ Logger = (*NullLoggerImpl)(nil)
	_ Logger = (*ConsoleLoggerImpl)(nil)
	_ Logger = (*FileLoggerImpl)(nil)
	_ Logger = (*SlogLoggerImpl)(nil)
)
