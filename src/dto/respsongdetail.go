package dto

import "github.com/lyric/songs/hw/src/repository/model"

// SongDetail
//
// swagger:model songDetail
type SongDetail struct {
	// swagger:strfmt date
	ReleaseDate model.Date `json:"release_date"`
	Lyric       string     `json:"text"`
	Link        string     `json:"link"`
}
