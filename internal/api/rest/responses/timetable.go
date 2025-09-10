package responses

import (
	"fmt"

	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/resources"
)

func Timetable(m models.Timetable) resources.Timetable {
	intervals := make([]resources.TimetableInterval, 0, len(m.Table))
	for _, i := range m.Table {
		intervals = append(intervals, TimetableInterval(i))

	}

	return resources.Timetable{
		Data: resources.TimetableData{
			Type: resources.TimetableType,
			Attributes: resources.TimetableDataAttributes{
				Table: intervals,
			},
		},
	}
}

func TimetableInterval(i models.TimeInterval) resources.TimetableInterval {
	return resources.TimetableInterval{
		From: resources.TimeMoment{
			Weekday: i.From.Weekday.String(),
			Time:    fmt.Sprintf("%v:%v", i.From.Time.Hours(), i.From.Time.Minutes()),
		},
		To: resources.TimeMoment{
			Weekday: i.To.Weekday.String(),
			Time:    fmt.Sprintf("%v:%v", i.From.Time.Hours(), i.From.Time.Minutes()),
		},
	}
}
