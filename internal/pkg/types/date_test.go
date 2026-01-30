package types

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDate_JSON(t *testing.T) {
	t.Run("Unmarshal success", func(t *testing.T) {
		var d Date
		err := json.Unmarshal([]byte(`"15.01.2024"`), &d)
		assert.NoError(t, err)
		assert.Equal(t, 2024, d.Year())
		assert.Equal(t, time.January, d.Month())
		assert.Equal(t, 15, d.Day())
	})

	t.Run("Unmarshal failure", func(t *testing.T) {
		var d Date
		err := json.Unmarshal([]byte(`"2024-01-15"`), &d)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid date format")

		err = json.Unmarshal([]byte(`123`), &d)
		assert.Error(t, err)
	})

	t.Run("Marshal", func(t *testing.T) {
		d := Date{Time: time.Date(2024, time.May, 20, 0, 0, 0, 0, time.UTC)}
		b, err := json.Marshal(&d)
		assert.NoError(t, err)
		assert.Equal(t, `"20.05.2024"`, string(b))
	})
}
