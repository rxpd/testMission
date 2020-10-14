package utils

import (
	"avitoTelegram/utils/logger"
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"time"
)

func PrintStructJSON(value interface{}) {
	result, err := json.MarshalIndent(value, "", "\t")
	logger.LogErrorIf(err)
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
