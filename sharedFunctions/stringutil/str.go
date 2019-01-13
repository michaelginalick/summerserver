package str

import (
	"../calendar"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// ExtractMonthDate :  extracts the month from a string
func ExtractMonthDate(line string) (string, int, error) {
	parseDate := parseFields(line)

	for i := 0; i < len(parseDate); i++ {

		monthValue := calendar.GetMonth(parseDate[i])

		if monthValue.Name != "" {
			return strings.ToLower(monthValue.Name), i + 1, nil
		}
	}
	return "", 0, errors.New("No date is listed with this event")
}

// ExtractDays :  extracts the days from a string
func ExtractDays(s string, i int) []string {
	date := parseFields(s)
	
	var x strings.Builder
	x.WriteString("")
	for j := i; j < len(date); j++ {
		x.WriteString(" ")
		x.WriteString(string(date[j]))
	}

	re := regexp.MustCompile("[0-9]+")
	return re.FindAllString(x.String(), -1)
}

// FirstAndLastElement :  gets first and last element in string slice
func FirstAndLastElement(days []string) (int, int) {
	first := days[0]
	last := days[len(days)-1]

	firstInt, _ := ConvInt(first)
	lastInt, _ := ConvInt(last)

	return firstInt, lastInt
}

func parseFields(s string) []string {
	return strings.Fields(s)
}

// ConvInt :  converts a string to int
func ConvInt(s string) (int, error) {
	i, err := strconv.Atoi(s)

	if err != nil {
		return 0, errors.New("cannot convert to integer")
	}

	return i, nil
}
