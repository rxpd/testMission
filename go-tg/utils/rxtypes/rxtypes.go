package rxtypes

const (
	// tg
	MsgInternalServerError          = "500 HTTP код, свяжитесь с разработчиком - @cbrrrrrrrrr."
	MsgForNewUser                   = "Привет name! Пришли мне ссылку на объявление и я буду за ним следить, если цена поменяется я сообщу тебе. Также ты можешь вручную проверять цену."
	SubscribeCommandHelp            = "Отправьте ссылку на объявление."
	UnsubscribeCommandHelp          = "Выбери объявление."
	DoesNotUnderstand               = "Я тебя вообще не понял. Введи /help."
	SuccessfullySubscribed          = "Теперь я слежу за ценой [title](url), на данный момент она составляет - *price ₽*."
	MsgInternalServerErrorOrNotFund = "Либо на сервере что-то не так, либо ты ввел не существующе объявление."
	AlreadySubscribed               = "Ты уже подписан на это объявление. Как только что-то поменяется - я сообщу."
	Successfully                    = "Успешно."
	PriceChanged                    = "Цена на [title](url) изменилась, с *oldPrice* на *newPrice ₽*."
	ChoseSubscribe                  = "Выбери подписку:"
	SuccessfulUnsubscribed          = "Теперь я не не слежу за [title](url). Последня цена - *price ₽*."
	CommandCanceled                 = "Команда отменена."
	PriceDoesNotChanged             = "Цена на [title](url) не изменилась. Последня цена - *price ₽*."
	CooldownLimit                   = "Слишком часто проверяешь, подожди. Последня цена на [title](url) - *price ₽*."
	//PriceChanged                    = "Цена на [title](url) изменилась, с oldPrice на newPrice e ₽*."

	CmdSubscribe   = "subscribe"
	CmdUnsubscribe = "unsubscribe"
	CmdManualCheck = "check_price"
	CmdCancel      = "cancel"

	// other
	EmptyString = ""
)
