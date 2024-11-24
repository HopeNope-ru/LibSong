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

func ModelSong2Song(mdl *model.SongModel, song *dto.Song) {
	song.Group = mdl.Group
	song.Song = mdl.Song
	song.ReleaseDate = mdl.ReleaseData.Time
	song.Text = mdl.Link.String
	song.Link = mdl.Link.String
}

func UnmarshalSong(r *http.Request) (*dto.Song, error) {
	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		return nil, errors.New(ser)
	}

	var song dto.Song
	if err = json.Unmarshal(b, &song); err != nil {
		ser := fmt.Sprintf(`{"error": "%s"}`, err)
		return nil, errors.New(ser)
	}

	return &song, nil
}

func GenerateUpdateQuery(table string, object any) (string, []any) {
	b := new(bytes.Buffer)
	b.Grow(255)
	b.WriteString("UPDATE ")
	b.WriteString(table)
	b.WriteString(" SET ")

	rs := reflect.ValueOf(object)
	ts := reflect.TypeOf(object)
	k := rs.Kind()
	fmt.Println(k)
	// Создаем и получаем массив с названием полей структуры
	num := rs.NumField()
	keyarr := make([]string, 0, num)
	valarr := make([]any, 0, num)
	for i := 0; i < num; i++ {
		if rs.Field(i).IsZero() {
			continue
		}

		nameMdb := ts.Field(i).Tag.Get("mdb")
		valMdb := rs.Field(i).String() // Минус то что тут постоянно строка, нужно тут типизировать
		keyarr = append(keyarr, nameMdb)
		valarr = append(valarr, valMdb)
	}

	first := true
	for i, val := range keyarr {
		if !first {
			b.WriteString(", ")
		} else {
			first = false
		}

		switch lval := strings.ToLower(val); lval {
		case "group":
			b.WriteRune('"')
			b.WriteString(lval)
			b.WriteRune('"')
			b.WriteString(" = ")
			b.WriteString(fmt.Sprintf("$%v ", i+1))
		default:
			b.WriteString(lval)
			b.WriteString(" = ")
			b.WriteString(fmt.Sprintf("$%v ", i+1))
		}
	}

	// Создаем условие WHERE
	// Можно было бы использовать еще один аргумент, чтобы посмторить условие.
	// Нужно думать как это делать, пока нет времени.
	// b.WriteString(predicate)
	// if predicate != nil {
	// 	b.WriteString(" WHERE ")
	// 	for _, val := range predicate {
	// 		switch r := reflect.TypeOf(val); r.Kind() {
	// 		case reflect.String:
	// 			b.WriteString(r.Name())
	// 			b.WriteString(" = ")
	// 			b.WriteRune('"')
	// 			b.WriteString(r.String())
	// 			b.WriteRune('"')
	// 		default:
	// 			b.WriteString(fmt.Sprintf("%v", predicate))
	// 		}
	// 	}
	// }

	return b.String(), valarr
}
