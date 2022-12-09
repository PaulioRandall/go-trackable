package example

import (
	"fmt"
	"time"

	"github.com/PaulioRandall/go-trackable/play/example/format"
)

const timePrintFormat = "Jan _2"

func printReadings(readings []format.PhReading) {
	fmt.Println(" Row |       Date       |  pH")
	fmt.Println("———————————————————————————————")

	for i, r := range readings {
		date := formatDateRange(r.From, r.To)
		fmt.Printf("%4d | %-16v | %.2f\n", i+1, date, r.Value)
	}
}

func formatDateRange(from, to time.Time) string {
	if from == to {
		return formatSingleDay(from)
	}

	if from.Month() == to.Month() {
		return formatInMonthDateRange(from, to)
	}

	return formatAcrossMonthDateRange(from, to)
}

func formatSingleDay(day time.Time) string {
	return day.Format(timePrintFormat)
}

func formatInMonthDateRange(from, to time.Time) string {
	fromStr := from.Format(timePrintFormat)
	return fmt.Sprintf("%s%s%d", fromStr, "-", to.Day())
}

func formatAcrossMonthDateRange(from, to time.Time) string {
	fromStr := from.Format(timePrintFormat)
	toStr := to.Format(timePrintFormat)
	return fmt.Sprintf("%s%s%s", fromStr, " ↷  ", toStr)
}
