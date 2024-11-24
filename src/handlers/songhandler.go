package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	w.Header().Add("Content-Type", "application/json")

	v := r.URL.Query()

	group := v.Get("group")
	if !v.Has("group") {
		http.Error(w, `{"error": "not found group"}`, http.StatusBadRequest)
		return
	}

	song := v.Get("song")
	if !v.Has("song") {
		http.Error(w, "not found song", http.StatusBadRequest)
		return
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
		return
	}

	w.Write(b)
}

func (sh *SongHandler) Lib(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	values := r.URL.Query()

	soffset := values.Get("offset")
	if !values.Has("offset") {
		soffset = "0"
	}

	slimit := values.Get("limit")
	if !values.Has("limit") {
		slimit = "5"
	}

	offset, err := strconv.Atoi(soffset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(slimit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filter := values.Get("filter")
	if !values.Has("filter") {
		filter = "song"
	}

	order := values.Get("order")
	if !values.Has("order") {
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

	if part < 0 {
		resp.Songs = songs
	} else {
		resp.Songs = songs[:limit]
	}

	b, err := json.Marshal(&resp)
	if err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func (sh *SongHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

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

	if ok, err := ValidReq(req); !ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	w.Header().Add("Content-Type", "application/json")

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

	if ok, err := ValidReq(req); !ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	w.Header().Add("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		http.Error(w, "request must be DELETE", http.StatusBadRequest)
		return
	}

	req, err := utils.UnmarshalSong(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ok, err := ValidReq(req); !ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

func ValidReq(req *dto.Song) (bool, error) {
	if len(req.Group) == 0 {
		return false, errors.New("need group")
	}

	if len(req.Song) == 0 {
		return false, errors.New("need song")
	}

	return true, nil
}
