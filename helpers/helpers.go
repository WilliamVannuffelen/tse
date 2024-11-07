package helpers

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

func NewErrorStackTraceString(message string) string {
	pc, file, line, _ := runtime.Caller(1)

	callerName := runtime.FuncForPC(pc).Name()
	callerParts := strings.Split(callerName, "/")
	callerName = callerParts[len(callerParts)-1]

	parts := strings.Split(file, "/")
	fileName := strings.Join(parts[len(parts)-2:], "/")

	return fmt.Sprintf("%s:%d - %s - %s \n from", fileName, line, callerName, message)
}

func GetCurrentWeekDate(now func() time.Time) string {
	today := now().Weekday()

	offset := int(time.Monday - today)
	if offset > 0 {
		offset = -6
	}
	monday := now().AddDate(0, 0, offset).Format("2006-01-02")

	return monday
}

func GetStartOfWeek(date string) (string, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", fmt.Errorf("%s %w", NewErrorStackTraceString("failed to parse date"), err)
	}

	weekday := parsedDate.Weekday()
	offset := int(time.Monday - weekday)
	if offset > 0 {
		offset = -6
	}
	startOfWeek := parsedDate.AddDate(0, 0, offset).Format("2006-01-02")

	return startOfWeek, nil
}

func getDayOffset(day string) (int, error) {
	days := map[string]int{
		"mon": 0,
		"tue": 1,
		"wed": 2,
		"thu": 3,
		"fri": 4,
		"sat": 5,
		"sun": 6,
	}

	offset, exists := days[day]
	if !exists {
		return 0, fmt.Errorf("%s %s", NewErrorStackTraceString("invalid day provided. Valid values: mon, tue, wed, thu, fri, sat, sun"), fmt.Errorf("keyerror: '%s' not in days map", day))
	}

	return offset, nil
}

func GetDateFromDay(day string, now func() time.Time) (string, error) {
	weekday := now().Weekday()

	offset := int(time.Monday - weekday)
	if offset > 0 {
		offset = -6
	}
	dayOffset, err := getDayOffset(day)
	if err != nil {
		return "", err
	}
	return now().AddDate(0, 0, offset+dayOffset).Format("2006-01-02"), nil
}
