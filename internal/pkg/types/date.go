package types

import (
	"encoding/json"
	"fmt"
	"time"
)

type Date struct {
	time.Time
}

func (cd *Date) UnmarshalJSON(b []byte) error {
	var dateStr string
	if err := json.Unmarshal(b, &dateStr); err != nil {
		return err
	}

	t, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format %s: %v", dateStr, err)
	}

	cd.Time = t
	return nil
}

func (cd *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(cd.Time.Format("02.01.2006"))
}
