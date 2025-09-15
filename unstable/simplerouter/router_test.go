package SimpleRouter

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"syscall"
	"testing"
	"time"

	Testing "github.com/lbatuska/goutils/testing"
	Type "github.com/lbatuska/goutils/type"
)

func setupServer(wg *sync.WaitGroup) {
	defer wg.Done()

	simpleRouter := SimpleRouter()
	simpleRouter.GET("test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 5)
		w.WriteHeader(200)
	}))
	simpleRouter.GET("test2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	simpleRouter.StartWithGracefulShutdown("0.0.0.0:8080", Type.None[ServerConfig]())
}

func sendKillSignal(wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(time.Second * 4)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}

func Test_gracefulShutdown(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	go setupServer(&wg)
	go sendKillSignal(&wg)

	time.Sleep(time.Second * 2)
	res, err := http.Get("http://localhost:8080/test")
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, http.StatusOK, res.StatusCode)

	wg.Wait()
}

func Test_middlewares(t *testing.T) {
	counter := 0
	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter += 1
			next.ServeHTTP(w, r)
		})
	}
	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter = counter * 2
			next.ServeHTTP(w, r)
		})
	}
	simpleRouter := SimpleRouter()
	simpleRouter.PushMiddleware(mw1)
	simpleRouter.PushMiddleware(mw2)
	simpleRouter.GET("test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	ts := httptest.NewServer(simpleRouter)
	req, _ := http.NewRequest("GET", ts.URL+"/test", nil)
	_, err := http.DefaultClient.Do(req)

	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, 2, counter)
	ts.Close()
	counter = 0

	simpleRouter.PopMiddleware()
	simpleRouter.PopMiddleware()
	simpleRouter.PushMiddleware(mw2)
	simpleRouter.PushMiddleware(mw1)
	simpleRouter.GET("test2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	ts = httptest.NewServer(simpleRouter)
	req, _ = http.NewRequest("GET", ts.URL+"/test2", nil)
	_, err = http.DefaultClient.Do(req)

	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, 1, counter)
	ts.Close()
	counter = 0

	simpleRouter.PopMiddleware()
	simpleRouter.PopMiddleware()
	simpleRouter.PushMiddleware(mw1)
	simpleRouter.PushMiddleware(mw2)
	simpleRouter.PushMiddleware(mw1)
	simpleRouter.PushMiddleware(mw2)
	simpleRouter.GET("test3", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	ts = httptest.NewServer(simpleRouter)
	defer ts.Close()
	req, _ = http.NewRequest("GET", ts.URL+"/test3", nil)
	_, err = http.DefaultClient.Do(req)

	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, 6, counter)
	ts.Close()
}

func Test_middlewares2(t *testing.T) {
	counter := 0
	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter += 1
			next.ServeHTTP(w, r)
		})
	}
	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter = counter * 2
			next.ServeHTTP(w, r)
		})
	}
	simpleRouter := SimpleRouter()
	simpleRouter.PushMiddleware(mw1, mw2)
	simpleRouter.GET("test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	ts := httptest.NewServer(simpleRouter)
	req, _ := http.NewRequest("GET", ts.URL+"/test", nil)
	_, err := http.DefaultClient.Do(req)

	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, 2, counter)
	ts.Close()
	counter = 0

	simpleRouter.PopMiddleware()
	simpleRouter.PopMiddleware()
	simpleRouter.PushMiddleware(mw2, mw1)
	simpleRouter.GET("test2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	ts = httptest.NewServer(simpleRouter)
	req, _ = http.NewRequest("GET", ts.URL+"/test2", nil)
	_, err = http.DefaultClient.Do(req)

	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, 1, counter)
	ts.Close()
	counter = 0

	simpleRouter.PopMiddleware()
	simpleRouter.PopMiddleware()
	simpleRouter.PushMiddleware(mw1, mw2, mw1, mw2)
	simpleRouter.GET("test3", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	ts = httptest.NewServer(simpleRouter)
	defer ts.Close()
	req, _ = http.NewRequest("GET", ts.URL+"/test3", nil)
	_, err = http.DefaultClient.Do(req)

	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, 6, counter)
	ts.Close()
}

