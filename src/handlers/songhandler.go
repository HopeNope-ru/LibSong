package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/lyric/songs/hw/src/dto"
	"github.com/lyric/songs/hw/src/enums"
	"github.com/lyric/songs/hw/src/repository"
	"github.com/lyric/songs/hw/src/utils"
	"github.com/rs/zerolog/log"
)

// swagger:parameters librarySong
type QuerySongLib struct {
	// in query
	// required: false
	// default: 0
	Offset int `json:"offset"`

	// in: query
	// required: false
	// default: 5
	Limit int `json:"limit"`

	// OrderBy
	// default: asc
	OrderBy enums.OrderBy `json:"order_by"`

	// Filter
	// default: song
	Filter enums.Filter `json:"filter"`
}

// swagger:parameters createSong deleteSong changeSong info
type QueryCreate struct {
	// in: query
	// required: true
	Song string `json:"song"`

	// in: query
	// required: true
	Group string `json:"group"`
}

// swagger:parameters changeSong
type SongBodyParams struct {

	// required: true
	// in: body
	Body *dto.Song `json:"song"`
}

// swagger:response respSongDetail
type ResponseSongDetail struct {
	Payload *dto.SongDetail `json:"song"`
}

// swagger:response respPaginationSong
type ResponsePaginationSong struct {
	Body *dto.PaginationSong `json:"song"`
}

// Модель ошибки
//
// swagger:response respError
type ErrorResponse struct {
	// in: body
	Payload *dto.ResponseError `json:"error"`
}

type SongHandler struct {
	storage *repository.SongRepository
}

func NewSongHandler(storage *repository.SongRepository) *SongHandler {
	return &SongHandler{storage: storage}
}

