package logger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sync"
)

// Use this in the init() function to initialize the size of the buffered channel
const logbuffersize int32 = 200

var DEBUG bool = true

var (
	loggerInstance Logger
	loggeronce     sync.Once
	loggerlogonce  sync.Once
)

func Create(instance Logger) {
	loggeronce.Do(func() {
		loggerInstance = instance
		loggerInstance.init()
	})
}

func LoggerInstance() Logger {
	return loggerInstance
}

func PrintJson[T any](entity *T) string {
	typename := reflect.TypeFor[T]().Name()
	outputStringJson, err := json.MarshalIndent((*entity), "", "     ")
	if err != nil {
		return "Error parsing json data"
	} else {
		return typename + ":\n" + string(outputStringJson) + "\n"
	}
}

func ExampleLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		LoggerInstance().WriteDebug(fmt.Sprintf("\tRequest from: %s to: Host: %s URL: %s \tWith HEADERS: %s \tWith BODY: %s", r.RemoteAddr, r.Host, r.URL, r.Header, r.Body))
		next.ServeHTTP(w, r)
	})
}