func Test_globalMiddlewares(t *testing.T) {
	counter := 0
	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter += 1
			next.ServeHTTP(w, r)
		})
	}
	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter = counter * 2
			next.ServeHTTP(w, r)
		})
	}
	gmw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter += 3
			next.ServeHTTP(w, r)
		})
	}
	gmw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter = counter * 4
			next.ServeHTTP(w, r)
		})
	}
	simpleRouter := SimpleRouter()
	simpleRouter2 := simpleRouter.SubPath("2")

	simpleRouter.PushGlobalMiddleware(gmw1)
	simpleRouter.PushGlobalMiddleware(gmw2)
	simpleRouter.PushMiddleware(mw1, mw2)
	simpleRouter.GET("test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	simpleRouter.GET("test2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	// Subpath
	simpleRouter2.GET("test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	_ts := httptest.NewServer(simpleRouter2)
	_req, _ := http.NewRequest("GET", _ts.URL+"/2/test", nil)
	_, _err := http.DefaultClient.Do(_req)
	Testing.AssertNotError(t, _err)
	Testing.AssertEqual(t, 12, counter)
	counter = 0
	// Subpath

	ts := httptest.NewServer(simpleRouter)
	req, _ := http.NewRequest("GET", ts.URL+"/test", nil)
	_, err := http.DefaultClient.Do(req)
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, 26, counter)
	counter = 0

	req2, _ := http.NewRequest("GET", ts.URL+"/test2", nil)
	_, err2 := http.DefaultClient.Do(req2)
	Testing.AssertNotError(t, err2)
	Testing.AssertEqual(t, 26, counter)
	ts.Close()
	counter = 0

	simpleRouter.PopGlobalMiddleware()
	simpleRouter.PopGlobalMiddleware()
	simpleRouter.PushGlobalMiddleware(gmw2)
	simpleRouter.PushGlobalMiddleware(gmw1)

	ts = httptest.NewServer(simpleRouter)
	req, _ = http.NewRequest("GET", ts.URL+"/test", nil)
	_, err = http.DefaultClient.Do(req)
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, 8, counter)
	counter = 0

	req2, _ = http.NewRequest("GET", ts.URL+"/test2", nil)
	_, err2 = http.DefaultClient.Do(req2)
	Testing.AssertNotError(t, err2)
	Testing.AssertEqual(t, 8, counter)
	ts.Close()
	counter = 0
}

func Test_globalMiddlewares2(t *testing.T) {
	counter := 0
	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter += 1
			next.ServeHTTP(w, r)
		})
	}
	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter = counter * 2
			next.ServeHTTP(w, r)
		})
	}
	gmw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter += 3
			next.ServeHTTP(w, r)
		})
	}
	gmw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter = counter * 4
			next.ServeHTTP(w, r)
		})
	}
	simpleRouter := SimpleRouter()
	simpleRouter.PushGlobalMiddleware(gmw1, gmw2)
	simpleRouter.PushMiddleware(mw1, mw2)
	simpleRouter.GET("test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	simpleRouter.GET("test2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	ts := httptest.NewServer(simpleRouter)
	req, _ := http.NewRequest("GET", ts.URL+"/test", nil)
	_, err := http.DefaultClient.Do(req)
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, 26, counter)
	counter = 0

	req2, _ := http.NewRequest("GET", ts.URL+"/test2", nil)
	_, err2 := http.DefaultClient.Do(req2)
	Testing.AssertNotError(t, err2)
	Testing.AssertEqual(t, 26, counter)
	ts.Close()
	counter = 0

	simpleRouter.PopGlobalMiddleware()
	simpleRouter.PopGlobalMiddleware()
	simpleRouter.PushGlobalMiddleware(gmw2, gmw1)

	ts = httptest.NewServer(simpleRouter)
	req, _ = http.NewRequest("GET", ts.URL+"/test", nil)
	_, err = http.DefaultClient.Do(req)
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, 8, counter)
	counter = 0

	req2, _ = http.NewRequest("GET", ts.URL+"/test2", nil)
	_, err2 = http.DefaultClient.Do(req2)
	Testing.AssertNotError(t, err2)
	Testing.AssertEqual(t, 8, counter)
	ts.Close()
	counter = 0
}

