package models

import "time"

type Moment struct {
	Weekday time.Weekday
	Time    time.Duration
}

type TimeInterval struct {
	From Moment
	To   Moment
}
type Timetable struct {
	Intervals []TimeInterval `json:"intervals"`
}

func NumberMinutesToMoment(minutes int) Moment {
	weekday := time.Weekday(minutes / (60 * 24))
	hour := (minutes % (60 * 24)) / 60

	return Moment{
		Weekday: weekday,
		Time:    time.Duration(hour) * time.Hour,
	}
}

func (ti TimeInterval) ToNumberMinutes() (int, int) {
	from := int(ti.From.Weekday)*24*60 + int(ti.From.Time.Hours())*60 + int(ti.From.Time.Minutes())%60
	to := int(ti.To.Weekday)*24*60 + int(ti.To.Time.Hours())*60 + int(ti.To.Time.Minutes())%60

	return from, to
}
