package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const isoDateLayout = "2006-01-02"

// ISODate — дата без времени, в JSON в формате "2006-01-02".
type ISODate struct {
	time.Time
}

func ParseISODate(s string) (ISODate, error) {
	t, err := time.Parse(isoDateLayout, s)
	if err != nil {
		return ISODate{}, fmt.Errorf("invalid date %q: expected format YYYY-MM-DD", s)
	}
	return ISODate{Time: t}, nil
}

func (d *ISODate) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	parsed, err := ParseISODate(s)
	if err != nil {
		return err
	}
	d.Time = parsed.Time
	return nil
}

func (d ISODate) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format(isoDateLayout))
}

// PG конвертирует дату в pgtype.Date для sqlc-параметров.
func (d ISODate) PG() pgtype.Date {
	return pgtype.Date{Time: d.Time, Valid: true}
}
