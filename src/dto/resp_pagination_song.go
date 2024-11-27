package dto

// swagger:model paginationSong
type PaginationSong struct {
	// Группа исполнителей
	Group string `json:"group,omitempty"`

	// Название песни исполнителей
	Song string `json:"song,omitempty"`

	// Дата выхода песни
	ReleaseDate string `json:"releaseDate,omitempty"`

	// Пагинация ложится на фронту, поскольку объем текста песен не большой и можно хранить у клиента
	Text []string `json:"text,omitempty"`

	// Источник песни
	Link string `json:"link,omitempty"`
}
