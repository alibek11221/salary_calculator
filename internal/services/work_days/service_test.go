package work_days

import (
	"salary_calculator/internal/pkg/http/work_calendar"
	"salary_calculator/internal/pkg/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestService_CalculateWorkDaysForMonth(t *testing.T) {
	service := New()

	tests := []struct {
		name  string
		month work_calendar.WorkdayResponse
		want  WorkdaysForMonth
		isNil bool
	}{
		{
			name: "Regular month",
			month: work_calendar.WorkdayResponse{
				Statistics: work_calendar.Statistics{WorkDays: 20},
				Days: []work_calendar.Day{
					{Date: types.Date{Time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}, TypeId: 1},
					{Date: types.Date{Time: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)}, TypeId: 1},
					{Date: types.Date{Time: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)}, TypeId: 1},
					{Date: types.Date{Time: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)}, TypeId: 1},
					{Date: types.Date{Time: time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC)}, TypeId: 1},
					{Date: types.Date{Time: time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)}, TypeId: 1},
					{Date: types.Date{Time: time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC)}, TypeId: 2},
				},
			},
			want: WorkdaysForMonth{
				TotalWorkdays:  20,
				FirstHalfDays:  4,
				SecondHalfDays: 2,
			},
		},
		{
			name: "Short day is work day",
			month: work_calendar.WorkdayResponse{
				Statistics: work_calendar.Statistics{WorkDays: 20},
				Days: []work_calendar.Day{
					{Date: types.Date{Time: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)}, TypeId: 5},
					{Date: types.Date{Time: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC)}, TypeId: 5},
				},
			},
			want: WorkdaysForMonth{
				TotalWorkdays:  20,
				FirstHalfDays:  1,
				SecondHalfDays: 1,
			},
		},
		{
			name:  "Nil response",
			month: work_calendar.WorkdayResponse{},
			want:  WorkdaysForMonth{},
			isNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var arg *work_calendar.WorkdayResponse
			if !tt.isNil {
				arg = &tt.month
			}
			got := service.CalculateWorkDaysForMonth(arg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_CalculateWorkDays(t *testing.T) {
	service := New()

	months := map[int]*work_calendar.WorkdayResponse{
		1: {
			Statistics: work_calendar.Statistics{WorkDays: 20},
			Days: []work_calendar.Day{
				{Date: types.Date{Time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}, TypeId: 1},
			},
		},
	}

	got := service.CalculateWorkDays(months)
	assert.Len(t, got, 1)
	assert.Equal(t, 20, got[1].TotalWorkdays)
	assert.Equal(t, 1, got[1].FirstHalfDays)
}
