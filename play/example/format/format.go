package format

import (
	"regexp"
	"strconv"
	"time"

	"github.com/PaulioRandall/go-trackable"
)

var (
	ErrFormatPkg = track.Checkpoint("play/example/format")
	ErrFormatRow = track.Error("Could not format row")

	zeroTime time.Time
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
		return zero, ErrFormatRow.CausedBy(e, "Could not parse date %q", row[0])
	}

	reading.Value, e = strconv.ParseFloat(row[1], 32)
	if e != nil {
		return zero, ErrFormatRow.CausedBy(e, "Could not parse pH value %q", row[1])
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

	return zeroTime, zeroTime, track.Untracked("Date string %q could not be parsed as day or day range", dateStr)
}

func parseInMonthDayRange(dateStr string) (time.Time, time.Time, bool) {
	//                   from-day   to-day      month
	inMonthDayRange := `^(\d{1,2})-(\d{1,2})([A-Za-z]{3})$`
	groups, ok := findGroups(inMonthDayRange, dateStr)

	if !ok {
		return zeroTime, zeroTime, false
	}

	var e error
	groups = groups[1:] // Remove full match group

	from, to, e := parseDayRangeFromParts(groups[0], groups[2], groups[1], groups[2])
	if e != nil {
		return zeroTime, zeroTime, false
	}

	return from, to, true
}

func parseAcrossMonthDayRange(dateStr string) (time.Time, time.Time, bool) {
	//                       from-day   from-month   to-day     to-month
	acrossMonthDayRange := `^(\d{1,2})([A-Za-z]{3})-(\d{1,2})([A-Za-z]{3})$`
	groups, ok := findGroups(acrossMonthDayRange, dateStr)

	if !ok {
		return zeroTime, zeroTime, false
	}

	var e error
	groups = groups[1:] // Remove full match group

	from, to, e := parseDayRangeFromParts(groups[0], groups[1], groups[2], groups[3])
	if e != nil {
		return zeroTime, zeroTime, false
	}

	return from, to, true
}

func parseDay(dateStr string) (time.Time, error) {
	const dayFmt = "_2 Jan"
	return time.Parse(dayFmt, dateStr)
}

func parseDayRangeFromParts(
	fromDay, fromMonth string,
	toDay, toMonth string,
) (from, to time.Time, e error) {
	if from, e = parseDay(fromDay + " " + fromMonth); e != nil {
		return
	}
	to, e = parseDay(toDay + " " + toMonth)
	return
}

func findGroups(reStr, str string) ([]string, bool) {
	re := regexp.MustCompile(reStr)
	matches := re.FindAllStringSubmatch(str, -1)
	if len(matches) == 0 {
		return nil, false
	}
	return matches[0], true
}
