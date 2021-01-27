package internalhttp

import (
	"net/http"
	"strings"
	"time"

	"github.com/eafomi4ev/golang_otus_homework/hw12_13_14_15_calendar/internal/app"
)

var dateLayout = "[02/Jan/2006:15:04:05 -0700]"

func loggingMiddleware(next http.Handler, logg app.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		logg.Info(strings.Join( // todo: добавить логирование кода ответа
			[]string{
				r.RemoteAddr,
				time.Now().Format(dateLayout),
				r.Method,
				r.URL.Path,
				r.Proto,
				r.Host,
				time.Since(start).String(),
				r.UserAgent(),
			},
			" ",
		))

		next.ServeHTTP(w, r)
	})
}
