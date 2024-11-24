package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
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

	for _, val := range m {
		fmt.Println(val)
	}
}
