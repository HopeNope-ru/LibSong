package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lyric/songs/hw/src/handlers/dto"
	"github.com/lyric/songs/hw/src/repository"
	"github.com/lyric/songs/hw/src/repository/model"
)

func main() {
	// GetApp()

	ctx := context.Background()
	dsn := "postgres://golang:1234@localhost:5432/song"
	dbpool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		panic(err)
	}

	sqlQuery := "select song, \"group\" from library.song limit $2 offset $1"
	row, err := dbpool.Query(ctx, sqlQuery, 2, 2)
	if err != nil {
		panic(err)
	}

	m := make([]model.SongDetail, 0, 2)
	for row.Next() {
		var s model.SongDetail
		err := row.Scan(&s.Song, &s.Group)
		if err != nil {
			panic(err)
		}

		m = append(m, s)
	}

	s := repository.NewSongRepository(context.Background(), dbpool)
	res, err := s.SelectSong("Wham!", "Last Christmas")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	ra, err := s.DeleteSong("Wham!1", "Last Christmas")
	if err != nil {
		panic(err)
	}
	fmt.Printf("rows affected %v\n", ra)

	song := dto.ReqSong{Group: "sdfg", Link: "sdfasdgfg"}
	// rsong := reflect.ValueOf(song)
	// for i := 0; i < rsong.NumField(); i++ {
	// 	fmt.Printf("Field: %d \t type: %T \t value: %v\n", i, rsong.Field(i), rsong.Field(i))
	// 	if rsong.Field(i).IsZero() {
	// 		fmt.Println("Field is nil")
	// 	}
	// 	field := rsong.Type()
	// 	vl := field.Field(i).Tag.Get("mdb")
	// 	valu, ok := field.Field(i).Tag.Lookup("mdb")
	// 	if !ok {
	// 		fmt.Println("NOT OK")
	// 	}
	// 	fmt.Printf("value of tag = %s and val of lookup = %s\n", vl, valu)
	// }

	s.ChangeSong("123", "123", song)
	// fmt.Println(resp)
}
