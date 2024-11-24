package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepSong interface {
}

type SongHandler struct {
	db *pgxpool.Pool
}

func Info(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	group, ok := v["group"]
	if !ok {
		http.Error(w, `{"error": "not found group"}`, http.StatusBadRequest)
	}

	song, ok := v["song"]
	if !ok {
		http.Error(w, `{"error": "not found group"}`, http.StatusBadRequest)
	}

	s, err := storage.GetSong(group, song)
	// Реализовать логирование ошибки

	b, err := json.Marshal(&s)
	if err != nil {
		http.Error(w, `{"error": "couldn't marshal info"}`, http.StatusInternalServerError)
	}

	w.Write(b)
}

func Lib(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	offset, ok := v["offset"]
	if !ok {
		offset = "0"
	}

	filter, ok := v["filter"]
	if !ok {
		filter = "song"
	}

	order, ok := v["filter"]
	if !ok {
		order = "asc"
	}

	s, err := storage.GetLibSongs(offset, filter, order)
	if err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusBadRequest)
	}

}

func Delete(w http.ResponseWriter, r *http.Request) {

}