// Info swagger:route GET /info song info
//
// # Получение информации по песне
//
// Responses:
//
//	200: respSongDetail
//	400: respError
//	500: respError
func (sh *SongHandler) Info(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	v := r.URL.Query()

	group := v.Get("group")
	if err := utils.ValidQuery("group", &v); err != nil {
		rerr := dto.NewRespError(http.StatusBadRequest, err.Error(), "")
		utils.JSONError(w, rerr, http.StatusBadRequest)
		return
	}

	song := v.Get("song")
	if err := utils.ValidQuery("song", &v); err != nil {
		rerr := dto.NewRespError(http.StatusBadRequest, err.Error(), "")
		utils.JSONError(w, rerr, http.StatusBadRequest)
		return
	}

	s, err := sh.storage.SelectSongDetail(group, song)
	if err != nil {
		log.Err(err).Send()
		rerr := dto.NewRespError(http.StatusInternalServerError, err.Error(), "Ошибка в базе данных")
		utils.JSONError(w, rerr, http.StatusInternalServerError)
	}

	b, err := json.Marshal(&s)
	if err != nil {
		log.Err(err).Send()
		rerr := dto.NewRespError(http.StatusInternalServerError, err.Error(), "couldn't marshal info")
		utils.JSONError(w, rerr, http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

// PaginationSong swagger:route GET /song song info
//
// # Получение песни с пагинацией
//
// Responses:
//
//	200: respPaginationSong
//	400: respError
//	500: respError
func (sh *SongHandler) PagninationSong(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	v := r.URL.Query()

	group := v.Get("group")
	if err := utils.ValidQuery("group", &v); err != nil {
		rerr := dto.NewRespError(http.StatusBadRequest, err.Error(), "")
		utils.JSONError(w, rerr, http.StatusBadRequest)
		return
	}

	song := v.Get("song")
	if err := utils.ValidQuery("song", &v); err != nil {
		rerr := dto.NewRespError(http.StatusBadRequest, err.Error(), "")
		utils.JSONError(w, rerr, http.StatusBadRequest)
		return
	}

	s, err := sh.storage.SelectPaginationSong(group, song)
	if err != nil {
		log.Err(err).Send()
		rerr := dto.NewRespError(http.StatusInternalServerError, err.Error(), "Ошибка в базе данных")
		utils.JSONError(w, rerr, http.StatusInternalServerError)
	}

	b, err := json.Marshal(&s)
	if err != nil {
		log.Err(err).Send()
		rerr := dto.NewRespError(http.StatusInternalServerError, err.Error(), "couldn't marshal info")
		utils.JSONError(w, rerr, http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

// LibrarySong swagger:route GET /library/songs song librarySong
//
// # Получение песен с пагинацией
//
// Responses:
//
//	200: respPaginationLib
//	400: respError
//	500: respError
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
		rerr := dto.NewRespError(http.StatusBadRequest, err.Error(), "Используйте число а не строку")
		utils.JSONError(w, rerr, http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(slimit)
	if err != nil {
		rerr := dto.NewRespError(http.StatusBadRequest, err.Error(), "Используйте число а не строку")
		utils.JSONError(w, rerr, http.StatusBadRequest)
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
		ser := "order is not compatible"
		log.Error().Str("error", ser).Send()
		rerr := dto.NewRespError(http.StatusBadRequest, ser, "")
		utils.JSONError(w, rerr, http.StatusBadRequest)
		return
	}

	// Заглядываем в будущее запроса, чтобы понять есть ли у нас еще строки в БД
	future := 5
	modelsongs, err := sh.storage.SelectFuturePaginationLibSong(offset, limit, future, filter, order)
	if err != nil {
		log.Error().Stack().Err(err).Send()
		rerr := dto.NewRespError(http.StatusInternalServerError, err.Error(), "Ошибка в базе данных")
		utils.JSONError(w, rerr, http.StatusInternalServerError)
		return
	}

	// Перемещаем из модели бд в ответ
	songs := make([]dto.Song, 0, len(modelsongs))
	for _, sng := range modelsongs {
		var ss dto.Song
		utils.ModelSong2SongToVerseList(&sng, &ss)
		songs = append(songs, ss)
	}

	resp := dto.RespPaginationLib{Next: false}

	part := len(songs) - limit
	if part > 0 {
		resp.Next = true
	}

	resp.Songs = songs[:limit]

	b, err := json.Marshal(&resp)
	if err != nil {
		rerr := dto.NewRespError(http.StatusInternalServerError, err.Error(), "Ошибка маршалирования")
		utils.JSONError(w, rerr, http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

// CreateSong swagger:route POST /create song createSong
//
// # Добавление песни
//
// Responses:
//
//	201:
//	400: respError
//	500: respError
func (sh *SongHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	m := r.Method
	if m != http.MethodPost {
		rerr := dto.NewRespError(http.StatusBadRequest, "request must be POST", "")
		utils.JSONError(w, rerr, http.StatusBadRequest)
		return
	}

	req, err := utils.UnmarshalSong(r)
	if err != nil {
		log.Error().Stack().Err(err).Send()
		rerr := dto.NewRespError(http.StatusInternalServerError, err.Error(), "Ошибка маршалирования")
		utils.JSONError(w, rerr, http.StatusInternalServerError)
		return
	}

	if ok, err := ValidReq(req); !ok {
		log.Error().Stack().Err(err).Send()
		rerr := dto.NewRespError(http.StatusBadRequest, err.Error(), "")
		utils.JSONError(w, rerr, http.StatusBadRequest)
		return
	}

	if err = sh.storage.CreateSong(req.Group, req.Song); err != nil {
		log.Error().Stack().Err(err).Send()
		rerr := dto.NewRespError(http.StatusInternalServerError, err.Error(), "Ошибка в базе данных")
		utils.JSONError(w, rerr, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// ChangeSong swagger:route PUT /change song changeSong
//
// # Изменение песни
//
// Responses:
//
//	201:
//	304:
//	400: respError
//	500: respError
func (sh *SongHandler) Change(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		rerr := dto.NewRespError(http.StatusBadRequest, "request must be PUT", "")
		utils.JSONError(w, rerr, http.StatusBadRequest)
		return
	}

	req, err := utils.UnmarshalSong(r)
	if err != nil {
		log.Error().Stack().Err(err).Send()
		rerr := dto.NewRespError(http.StatusInternalServerError, err.Error(), "Ошибка маршалирования")
		utils.JSONError(w, rerr, http.StatusInternalServerError)
		return
	}

	if ok, err := ValidReq(req); !ok {
		rerr := dto.NewRespError(http.StatusBadRequest, err.Error(), "")
		utils.JSONError(w, rerr, http.StatusBadRequest)
		return
	}

	rowAffected, err := sh.storage.ChangeSong(req.Group, req.Song, req)
	if err != nil {
		log.Error().Stack().Err(err).Send()
		rerr := dto.NewRespError(http.StatusInternalServerError, err.Error(), "Ошибка в базе данных")
		utils.JSONError(w, rerr, http.StatusInternalServerError)
		return
	}

	if rowAffected == 0 {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// DeleteSong swagger:route DELETE /delete song deleteSong
//
// # Удаление песни
//
// Responses:
//
//	200:
//	304:
//	500: respError
func (sh *SongHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		rerr := dto.NewRespError(http.StatusBadRequest, "request must be DELETE", "")
		utils.JSONError(w, rerr, http.StatusBadRequest)
		return
	}

	req, err := utils.UnmarshalSong(r)
	if err != nil {
		log.Error().Stack().Err(err).Send()
		rerr := dto.NewRespError(http.StatusInternalServerError, err.Error(), "Ошибка маршалирования")
		utils.JSONError(w, rerr, http.StatusInternalServerError)
		return
	}

	if ok, err := ValidReq(req); !ok {
		rerr := dto.NewRespError(http.StatusBadRequest, err.Error(), "")
		utils.JSONError(w, rerr, http.StatusBadRequest)
		return
	}

	rowAffected, err := sh.storage.DeleteSong(req.Group, req.Song)
	if err != nil {
		log.Error().Stack().Err(err).Send()
		rerr := dto.NewRespError(http.StatusInternalServerError, err.Error(), "Ошибка в базе данных")
		utils.JSONError(w, rerr, http.StatusInternalServerError)
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
