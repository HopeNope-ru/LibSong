package middlewares

import (
	"bytes"
	"io"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if l := log.Debug(); l.Enabled() {
			dict := zerolog.Dict()
			for key, val := range r.Header {
				dict = dict.Strs(key, val)
			}

			// Логировать тело запроса затратное дело.
			// Поэтому пускай будет в дебаге
			buf := new(bytes.Buffer)
			buf.Grow(255)

			_, err := io.Copy(buf, r.Body)
			if err != nil {
				http.Error(w, "err", http.StatusInternalServerError)
			}

			// Надо изучить либу zerolog, чтобы не копипастить
			// path и method. Возможное решение через context.
			l.Str("method", r.Method).
				Str("path", r.URL.Path).
				Dict("headers", dict).
				Bytes("body", buf.Bytes()).
				Interface("query", r.URL.Query()).
				Msg("logging")

			r.Body = io.NopCloser(buf)
		} else {
			log.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Msg("logging")
		}

		next.ServeHTTP(w, r)
	})
}
