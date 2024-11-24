package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lyric/songs/hw/src/repository/model"
	"github.com/lyric/songs/hw/src/storage"
	"github.com/lyric/songs/hw/src/utils"
)

func stub(w http.ResponseWriter, r *http.Request) {

}

func GetApp(args ...string) {
	dsn := "postgres://golang:1234@localhost:5432/song"
	db, err := storage.NewDbPool(dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	group := "Wham!"
	song := "Last Christmas"

	s := model.SongDetail{}
	ctx := context.Background()
	err = db.QueryRow(ctx, "select release_date, lyrics, link from library.song where \"group\" = $1 and song = $2", group, song).
		Scan(&s.ReleaseData, &s.Lyric, &s.Link)

	if err != nil {
		panic(err)
	}

	txt, _ := s.Lyric.Value()
	if txt == nil {
		panic("not value in txt db")
	}
	pag := utils.ToVerseList(txt.(string))
	for i, val := range pag {
		fmt.Printf("%v. %s\n", i, val)
	}

	r := mux.NewRouter()

	r.HandleFunc("/info", stub).
		Methods(http.MethodGet)
	r.HandleFunc("/library/songs", stub).
		Methods(http.MethodGet)
	r.HandleFunc("/text/song", stub).
		Methods(http.MethodPost)

}
