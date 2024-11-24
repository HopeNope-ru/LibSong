package main

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lyric/songs/hw/src/handlers"
	"github.com/lyric/songs/hw/src/repository"
)

func stub(w http.ResponseWriter, r *http.Request) {

}

func GetApp(args ...string) {
	dsn := "postgres://golang:1234@localhost:5432/song"
	dbpool, err := pgxpool.New(context.Background(), dsn)
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
		Methods(http.MethodDelete).
		Queries(
			"group", "{group}",
			"song", "{song}",
		)

	http.ListenAndServe(":8080", r)
}
