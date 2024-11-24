package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lyric/songs/hw/src/repository/model"
)

type SongRepository struct {
	table string
	ctx   context.Context
	db    *pgxpool.Pool
}

func NewSongRepository(ctx context.Context, db *pgxpool.Pool) *SongRepository {
	return &SongRepository{db: db, ctx: ctx, table: "library.song"}
}

func (s *SongRepository) SelectSong(group, song string) (*model.SongDetail, error) {

	var sng model.SongDetail
	q := fmt.Sprintf(`
		select release_date, lyrics, link 
		from %s 
		where 1=1
			and "group" = $1 
			and song = $2`, s.table,
	)

	err := s.db.QueryRow(s.ctx, q, group, song).
		Scan(&sng.ReleaseData, &sng.Lyric, &sng.Link)

	if err != nil {
		return nil, err
	}

	return &sng, nil
}

func (s *SongRepository) DeleteSong(group, song string) (int64, error) {
	q := fmt.Sprintf("delete from %s where song = $1 and \"group\" = $2", s.table)
	ct, err := s.db.Exec(s.ctx, q, song, group)
	if err != nil {
		return 0, err
	}

	ra := ct.RowsAffected()
	return ra, nil
}

// func (s *SongRepository) ChangeSong(song dto.ReqSong) (int64, error) {

// 	q := utils.GenerateUpdateQuery(s.table, song)

// 	return nil
// }

// func (s *SongRepository) execUpdate(query string, song dto.ReqSong) {

// }
