package dto

type ReqSong struct {
	Group       string `json:"group,omitempty" mdb:"group"`
	Song        string `json:"song,omitempty" mdb:"song"`
	ReleaseDate string `json:"releaseDate,omitempty" mdb:"release_date"`
	Text        string `json:"text,omitempty" mdb:"lyrics"`
	Link        string `json:"link,omitempty" mdb:"link"`
	Som         int    `mdb:"som"`
}
