package dto

// swagger:model song
type Song struct {
	// Группа исполнителей
	Group string `json:"group,omitempty" mdb:"group"`

	// Название песни исполнителей
	Song string `json:"song,omitempty" mdb:"song"`

	// Дата выхода песни
	ReleaseDate string `json:"releaseDate,omitempty" mdb:"release_date"`

	// Первый куплет песни
	Text string `json:"text,omitempty" mdb:"lyrics"`

	// Источник песни
	Link string `json:"link,omitempty" mdb:"link"`
}
