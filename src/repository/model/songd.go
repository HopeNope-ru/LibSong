package model

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"errors"
	"strconv"
	"time"
)

type Date struct {
	Valid bool
	Time  string
}

func (d *Date) Scan(value interface{}) error {
	if value == nil {
		d.Valid, d.Time = false, ""
		return nil
	}

	source, ok := value.(time.Time)
	if !ok {
		return errors.New("incompatible type")
	}

	dd, err := d.timeToDate(&source)
	if err != nil {
		return err
	}

	d.Time = dd

	return nil
}

func (d Date) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Time, nil
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

type SongModel struct {
	Group       string
	Song        string
	ReleaseDate Date
	Lyric       sql.NullString
	Link        sql.NullString
}
