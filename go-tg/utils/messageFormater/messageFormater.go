package messageFormater

import "strings"

const (
	InternalServerError    = "500 HTTP код, свяжись с разработчиком - [rxpd](https://github.com/rxpd)"
	NewUserGreeting        = "Привет name! Пришли мне ссылку на объявление и я буду за ним следить, если цена поменяется я сообщу тебе. Также ты можешь вручную проверять цену. Используй */subscribe*"
	SubscribeCommandHelp   = "Отправь ссылку на объявление."
	UnsubscribeCommandHelp = "Выбери объявление."
	SuccessfullySubscribed = "Теперь я слежу за ценой [title](url), на данный момент она составляет - *price ₽*."
	GoodNotFound           = "Я не нашел объявление. Проверь ссылку или разработчиком - [rxpd](https://github.com/rxpd)."
	AlreadySubscribed      = "Ты уже подписан на это объявление. Как только что-то поменяется - я сообщу."
	Successfully           = "Успешно."
	PriceChanged           = "Цена на [title](url) изменилась, с *oldPrice* на *newPrice ₽*."
	ChoseSubscribe         = "Выбери подписку:"
	SuccessfulUnsubscribed = "Теперь я не не слежу за [title](url). Последня цена - *price ₽*."
	CommandCanceled        = "Команда отменена."
	PriceDoesNotChanged    = "Цена на [title](url) не изменилась. Последня цена - *price ₽*."
	CooldownLimit          = "Слишком часто проверяешь, подожди. Последня цена на [title](url) - *price ₽*."
	NoSubscribes           = "У тебя нет подписок."
	DoesNotUnderstand      = "Я тебя вообще не понял. Введи /help."
	Help                   = "Помощь."
)

func SuccessfullySubscribedFormat(title, url, price string) string {
	result := SuccessfullySubscribed
	result = strings.Replace(result, "title", title, 1)
	result = strings.Replace(result, "url", url, 1)
	result = strings.Replace(result, "price", price, 1)
	return result
}

func SuccessfulUnsubscribedFormat(title, url, price string) string {
	result := SuccessfulUnsubscribed
	result = strings.Replace(result, "title", title, 1)
	result = strings.Replace(result, "url", url, 1)
	result = strings.Replace(result, "price", price, 1)
	return result
}

func PriceChangedFormat(title, url, oldPrice, newPrice string) string {
	result := PriceChanged
	result = strings.Replace(result, "title", title, 1)
	result = strings.Replace(result, "url", url, 1)
	result = strings.Replace(result, "oldPrice", oldPrice, 1)
	result = strings.Replace(result, "newPrice", newPrice, 1)
	return result
}

func PriceDoesNotChangedFormat(title, url, price string) string {
	result := PriceDoesNotChanged
	result = strings.Replace(result, "title", title, 1)
	result = strings.Replace(result, "url", url, 1)
	result = strings.Replace(result, "price", price, 1)
	return result
}
func CooldownLimitFormat(title, url, price string) string {
	result := CooldownLimit
	result = strings.Replace(result, "title", title, 1)
	result = strings.Replace(result, "url", url, 1)
	result = strings.Replace(result, "price", price, 1)
	return result
}
