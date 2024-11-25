package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lyric/songs/hw/src/configs"
	"github.com/lyric/songs/hw/src/handlers"
	"github.com/lyric/songs/hw/src/middlewares"
	"github.com/lyric/songs/hw/src/repository"
	"github.com/rs/zerolog"
)

func GetApp(args ...string) {
	conf := configs.New()
	zerolog.SetGlobalLevel(conf.Levellog)

	dbpool, err := pgxpool.New(context.Background(), conf.Dsn)
	if err != nil {
		panic(err)
	}
	defer dbpool.Close()

	songhandlers := handlers.NewSongHandler(
		repository.NewSongRepository(
			context.Background(),
			dbpool,
		),
	)

	r := mux.NewRouter()

	r.HandleFunc("/info", songhandlers.Info).
		Methods(http.MethodGet)
	r.HandleFunc("/library/songs", songhandlers.Lib).
		Methods(http.MethodGet)
	r.HandleFunc("/delete", songhandlers.Delete).
		Methods(http.MethodDelete)
	r.HandleFunc("/song", songhandlers.Change).
		Methods(http.MethodPut)
	r.Use(middlewares.LoggingMiddleware)

	log.Println("Listen and serve port 8080")
	http.ListenAndServe(":8080", r)
}
