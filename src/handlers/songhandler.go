package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lyric/songs/hw/src/handlers/dto"
	"github.com/lyric/songs/hw/src/repository"
	"github.com/lyric/songs/hw/src/repository/model"
	"github.com/lyric/songs/hw/src/utils"
)

type RepSong interface {
}

type SongHandler struct {
	db      *pgxpool.Pool
	storage *repository.SongRepository
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
	}

	offset, err := strconv.Atoi(slimit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	}

	// Заглядываем в будущее запроса, чтобы понять есть ли у нас еще строки в БД
	future := 5
	songs, err := sh.storage.SelectFuturePaginationLibSong(offset, limit, future, filter, order)
	if err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusBadRequest)
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
	}

	w.Write(b)
}

func (sh *SongHandler) Create(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodPost {
		ser := `{"error": "request must be POST"}`
		http.Error(w, ser, http.StatusBadRequest)
	}

	req, err := utils.UnmarshalSong(r)
	if err != nil {
		ser := `{"error": "request must be POST"}`
		http.Error(w, ser, http.StatusBadRequest)
	}

	if err = storage.CreateSong(req); err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}

func (sh *SongHandler) Change(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		ser := `{"error": "request must be PUT"}`
		http.Error(w, ser, http.StatusBadRequest)
	}

	// !!!!!!!! Сюда добавить util.UnmarshalSong()
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
	// !!!!!!!!!!!!!!!!!!!!!

	if err := storage.ChangeSong(); err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}

func (sh *SongHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		ser := `{"error": "request must be DELETE"}`
		http.Error(w, ser, http.StatusBadRequest)
	}
	// !!!!!!!! Сюда добавить util.UnmarshalSong()
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
	// !!!!!!!!!!!!!!!!!!!!!

	if err = storage.DeleteSong(); err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		http.Error(w, ser, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
