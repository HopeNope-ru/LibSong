package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lyric/songs/hw/src/handlers/dto"
	"github.com/lyric/songs/hw/src/repository"
	"github.com/lyric/songs/hw/src/utils"
)

type SongHandler struct {
	storage *repository.SongRepository
}

func NewSongHandler(storage *repository.SongRepository) *SongHandler {
	return &SongHandler{storage: storage}
}

func (sh *SongHandler) Info(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	group, ok := v["group"]
	if !ok {
		http.Error(w, `{"error": "not found group"}`, http.StatusBadRequest)
	}

	song, ok := v["song"]
	if !ok {
		http.Error(w, `{"error": "not found group"}`, http.StatusBadRequest)
	}

	s, err := sh.storage.SelectSong(group, song)
	if err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusInternalServerError)
	}
	// Реализовать логирование ошибки

	b, err := json.Marshal(&s)
	if err != nil {
		http.Error(w, `{"error": "couldn't marshal info"}`, http.StatusInternalServerError)
	}

	w.Write(b)
}

func (sh *SongHandler) Lib(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	soffset, ok := v["offset"]
	if !ok {
		soffset = "0"
	}

	slimit, ok := v["limit"]
	if !ok {
		slimit = "5"
	}

	limit, err := strconv.Atoi(soffset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offset, err := strconv.Atoi(slimit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filter, ok := v["filter"]
	if !ok {
		filter = "song"
	}

	order, ok := v["order"]
	if !ok {
		order = "asc"
	}

	switch lo := strings.ToLower(order); lo {
	case "asc":
		order = lo
	case "desc":
		order = lo
	default:
		ser := `{"error": "order is not compatible"}`
		http.Error(w, ser, http.StatusBadRequest)
		return
	}

	// Заглядываем в будущее запроса, чтобы понять есть ли у нас еще строки в БД
	future := 5
	songs, err := sh.storage.SelectFuturePaginationLibSong(offset, limit, future, filter, order)
	if err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusBadRequest)
		return
	}

	resp := dto.RespPaginationLib{Next: false}

	part := len(songs) - limit
	if part > 0 {
		resp.Next = true
	}

	resp.Songs = songs

	b, err := json.Marshal(&resp)
	if err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func (sh *SongHandler) Create(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodPost {
		ser := `{"error": "request must be POST"}`
		http.Error(w, ser, http.StatusBadRequest)
		return
	}

	req, err := utils.UnmarshalSong(r)
	if err != nil {
		ser := `{"error": "request must be POST"}`
		http.Error(w, ser, http.StatusBadRequest)
		return
	}

	if len(req.Group) == 0 {
		http.Error(w, "need group", http.StatusBadRequest)
		return
	}

	if len(req.Song) == 0 {
		http.Error(w, "need song", http.StatusBadRequest)
		return
	}

	if err = sh.storage.CreateSong(req.Group, req.Song); err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (sh *SongHandler) Change(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		ser := `{"error": "request must be PUT"}`
		http.Error(w, ser, http.StatusBadRequest)
		return
	}

	req, err := utils.UnmarshalSong(r)
	if err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusInternalServerError)
		return
	}

	if len(req.Group) == 0 {
		http.Error(w, "need group", http.StatusBadRequest)
		return
	}

	if len(req.Song) == 0 {
		http.Error(w, "need song", http.StatusBadRequest)
		return
	}

	rowAffected, err := sh.storage.ChangeSong(req.Group, req.Song, req)
	if err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusInternalServerError)
		return
	}

	if rowAffected == 0 {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (sh *SongHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "request must be DELETE", http.StatusBadRequest)
		return
	}

	req, err := utils.UnmarshalSong(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ValidReq(w, req)

	rowAffected, err := sh.storage.DeleteSong(req.Group, req.Song)
	if err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusInternalServerError)
		return
	}

	if rowAffected == 0 {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ValidReq(w http.ResponseWriter, req *dto.Song) {
	if len(req.Group) == 0 {
		http.Error(w, "need group", http.StatusBadRequest)
	}

	if len(req.Song) == 0 {
		http.Error(w, "need song", http.StatusBadRequest)
	}
}
