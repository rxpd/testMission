package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"net/url"
	"time"
)

func PrintStructJSON(value interface{}) {
	result, err := json.MarshalIndent(value, "", "\t")
	LogErrorIf(err)
	fmt.Printf("\n" + string(result) + "\n")
}

type TimeSinceCounter struct {
	time  time.Time
	title string
}

func (t *TimeSinceCounter) StartTimeSince(title string) {
	t.title = title
	t.time = time.Now()
}
func (t *TimeSinceCounter) LogTimeSince() {
	cyan := color.FgCyan.Render
	fmt.Printf("\n%s %s took - %s\n", cyan("INFO"), t.title, cyan(time.Since(t.time)))
}

func RemoveDuplicatesFromSlice(slice []string) []string {
	m := make(map[string]bool)
	for _, item := range slice {
		if _, ok := m[item]; ok {
			// duplicate item
			//fmt.Println(item, "is a duplicate")
		} else {
			m[item] = true
		}
	}

	var result []string
	for item, _ := range m {
		result = append(result, item)
	}
	return result
}

func GetRequestParams(values url.Values) map[string]string {
	urlValues := make(map[string]string, len(values))
	for k, v := range values {
		urlValues[k] = v[0]
	}
	return urlValues
}

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

func PriceBeautify(price string) string {
	var result string
	for i := len(price); i > 0; i-- {
		if len(price) > 6 && (i == 6 || i == 3) {
			result += " " + string(price[len(price)-i])
			continue
		} else if len(price) > 4 && (i == 3) {
			result += " " + string(price[len(price)-i])
			continue
		}
		result += string(price[len(price)-i])

	}
	return result
}
