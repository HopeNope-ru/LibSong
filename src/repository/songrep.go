package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lyric/songs/hw/src/dto"
	"github.com/lyric/songs/hw/src/repository/model"
	"github.com/lyric/songs/hw/src/utils"
)

type SongRepository struct {
	table string
	ctx   context.Context
	db    *pgxpool.Pool
}

func NewSongRepository(ctx context.Context, db *pgxpool.Pool) *SongRepository {
	return &SongRepository{db: db, ctx: ctx, table: "library.song"}
}

func (s *SongRepository) SelectSongDetail(group, song string) (*dto.SongDetail, error) {

	var sng model.SongModel
	q := fmt.Sprintf(`
		select release_date, lyrics, link 
		from %s 
		where 1=1
			and "group" = $1 
			and song = $2`, s.table,
	)

	err := s.db.QueryRow(s.ctx, q, group, song).
		Scan(&sng.ReleaseDate, &sng.Lyric, &sng.Link)

	if err != nil {
		return nil, err
	}

	return &dto.SongDetail{
		Lyric:       sng.Lyric.String,
		Link:        sng.Link.String,
		ReleaseDate: sng.ReleaseDate,
	}, nil
}

func (s *SongRepository) SelectPaginationSong(group, song string) (*dto.PaginationSong, error) {

	var sng model.SongModel
	q := fmt.Sprintf(`
		select release_date, lyrics, link, "group", song
		from %s 
		where 1=1
			and "group" = $1 
			and song = $2`, s.table,
	)

	err := s.db.QueryRow(s.ctx, q, group, song).
		Scan(&sng.ReleaseDate, &sng.Lyric, &sng.Link, &sng.Group, &sng.Song)

	if err != nil {
		return nil, err
	}

	return &dto.PaginationSong{
		Text:        utils.ToVerseList(sng.Lyric.String),
		Link:        sng.Link.String,
		ReleaseDate: sng.ReleaseDate.Time,
		Group:       sng.Group,
		Song:        sng.Song,
	}, nil
}

func (s *SongRepository) DeleteSong(group, song string) (int64, error) {
	q := fmt.Sprintf("delete from %s where song = $1 and \"group\" = $2", s.table)
	ct, err := s.db.Exec(s.ctx, q, song, group)
	if err != nil {
		return 0, err
	}

	return ct.RowsAffected(), nil
}

func (s *SongRepository) ChangeSong(group, song string, req *dto.Song) (int64, error) {

	q, args := utils.GenerateUpdateQuery(s.table, req)

	l := len(args)
	// Большой костыль, от которого надо избавляться
	q = fmt.Sprintf("%s WHERE \"group\" = $%v and song = $%v", q, l+1, l+2)
	args = append(args, group)
	args = append(args, song)

	excd, err := s.db.Exec(s.ctx, q, args...)
	if err != nil {
		return 0, err
	}

	return excd.RowsAffected(), nil
}

func (s *SongRepository) CreateSong(group, song string) error {
	q := fmt.Sprintf(`insert into %s ("group", song) values ($1, $2)`, s.table)
	_, err := s.db.Exec(s.ctx, q, group, song)
	if err != nil {
		return err
	}

	return nil
}

func (s *SongRepository) SelectFuturePaginationLibSong(offset, limit, future int, filter, order_by string) ([]model.SongModel, error) {
	q := fmt.Sprintf("select song, \"group\", release_date, lyrics, link from library.song order by %s %s limit $2 offset $1", filter, order_by)
	row, err := s.db.Query(s.ctx, q, offset, limit+future)
	if err != nil {
		return nil, err
	}

	songs := make([]model.SongModel, 0, limit)
	for row.Next() {

		var ms model.SongModel
		err := row.Scan(&ms.Song, &ms.Group, &ms.ReleaseDate, &ms.Lyric, &ms.Link)
		if err != nil {
			panic(err)
		}

		songs = append(songs, ms)
	}

	return songs, nil
}
