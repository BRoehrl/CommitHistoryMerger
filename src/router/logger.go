package router

import (
	"fmt"
	"git"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"user"
)

func init() {
	logfile, err := os.OpenFile("commit-finder.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		multi := io.MultiWriter(os.Stdout)
		log.SetOutput(multi)
		log.Fatalf("error opening file: %v", err)
	} else {
		multi := io.MultiWriter(os.Stdout, logfile)
		log.SetOutput(multi)
	}
}

// Logger wraps a log output around a http handler
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		userID, _ := checkJWTandGetUserID(r)
		inner.ServeHTTP(w, r)

		duration := time.Since(start)

		authToken, _ := user.GetAccessToken(userID)
		metaData := git.AuthTokenToLastResponse[authToken]
		log.Printf(
			"%s\t%s\t%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			duration,
			userID,
			fmt.Sprint(metaData.RateLimitRemaining, "/", metaData.RateLimit),
		)
	})
}
