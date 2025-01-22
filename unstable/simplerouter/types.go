package SimpleRouter

import "net/http"

type Middleware func(http.Handler) http.Handler

type RouteGroup struct {
	mux                *http.ServeMux // A ponter to the underlying ServeMux, this allows us to call ListenAndServe on any instance
	basePath           string         // The current path we are defining handlers on / appending to
	middlewares        []Middleware   // Stack of middlewares that will be applied on a handler in order
	global_middlewares *[]Middleware
}
