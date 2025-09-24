package controller

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/requests"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
)

var hhmmRe = regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)

var wdMap = map[string]time.Weekday{
	"monday":    time.Monday,
	"tuesday":   time.Tuesday,
	"wednesday": time.Wednesday,
	"thursday":  time.Thursday,
	"friday":    time.Friday,
	"saturday":  time.Saturday,
	"sunday":    time.Sunday,
}

func parseWeekday(s string) (time.Weekday, error) {
	w, ok := wdMap[strings.ToLower(strings.TrimSpace(s))]
	if !ok {
		return 0, fmt.Errorf("invalid weekday: %q", s)
	}
	return w, nil
}

func parseHHMM(s string) (time.Duration, error) {
	if !hhmmRe.MatchString(s) {
		return 0, fmt.Errorf("invalid time format (HH:MM): %q", s)
	}
	h := (int(s[0]-'0')*10 + int(s[1]-'0'))
	m := (int(s[3]-'0')*10 + int(s[4]-'0'))
	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute, nil
}

type daySpan struct{ start, end int } // минуты с начала суток [0..1440)

func validateNoOverlaps(spans []daySpan) error {
	if len(spans) <= 1 {
		return nil
	}
	sort.Slice(spans, func(i, j int) bool { return spans[i].start < spans[j].start })
	prev := spans[0]
	for i := 1; i < len(spans); i++ {
		cur := spans[i]
		if cur.start < prev.end { // пересечение (стык в prev.end == cur.start допустим)
			return fmt.Errorf("overlap: [%02d:%02d-%02d:%02d) with [%02d:%02d-%02d:%02d)",
				prev.start/60, prev.start%60, prev.end/60, prev.end%60,
				cur.start/60, cur.start%60, cur.end/60, cur.end%60)
		}
		prev = cur
	}
	return nil
}

func (h Service) SetTimetable(w http.ResponseWriter, r *http.Request) {
	req, err := requests.SetTimetable(r)
	if err != nil {
		h.log.WithError(err).Error("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	params := models.Timetable{Table: make([]models.TimeInterval, 0, len(req.Data.Attributes.Table))}
	perDay := map[time.Weekday][]daySpan{}

	for i, interval := range req.Data.Attributes.Table {
		fromWD, err := parseWeekday(interval.From.Weekday)
		if err != nil {
			ape.RenderErr(w, problems.InvalidPointer(fmt.Sprintf("data/attributes/table/%d/from/weekday", i), err))
			return
		}
		toWD, err := parseWeekday(interval.To.Weekday)
		if err != nil {
			ape.RenderErr(w, problems.InvalidPointer(fmt.Sprintf("data/attributes/table/%d/to/weekday", i), err))
			return
		}
		fromT, err := parseHHMM(interval.From.Time)
		if err != nil {
			ape.RenderErr(w, problems.InvalidPointer(fmt.Sprintf("data/attributes/table/%d/from/time", i), err))
			return
		}
		toT, err := parseHHMM(interval.To.Time)
		if err != nil {
			ape.RenderErr(w, problems.InvalidPointer(fmt.Sprintf("data/attributes/table/%d/to/time", i), err))
			return
		}

		if fromWD != toWD {
			ape.RenderErr(w, problems.InvalidPointer(
				fmt.Sprintf("data/attributes/table/%d", i),
				fmt.Errorf("interval must be within one day (from.weekday == to.weekday)"),
			))
			return
		}
		if !(fromT < toT) {
			ape.RenderErr(w, problems.InvalidPointer(
				fmt.Sprintf("data/attributes/table/%d", i),
				fmt.Errorf("from must be earlier than to: %s >= %s", interval.From.Time, interval.To.Time),
			))
			return
		}

		startMin := int(fromT / time.Minute)
		endMin := int(toT / time.Minute)
		perDay[fromWD] = append(perDay[fromWD], daySpan{start: startMin, end: endMin})

		params.Table = append(params.Table, models.TimeInterval{
			From: models.Moment{Weekday: fromWD, Time: fromT},
			To:   models.Moment{Weekday: toWD, Time: toT},
		})
	}

	for wd, spans := range perDay {
		if err := validateNoOverlaps(spans); err != nil {
			ape.RenderErr(w, problems.InvalidParameter(strings.ToLower(wd.String()), err))
			return
		}
	}

	res, err := h.domain.Place.SetTimetable(r.Context(), req.Data.Id, params)
	if err != nil {
		h.log.WithError(err).Error("could not set timetable")
		switch {
		case errors.Is(err, errx.ErrorPlaceNotFound):
			ape.RenderErr(w, problems.NotFound(fmt.Sprintf("place %s not found", req.Data.Id)))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Place(res))
}
