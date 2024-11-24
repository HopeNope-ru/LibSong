package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/lyric/songs/hw/src/handlers/dto"
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

func GenerateUpdateQuery(table string, song dto.ReqSong) string {
	b := new(bytes.Buffer)
	b.Grow(255)
	b.WriteString("UPDATE ")
	b.WriteString(table)
	b.WriteString(" SET ")

	rs := reflect.ValueOf(song)
	ts := reflect.TypeOf(song)

	num := rs.NumField()
	keyarr := make([]string, 0, num)
	for i := 0; i < num; i++ {
		if rs.Field(i).IsZero() {
			continue
		}

		nameMdb := ts.Field(i).Tag.Get("mdb")
		keyarr = append(keyarr, nameMdb)
	}

	first := true
	for i, val := range keyarr {
		if !first {
			b.WriteString(", ")
		} else {
			first = false
		}

		b.WriteString(val)
		b.WriteString(" = ")
		b.WriteString(fmt.Sprintf("$%v ", i+1))
	}

	return b.String()
}
