package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lyric/songs/hw/src/handlers/dto"
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

func (s *SongRepository) ChangeSong(group, song string, req dto.ReqSong) (int64, error) {

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

// func (s *SongRepository) execUpdate(query string, song dto.ReqSong) {

// }
