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
	Table []TimeInterval
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
	from := ti.From.ToNumberMinutes()
	to := ti.To.ToNumberMinutes()

	return from, to
}

func (m Moment) ToNumberMinutes() int {
	return int(m.Weekday)*24*60 + int(m.Time.Hours())*60 + int(m.Time.Minutes())%60
}
