package middlewares

import (
	"bytes"
	"io"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type HookMaster struct {
	Body
	Method
	Path
	http.Header
	zerolog.Hook
}

func (h HookMaster) Run(e *zerolog.Event, level zerolog.Level, message string) {
	switch level {
	case zerolog.DebugLevel:

	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		compatibleLog := new(zerolog.Event)
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

			l.Dict("headers", dict).
				Bytes("body", buf.Bytes())
			r.Body = io.NopCloser(buf)
		}

		log.Warn().Msg("")

		l.Str("path", r.URL.Path).
			Str("method", r.Method)
		// Релизация логгирования с уровнем DEBUG
		// v := mux.Vars(r)
		// Реализаци логгирования с уровнем INFO
		next.ServeHTTP(w, r)
	})
}
