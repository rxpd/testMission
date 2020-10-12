package utils

import (
	"avitoTelegram/utils/logger"
	"avitoTelegram/utils/rxtypes"
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"strings"
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

func GetTextFromCommand(command string, text string) string {
	return strings.Replace(text, "/"+command+" ", "", 1)
}

func SuccessfullySubscribedFormat(title, url, price string) string {
	result := rxtypes.SuccessfullySubscribed
	result = strings.Replace(result, "title", title, 1)
	result = strings.Replace(result, "url", url, 1)
	result = strings.Replace(result, "price", price, 1)
	return result
}

func SuccessfulUnsubscribedFormat(title, url, price string) string {
	result := rxtypes.SuccessfulUnsubscribed
	result = strings.Replace(result, "title", title, 1)
	result = strings.Replace(result, "url", url, 1)
	result = strings.Replace(result, "price", price, 1)
	return result
}

func PriceChangeFormat(title, url, oldPrice, newPrice string) string {
	//"Цена на [title](url) изменилась, с *oldPrice* на *newPrice ₽*."
	result := rxtypes.PriceChanged
	result = strings.Replace(result, "title", title, 1)
	result = strings.Replace(result, "url", url, 1)
	result = strings.Replace(result, "oldPrice", oldPrice, 1)
	result = strings.Replace(result, "newPrice", newPrice, 1)
	return result
}

func PriceDoesNotChangeFormat(title, url, price string) string {
	//Цена на [title](url) не изменилась. Последня цена - *price ₽*.
	result := rxtypes.PriceDoesNotChanged
	result = strings.Replace(result, "title", title, 1)
	result = strings.Replace(result, "url", url, 1)
	result = strings.Replace(result, "price", price, 1)
	return result
}
func CooldownLimitFormatter(title, url, price string) string {
	//Слишком часто проверяешь, подожди. Последня цена на [title](url) - *price ₽*."
	result := rxtypes.CooldownLimit
	result = strings.Replace(result, "title", title, 1)
	result = strings.Replace(result, "url", url, 1)
	result = strings.Replace(result, "price", price, 1)
	return result
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
