package enums

// swagger:enum Filter
type Filter string

const (
	GROUP        Filter = "group"
	SONG         Filter = "song"
	RELEASE_DATE Filter = "release_date"
	LYRIC        Filter = "text"
	LINK         Filter = "link"
)
