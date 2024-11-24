package model

import (
	"bytes"
	"database/sql"
	"errors"
	"strconv"
	"time"
)

type Date string

func (d *Date) Scan(value interface{}) error {
	source, ok := value.(time.Time)
	if !ok {
		return errors.New("incompatible type")
	}

	dd, err := d.timeToDate(&source)
	if err != nil {
		return err
	}

	*d = Date(dd)

	return nil
}

func (d *Date) timeToDate(t *time.Time) (string, error) {
	b := new(bytes.Buffer)

	if _, err := b.WriteString(strconv.Itoa(t.Day())); err != nil {
		return "", err
	}

	if _, err := b.WriteRune('.'); err != nil {
		return "", err
	}

	month := strconv.Itoa(int(t.Month()))

	if len(month) == 1 {
		if _, err := b.WriteRune('0'); err != nil {
			return "", err
		}
		if _, err := b.WriteString(month); err != nil {
			return "", err
		}
	} else {
		if _, err := b.WriteString(month); err != nil {
			return "", err
		}
	}

	if _, err := b.WriteRune('.'); err != nil {
		return "", err
	}

	if _, err := b.WriteString(strconv.Itoa(t.Year())); err != nil {
		return "", err
	}

	return b.String(), nil
}

// func (d *Date) UnmarshalJSON(bytes []byte) error {
// 	if string(bytes) == "null" {
// 		return nil
// 	}

// }

type SongDetail struct {
	Song        string         `json:"song"`
	Group       string         `json:"group"`
	ReleaseData Date           `json:"releaseDate,omitempty"`
	Lyric       sql.NullString `json:"text,omitempty"`
	Link        sql.NullString `json:"link,omitempty"`
	Format      sql.NullString `json:"format,omitempty"`
}

type SongModel struct {
	Song        string         `mdb:"song"`
	Group       string         `mdb:"group"`
	ReleaseData Date           `mdb:"releaseDate,omitempty"`
	Lyric       sql.NullString `mdb:"text,omitempty"`
	Link        sql.NullString `mdb:"link,omitempty"`
}
