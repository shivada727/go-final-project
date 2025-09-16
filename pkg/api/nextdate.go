package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DateLayout = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	repeat = strings.TrimSpace(repeat)

	if repeat == "" {
		return "", errors.New("repeat is empty")
	}

	start, err := time.Parse(DateLayout, dstart)

	if err != nil {
		return "", fmt.Errorf("bad dstart: %w", err)
	}

	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	parts := strings.Fields(repeat)

	switch parts[0] {
	case "y":
		if len(parts) != 1 {
			return "", errors.New("bad format for 'y'")
		}

		cur := start

		for {
			cur = cur.AddDate(1, 0, 0)

			if cur.After(now) {
				break
			}
		}

		return cur.Format(DateLayout), nil

	case "d":
		if len(parts) != 2 {
			return "", errors.New("bad format for 'd'")
		}

		n, err := strconv.Atoi(parts[1])

		if err != nil || n < 1 || n > 400 {
			return "", errors.New("bad days interval (1..400)")
		}

		cur := start

		for {
			cur = cur.AddDate(0, 0, n)
			if cur.After(now) {
				break
			}
		}

		return cur.Format(DateLayout), nil

	case "w":
		if len(parts) != 2 {
			return "", errors.New("bad format for 'w'")
		}

		allowed := make(map[int]bool, 7)

		for _, s := range strings.Split(parts[1], ",") {
			s = strings.TrimSpace(s)
			v, err := strconv.Atoi(s)

			if err != nil || v < 1 || v > 7 {
				return "", errors.New("weekday must be 1..7")
			}

			allowed[v] = true
		}

		cur := start

		for i := 0; i < 2000; i++ {
			cur = cur.AddDate(0, 0, 1)
			if cur.After(now) {
				if wd := weekday127(cur.Weekday()); allowed[wd] {
					return cur.Format(DateLayout), nil
				}
			}
		}

		return "", errors.New("no next date found for 'w'")

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("bad format for 'm'")
		}

		type tMonthSet [13]bool
		type tDaySet struct {
			d          [32]bool
			last       bool
			lastMinus1 bool
		}

		var days tDaySet

		for _, s := range strings.Split(parts[1], ",") {

			s = strings.TrimSpace(s)

			if s == "" {
				return "", errors.New("bad day entry")
			}

			switch s {
			case "-1":
				days.last = true
			case "-2":
				days.lastMinus1 = true
			default:
				v, err := strconv.Atoi(s)
				if err != nil || v < 1 || v > 31 {
					return "", errors.New("day must be 1..31 or -1/-2")
				}

				days.d[v] = true
			}
		}

		var months tMonthSet

		limitByMonths := false

		if len(parts) == 3 {
			limitByMonths = true

			for _, s := range strings.Split(parts[2], ",") {

				s = strings.TrimSpace(s)

				if s == "" {
					return "", errors.New("bad month entry")
				}

				v, err := strconv.Atoi(s)

				if err != nil || v < 1 || v > 12 {
					return "", errors.New("month must be 1..12")
				}

				months[v] = true
			}
		}

		monthAllowed := func(m time.Month) bool {
			if !limitByMonths {
				return true
			}

			return months[int(m)]
		}

		dayAllowed := func(t time.Time) bool {

			if t.Day() <= 31 && days.d[t.Day()] {
				return true
			}

			last := lastDayOfMonth(t)

			if days.last && t.Day() == last {
				return true
			}

			if days.lastMinus1 && t.Day() == last-1 {
				return true
			}

			return false
		}

		cur := start
		for i := 0; i < 5000; i++ {
			cur = cur.AddDate(0, 0, 1)

			if cur.After(now) && monthAllowed(cur.Month()) && dayAllowed(cur) {
				return cur.Format(DateLayout), nil
			}
		}

		return "", errors.New("no next date found for 'm'")

	default:
		return "", errors.New("unsupported repeat rule")
	}
}

func weekday127(w time.Weekday) int {
	if w == time.Sunday {
		return 7
	}

	return int(w)
}

func lastDayOfMonth(t time.Time) int {
	firstNext := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())

	last := firstNext.AddDate(0, 0, -1)

	return last.Day()
}
