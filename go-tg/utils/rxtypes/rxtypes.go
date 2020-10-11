package rxtypes

const (
	// tg
	MsgInternalServerError          = "500 HTTP код, свяжитесь с разработчиком - @cbrrrrrrrrr."
	MsgForNewUser                   = "Привет name! Пришли мне ссылку на объявление и я буду за ним следить, если цена поменяется я сообщу тебе. Также ты можешь вручную проверять цену."
	SubscribeCommandHelp            = "Отправьте ссылку на объявление."
	UnsubscribeCommandHelp          = "Выбери объявление."
	DoesNotUnderstand               = "Я тебя вообще не понял. Введи /help."
	SuccessfullySubscribed          = "Теперь я слежу за ценой [title](url), на данный момент она составляет - price ₽."
	MsgInternalServerErrorOrNotFund = "Либо на сервере что-то не так, либо ты ввел не существующе объявление."
	AlreadySubscribed               = "Ты уже подписан на это объявление. Как только что-то поменяется - я сообщу."
	Successfully                    = "Успешно."
	PriceChanged                    = "Цена на [title](url) изменилась, с oldPrice на newPrice ₽."

	CmdSubscribe   = "subscribe"
	CmdUnsubscribe = "unsubscribe"
	CmdManualCheck = "check_price"

	// other
	EmptyString = ""
)
