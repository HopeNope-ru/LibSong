package dto

// SongDetail
//
// swagger:model respPaginationLib
type RespPaginationLib struct {
	// Признак, того, что имеется ли еще песни
	Next bool `json:"next"`

	// Массив полученных песен
	Songs []Song `json:"songs"`
}
