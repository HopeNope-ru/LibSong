package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lyric/songs/hw/src/repository/model"
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

	order, ok := v["order"]
	if !ok {
		order = "asc"
	}

	s, err := storage.GetLibSongs(offset, filter, order)
	if err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusBadRequest)
	}

	b, err := json.Marshal(&s)
	if err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusInternalServerError)
	}

	w.Write(b)
}

func Create(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodPost {
		ser := `{"error": "request must be POST"}`
		http.Error(w, ser, http.StatusBadRequest)
	}

	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusInternalServerError)
	}

	var req model.SongDetail
	if err = json.Unmarshal(b, &resp); err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusInternalServerError)
	}

	if err = storage.CreateSong(req); err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}

func Change(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		ser := `{"error": "request must be PUT"}`
		http.Error(w, ser, http.StatusBadRequest)
	}

	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusInternalServerError)
	}

	var song model.SongDetail
	if err = json.Unmarshal(b, &song); err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusInternalServerError)
	}

	if err := storage.ChangeSong(); err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		ser := `{"error": "request must be DELETE"}`
		http.Error(w, ser, http.StatusBadRequest)
	}

	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusInternalServerError)
	}

	var song model.SongDetail
	if err = json.Unmarshal(b, &song); err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusInternalServerError)
	}

	if err = storage.DeleteSong(); err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
