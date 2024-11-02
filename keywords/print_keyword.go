package keywords

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"sort"
)

func PrintKeywordAsJSON(keyword map[string]string) error {
	key := keyword["keyword"]
	data := make(map[string]string)
	for k, v := range keyword {
		if k != "keyword" {
			data[k] = v
		}
	}
	output := map[string]map[string]string{
		key: data,
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonData))
	return nil
}

func PrettyPrintKeyword(keyword map[string]string) {
	blue := color.New(color.FgBlue, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	key := keyword["keyword"]
	fmt.Println("{")
	color.New(color.FgBlue, color.Bold).Printf("  \"%s\": {\n", key)
	keys := make([]string, 0, len(keyword))
	for k := range keyword {
		if k != "keyword" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for i, k := range keys {
		v := keyword[k]
		var coloredKey string
		switch k {
		case "description":
			coloredKey = blue(k)
		case "project":
			coloredKey = cyan(k)
		case "jira-ref":
			coloredKey = magenta(k)
		case "app-ref":
			coloredKey = yellow(k)
		default:
			coloredKey = k
		}
		if i == len(keys)-1 {
			fmt.Printf("    \"%s\": \"%s\"\n", coloredKey, green(v))
		} else {
			fmt.Printf("    \"%s\": \"%s\",\n", coloredKey, green(v))
		}
	}
	color.New(color.FgBlue, color.Bold).Println("  }")
	fmt.Println("}")
}
