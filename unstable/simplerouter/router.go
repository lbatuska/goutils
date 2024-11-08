package SimpleRouter

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (rg *RouteGroup) StartServer(addr string) *http.Server {
	server := &http.Server{
		Addr:    addr,
		Handler: rg,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("HTTP server error: %v", err))
		}
	}()

	return server
}

func (rg *RouteGroup) StartWithGracefulShutdown(addr string) {
	server := &http.Server{
		Addr:    addr,
		Handler: rg, // This will use your RouteGroup (with ServeMux) as the handler
	}

	// Start the server in a separate goroutine
	go func() {
		fmt.Println("Starting server on", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("HTTP server error: %v", err))
		}
	}()

	// Set up signal catching for graceful shutdown (Ctrl+C, SIGINT, SIGTERM)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-sigChan

	fmt.Println("Shutting down server...")

	// Set a timeout for the shutdown process (10 seconds)
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	// Perform graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		panic(fmt.Sprintf("HTTP shutdown error: %v", err))
	}

	fmt.Println("Graceful shutdown complete.")
}

// http.ListenAndServe takes a Handler interface defined as: ServeHTTP(ResponseWriter, *Request)
func (rg *RouteGroup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rg.mux.ServeHTTP(w, r)
}

// / Handler returns the handler to use for the given request, ... propagate call to http library
func (rg *RouteGroup) Handler(r *http.Request) (h http.Handler, pattern string) {
	return rg.mux.Handler(r)
}

// Creates a new SimpleRouter with the basePath of "/"
func SimpleRouter() *RouteGroup {
	return &RouteGroup{mux: http.NewServeMux(), basePath: "/"}
}

// Puts middlewares on top of existing ones in order
func (rg *RouteGroup) PushMiddleware(first Middleware, others ...Middleware) *RouteGroup {
	rg.middlewares = append(rg.middlewares, first)
	// append supposedly keeps the order of elements, so no extra work is needed
	rg.middlewares = append(rg.middlewares, others...)
	return rg
}

// Removes the top middleware and returns it to the caller
func (rg *RouteGroup) PopMiddleware() Middleware {
	if len(rg.middlewares) == 0 {
		return nil
	}

	lastIndex := len(rg.middlewares) - 1
	lastMiddleware := rg.middlewares[lastIndex]
	rg.middlewares = rg.middlewares[:lastIndex]

	return lastMiddleware
}

// applies the array of middlewares on the handler
func (rg *RouteGroup) applyMiddlewares(handler http.Handler) http.Handler {
	mwlen := len(rg.middlewares) - 1 // We need to start from the last element going to 0 for correct ordering
	for i := range rg.middlewares {
		handler = rg.middlewares[mwlen-i](handler)
	}
	return handler
}

func (rg *RouteGroup) registerRoute(method string, path string, handler http.HandlerFunc) {
	if path == "/" {
		path = rg.basePath[:len(rg.basePath)-1]
	} else {
		path = rg.basePath + path
	}

	fullPath := method + " " + path

	rg.mux.HandleFunc(fullPath, rg.applyMiddlewares(handler).ServeHTTP)
}

func (rg *RouteGroup) SubPath(path string) *RouteGroup {
	rgc := &RouteGroup{
		mux:      rg.mux,
		basePath: rg.basePath + path + "/",
	}
	middlewares := make([]Middleware, len(rg.middlewares))
	copy(middlewares, rg.middlewares)
	rgc.middlewares = middlewares
	return rgc
}

func (rg *RouteGroup) HandleFunc(method string, path string, handler http.HandlerFunc) {
	rg.registerRoute(method, path, handler)
}

func (rg *RouteGroup) Handle(method string, path string, handler http.Handler) {
	rg.HandleFunc(method, path, handler.ServeHTTP)
}

func (rg *RouteGroup) SubpathHandle(path string, handler http.Handler) {
	rg.registerRoute("", path, handler.ServeHTTP)
}

func (rg *RouteGroup) GET(path string, handler http.Handler) {
	rg.Handle(http.MethodGet, path, handler)
}

func (rg *RouteGroup) HEAD(path string, handler http.Handler) {
	rg.Handle(http.MethodHead, path, handler)
}

func (rg *RouteGroup) OPTIONS(path string, handler http.Handler) {
	rg.Handle(http.MethodOptions, path, handler)
}

func (rg *RouteGroup) POST(path string, handler http.Handler) {
	rg.Handle(http.MethodPost, path, handler)
}

func (rg *RouteGroup) PUT(path string, handler http.Handler) {
	rg.Handle(http.MethodPut, path, handler)
}

func (rg *RouteGroup) PATCH(path string, handler http.Handler) {
	rg.Handle(http.MethodPatch, path, handler)
}

func (rg *RouteGroup) DELETE(path string, handler http.Handler) {
	rg.Handle(http.MethodDelete, path, handler)
}
