package value_objects

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSalaryDate(t *testing.T) {
	tests := []struct {
		name    string
		date    string
		wantErr bool
	}{
		{"Valid date", "2024_01", false},
		{"Valid date late year", "2025_12", false},
		{"Invalid format", "2024-01", true},
		{"Invalid year (too early)", "2023_01", true},
		{"Invalid year (not a number)", "abcd_01", true},
		{"Invalid month (too high)", "2024_13", true},
		{"Invalid month (too low)", "2024_00", true},
		{"Invalid month (not a number)", "2024_ab", true},
		{"Empty string", "", true},
		{"Missing month", "2024", true},
		{"Too many parts", "2024_01_01", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSalaryDate(tt.date)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.date, got.String())
			}
		})
	}
}

func TestSalaryDate_Compare(t *testing.T) {
	d202401, _ := NewSalaryDate("2024_01")
	d202402, _ := NewSalaryDate("2024_02")
	d202501, _ := NewSalaryDate("2025_01")
	d202501Dup, _ := NewSalaryDate("2025_01")

	tests := []struct {
		name  string
		s     *SalaryDate
		other *SalaryDate
		want  int
	}{
		{"Same date", d202501, d202501Dup, 0},
		{"Earlier month", d202401, d202402, -1},
		{"Later month", d202402, d202401, 1},
		{"Earlier year", d202401, d202501, -1},
		{"Later year", d202501, d202401, 1},
		{"Nil s", nil, d202401, -1},
		{"Nil other", d202401, nil, 1},
		{"Both nil", nil, nil, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.s.Compare(tt.other))
		})
	}
}

func TestSalaryDate_JSON(t *testing.T) {
	d, _ := NewSalaryDate("2024_05")

	data, err := json.Marshal(d)
	assert.NoError(t, err)
	assert.Equal(t, `"2024_05"`, string(data))

	var d2 SalaryDate
	err = json.Unmarshal(data, &d2)
	assert.NoError(t, err)
	assert.Equal(t, d.String(), d2.String())
	assert.Equal(t, d.month, d2.month)
	assert.Equal(t, d.year, d2.year)

	err = json.Unmarshal([]byte(`"2023_01"`), &d2)
	assert.Error(t, err)
}