func Test_handle(t *testing.T) {
	counter := 0
	subpathHandler := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter = 1
		})
	}
	regularHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counter = 2
	})
	simpleRouter := SimpleRouter()
	simpleRouter.SubpathHandle("test/", subpathHandler())
	simpleRouter.Handle("GET", "test2", regularHandler)

	ts := httptest.NewServer(simpleRouter)
	defer ts.Close()
	req, _ := http.NewRequest("GET", ts.URL+"/test/1", nil)
	res, err := http.DefaultClient.Do(req)
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, http.StatusOK, res.StatusCode)
	Testing.AssertEqual(t, 1, counter)

	req, _ = http.NewRequest("GET", ts.URL+"/test/2", nil)
	res, err = http.DefaultClient.Do(req)
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, http.StatusOK, res.StatusCode)
	Testing.AssertEqual(t, 1, counter)

	req, _ = http.NewRequest("GET", ts.URL+"/test2", nil)
	res, err = http.DefaultClient.Do(req)
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, http.StatusOK, res.StatusCode)
	Testing.AssertEqual(t, 2, counter)
}

func Test_methods(t *testing.T) {
	simpleRouter := SimpleRouter()
	simpleRouter.GET("test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	simpleRouter.POST("test2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	simpleRouter.PATCH("test3", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	simpleRouter.PUT("test4", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	simpleRouter.DELETE("test5", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	simpleRouter.OPTIONS("test6", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	simpleRouter.HEAD("test7", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	ts := httptest.NewServer(simpleRouter)
	defer ts.Close()
	req, _ := http.NewRequest("GET", ts.URL+"/test", nil)
	req2, _ := http.NewRequest("POST", ts.URL+"/test", nil)
	res, err := http.DefaultClient.Do(req)
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, http.StatusOK, res.StatusCode)
	res2, err := http.DefaultClient.Do(req2)
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, http.StatusMethodNotAllowed, res2.StatusCode)

	req, _ = http.NewRequest("POST", ts.URL+"/test2", nil)
	req2, _ = http.NewRequest("PATCH", ts.URL+"/test2", nil)
	res, err = http.DefaultClient.Do(req)
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, http.StatusOK, res.StatusCode)
	res2, err = http.DefaultClient.Do(req2)
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, http.StatusMethodNotAllowed, res2.StatusCode)

	req, _ = http.NewRequest("PATCH", ts.URL+"/test3", nil)
	req2, _ = http.NewRequest("PUT", ts.URL+"/test3", nil)
	res, err = http.DefaultClient.Do(req)
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, http.StatusOK, res.StatusCode)
	res2, err = http.DefaultClient.Do(req2)
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, http.StatusMethodNotAllowed, res2.StatusCode)

	req, _ = http.NewRequest("PUT", ts.URL+"/test4", nil)
	req2, _ = http.NewRequest("DELETE", ts.URL+"/test4", nil)
	res, err = http.DefaultClient.Do(req)
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, http.StatusOK, res.StatusCode)
	res2, err = http.DefaultClient.Do(req2)
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, http.StatusMethodNotAllowed, res2.StatusCode)

	req, _ = http.NewRequest("DELETE", ts.URL+"/test5", nil)
	req2, _ = http.NewRequest("OPTIONS", ts.URL+"/test5", nil)
	res, err = http.DefaultClient.Do(req)
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, http.StatusOK, res.StatusCode)
	res2, err = http.DefaultClient.Do(req2)
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, http.StatusMethodNotAllowed, res2.StatusCode)

	req, _ = http.NewRequest("OPTIONS", ts.URL+"/test6", nil)
	req2, _ = http.NewRequest("HEAD", ts.URL+"/test6", nil)
	res, err = http.DefaultClient.Do(req)
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, http.StatusOK, res.StatusCode)
	res2, err = http.DefaultClient.Do(req2)
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, http.StatusMethodNotAllowed, res2.StatusCode)

	req, _ = http.NewRequest("HEAD", ts.URL+"/test7", nil)
	req2, _ = http.NewRequest("GET", ts.URL+"/test7", nil)
	res, err = http.DefaultClient.Do(req)
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, http.StatusOK, res.StatusCode)
	res2, err = http.DefaultClient.Do(req2)
	Testing.AssertNotError(t, err)
	Testing.AssertEqual(t, http.StatusMethodNotAllowed, res2.StatusCode)
}
