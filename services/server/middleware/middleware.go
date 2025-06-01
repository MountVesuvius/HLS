package middleware

// thanks Dreams of Code

import (
	"log"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func CreateChain(ms ...Middleware) Middleware {
    return func(next http.Handler) http.Handler {
        for i := len(ms)-1; i >= 0; i-- {
            m := ms[i]
            next = m(next)
        }

        return next
    }
}

type wrappedWriter struct {
    http.ResponseWriter
    statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
    w.ResponseWriter.WriteHeader(statusCode)
    w.statusCode = statusCode
}

func Logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        start := time.Now()

        wrapped := &wrappedWriter{
            ResponseWriter: w,
            statusCode: http.StatusOK,
        }

        next.ServeHTTP(wrapped, req)

        log.Println(wrapped.statusCode, req.Method, req.URL.Path, time.Since(start))
    })
}

func EnableCors(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        log.Println("Pretend this enabled cors")
        next.ServeHTTP(w, req)
    })
}

func OnlyOnSome(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        log.Println("This is only run on some routes")
        next.ServeHTTP(w, req)
    })
}
