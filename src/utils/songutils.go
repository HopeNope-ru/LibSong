package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lyric/songs/hw/src/repository/model"
)

func ToVerseList(text string) []string {
	return strings.Split(text, "\\n\\n")
}

func UnmarshalSong(r *http.Request) (*model.SongDetail, error) {
	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		return nil, errors.New(ser)
	}

	var song model.SongDetail
	if err = json.Unmarshal(b, &song); err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		return nil, errors.New(ser)
	}

	return &song, nil
}
