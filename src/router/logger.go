package router

import (
	"bytes"
	"fmt"
	"git"
	"log"
	"net/http"
	"time"
)

var LogBuffer *bytes.Buffer

/*func init() {
	LogBuffer = new(bytes.Buffer)
	log.SetOutput(LogBuffer)
}*/

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
			fmt.Sprint(git.RateLimitRemaining, "/", git.RateLimit),
		)
	})
}
