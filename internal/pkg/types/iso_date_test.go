package types_test

import (
	"encoding/json"
	"testing"
	"time"

	"salary_calculator/internal/pkg/types"

	"github.com/stretchr/testify/assert"
)

func TestISODate_JSON(t *testing.T) {
	t.Run("unmarshal valid", func(t *testing.T) {
		var d types.ISODate
		err := json.Unmarshal([]byte(`"2026-08-03"`), &d)
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2026, 8, 3, 0, 0, 0, 0, time.UTC), d.Time)
	})

	t.Run("unmarshal invalid format", func(t *testing.T) {
		var d types.ISODate
		err := json.Unmarshal([]byte(`"03.08.2026"`), &d)
		assert.Error(t, err)
	})

	t.Run("marshal round-trip", func(t *testing.T) {
		d := types.ISODate{Time: time.Date(2026, 8, 3, 0, 0, 0, 0, time.UTC)}
		b, err := json.Marshal(d)
		assert.NoError(t, err)
		assert.Equal(t, `"2026-08-03"`, string(b))
	})
}

func TestParseISODate(t *testing.T) {
	d, err := types.ParseISODate("2026-01-31")
	assert.NoError(t, err)
	assert.Equal(t, 31, d.Day())

	_, err = types.ParseISODate("not-a-date")
	assert.Error(t, err)
}

func TestISODate_PG(t *testing.T) {
	d, _ := types.ParseISODate("2026-08-03")
	pg := d.PG()
	assert.True(t, pg.Valid)
	assert.Equal(t, d.Time, pg.Time)
}
