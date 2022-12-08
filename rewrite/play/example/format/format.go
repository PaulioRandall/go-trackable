package format

import (
	"regexp"
	"strconv"
	"time"

	"github.com/PaulioRandall/trackable/rewrite"
)

var (
	ErrFormatPkg = track.Checkpoint("play/example/format")
	ErrFormatRow = track.Error("Could not format row")
)

type PhReading struct {
	From  time.Time
	To    time.Time
	Value float64
}

func Format(data [][]string) ([]PhReading, error) {
	var results []PhReading

	for i := 0; i < len(data); i++ {
		phReading, e := formatRow(data[i])

		if e != nil {
			return nil, ErrFormatPkg.Wrap(e)
		}

		results = append(results, phReading)
	}

	return results, nil
}

func formatRow(row []string) (PhReading, error) {
	var reading, zero PhReading
	var e error

	reading.From, reading.To, e = parseDate(row[0])
	if e != nil {
		return zero, ErrFormatRow.CausedBy(e, "Could not parse date")
	}

	reading.Value, e = strconv.ParseFloat(row[1], 32)
	if e != nil {
		return zero, ErrFormatRow.CausedBy(e, "Could not parse pH value")
	}

	return reading, nil
}

func parseDate(dateStr string) (time.Time, time.Time, error) {
	if day, e := parseDay(dateStr); e == nil {
		return day, day, nil
	}

	findWhitespace := regexp.MustCompile(`\s`)
	dateStr = findWhitespace.ReplaceAllString(dateStr, "")

	if from, to, ok := parseInMonthDayRange(dateStr); ok {
		return from, to, nil
	}

	if from, to, ok := parseAcrossMonthDayRange(dateStr); ok {
		return from, to, nil
	}

	var zero time.Time
	return zero, zero, track.Untracked("Date string %q could not be parsed as day or day range", dateStr)
}

func parseInMonthDayRange(dateStr string) (from time.Time, to time.Time, ok bool) {
	inMonthDayRange := regexp.MustCompile(`^(\d{1,2})-(\d{1,2})([A-Za-z]{3})$`)
	groups, ok := findGroups(inMonthDayRange, dateStr)

	if !ok {
		return
	}

	var e error
	groups = groups[1:] // Remove full match group

	if from, e = parseDayFromParts(groups[0], groups[2]); e != nil {
		return
	}

	if to, e = parseDayFromParts(groups[1], groups[2]); e != nil {
		return
	}

	return from, to, true
}

func parseAcrossMonthDayRange(dateStr string) (from time.Time, to time.Time, ok bool) {
	acrossMonthDayRange := regexp.MustCompile(`^(\d{1,2})([A-Za-z]{3})-(\d{1,2})([A-Za-z]{3})$`)
	groups, ok := findGroups(acrossMonthDayRange, dateStr)

	if !ok {
		return
	}

	var e error
	groups = groups[1:] // Remove full match group

	if from, e = parseDayFromParts(groups[0], groups[1]); e != nil {
		return
	}

	if to, e = parseDayFromParts(groups[2], groups[3]); e != nil {
		return
	}

	return from, to, true
}

func parseDay(dateStr string) (time.Time, error) {
	const dayFmt = "_2 Jan"
	return time.Parse(dayFmt, dateStr)
}

func parseDayFromParts(day, month string) (time.Time, error) {
	return parseDay(day + " " + month)
}

func findGroups(re *regexp.Regexp, str string) ([]string, bool) {
	matches := re.FindAllStringSubmatch(str, -1)
	if len(matches) == 0 {
		return nil, false
	}
	return matches[0], true
}
