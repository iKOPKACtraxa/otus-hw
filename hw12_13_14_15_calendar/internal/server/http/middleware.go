package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

// LoggingMiddleware is for logging at requests.
func LoggingMiddleware(app Application, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		next.ServeHTTP(w, r)
		latency := fmt.Sprintf("%vms", time.Since(now).Milliseconds())
		app.Info(fmt.Sprintf("Вывод в Logger: %v [%v] %v %v %v %v %v %v",
			r.RemoteAddr,
			now.Format(time.RFC822Z),
			r.Method,
			r.RequestURI,
			r.Proto,
			http.StatusOK,
			latency,
			r.UserAgent(),
		))
	})
}
