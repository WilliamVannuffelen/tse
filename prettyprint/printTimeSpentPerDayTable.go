package prettyprint

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/williamvannuffelen/tse/workitem"
	"os"
	"reflect"
	"sort"
	"text/tabwriter"
	"time"
)

func PrintTimeSpentPerDayTable(timeSpentPerDay []workitem.TimeSpentPerDay, filterDate string) {
	sort.Slice(timeSpentPerDay, func(i, j int) bool {
		date1, _ := time.Parse("2006-01-02", timeSpentPerDay[i].Date)
		date2, _ := time.Parse("2006-01-02", timeSpentPerDay[j].Date)
		return date1.Before(date2)
	})

	tableColor := color.New(color.FgBlue, color.Bold).SprintFunc()
	dayHeaderColor := color.New(color.FgYellow, color.Bold).SprintFunc()
	dateHeaderColor := color.New(color.FgCyan, color.Bold).SprintFunc()
	timeSpentHeaderColor := color.New(color.FgGreen, color.Bold).SprintFunc()
	greenColor := color.New(color.FgGreen, color.Bold).SprintFunc()
	redColor := color.New(color.FgRed, color.Bold).SprintFunc()

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(writer, tableColor("+-----------------+-----------------+-----------------+"))
	fmt.Fprintln(writer, tableColor("|"), dayHeaderColor("Day            "), tableColor("|"), dateHeaderColor("Date           "), tableColor("|"), timeSpentHeaderColor("TimeSpent      "), tableColor("|"))
	fmt.Fprintln(writer, tableColor("+-----------------+-----------------+-----------------+"))
	for _, entry := range timeSpentPerDay {
		if filterDate != "" && entry.Date != filterDate {
			continue
		}
		timeSpentColor := redColor
		if entry.TimeSpent > 8 {
			timeSpentColor = greenColor
		}
		fmt.Fprintln(writer, tableColor("|"), dayHeaderColor(fmt.Sprintf("%-15s", entry.Day)), tableColor("|"), dateHeaderColor(fmt.Sprintf("%-15s", entry.Date)), tableColor("|"), timeSpentColor(fmt.Sprintf("%-15.2f", entry.TimeSpent)), tableColor("|"))
	}
	fmt.Fprintln(writer, tableColor("+-----------------+-----------------+-----------------+"))
	writer.Flush()
}

// func PrintTimeSpentPerDayTable(timeSpentPerDay []workitem.TimeSpentPerDay) {
// 	sort.Slice(timeSpentPerDay, func(i, j int) bool {
// 		date1, _ := time.Parse("2006-01-02", timeSpentPerDay[i].Date)
// 		date2, _ := time.Parse("2006-01-02", timeSpentPerDay[j].Date)
// 		return date1.Before(date2)
// 	})

// 	tableColor := color.New(color.FgBlue, color.Bold).SprintFunc()
// 	dayHeaderColor := color.New(color.FgYellow, color.Bold).SprintFunc()
// 	dateHeaderColor := color.New(color.FgCyan, color.Bold).SprintFunc()
// 	timeSpentHeaderColor := color.New(color.FgGreen, color.Bold).SprintFunc()
// 	greenColor := color.New(color.FgGreen, color.Bold).SprintFunc()
// 	redColor := color.New(color.FgRed, color.Bold).SprintFunc()

// 	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
// 	fmt.Fprintln(writer, tableColor("+-----------------+-----------------+-----------------+"))
// 	fmt.Fprintln(writer, tableColor("|"), dayHeaderColor("Day            "), tableColor("|"), dateHeaderColor("Date           "), tableColor("|"), timeSpentHeaderColor("TimeSpent      "), tableColor("|"))
// 	fmt.Fprintln(writer, tableColor("+-----------------+-----------------+-----------------+"))
// 	for _, entry := range timeSpentPerDay {
// 		timeSpentColor := redColor
// 		if entry.TimeSpent > 8 {
// 			timeSpentColor = greenColor
// 		}
// 		fmt.Fprintln(writer, tableColor("|"), dayHeaderColor(fmt.Sprintf("%-15s", entry.Day)), tableColor("|"), dateHeaderColor(fmt.Sprintf("%-15s", entry.Date)), tableColor("|"), timeSpentColor(fmt.Sprintf("%-15.2f", entry.TimeSpent)), tableColor("|"))
// 	}
// 	fmt.Fprintln(writer, tableColor("+-----------------+-----------------+-----------------+"))
// 	writer.Flush()
// }

func PrintStructFields(v interface{}) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		fmt.Println("Expected a struct")
		return
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i).Interface()
		fmt.Printf("%s: %v\n", field.Name, value)
	}
}

func PrintTimeSpentWeekTotal(totalTimeSpent float64) {
	blueBoldColor := color.New(color.FgBlue, color.Bold).SprintFunc()
	greenColor := color.New(color.FgGreen, color.Bold).SprintFunc()
	redColor := color.New(color.FgRed, color.Bold).SprintFunc()

	timeSpentColor := redColor
	if totalTimeSpent > 40 {
		timeSpentColor = greenColor
	}

	fmt.Println(blueBoldColor("Week Total Time Spent:"), timeSpentColor(fmt.Sprintf("%.2f hours", totalTimeSpent)))
}
