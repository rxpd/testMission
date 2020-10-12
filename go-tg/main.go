package main

import (
	"avitoTelegram/DBModule"
	"avitoTelegram/config"
	"avitoTelegram/utils"
	"avitoTelegram/utils/dbResponses"
	"avitoTelegram/utils/logger"
	"avitoTelegram/utils/rxtypes"
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/tebeka/selenium"
	"strconv"
	"strings"
	"time"
)

var bot *tgapi.BotAPI
var wd selenium.WebDriver
var userManualCheck []UserManualCheck

func main() {
	initSelenium()
	go startPriceChecker()

	runCyborg(false, config.TelegramToken)
}

func runCyborg(debugMode bool, token string) {
	DB := DBModule.GetDB()
	var err error
	bot, err = tgapi.NewBotAPI(token)
	logger.LogFatalIf(err)
	bot.Debug = debugMode

	u := tgapi.NewUpdate(0)
	u.Timeout = 60

	lastUserCommand := map[int64]string{}

	updates, err := bot.GetUpdatesChan(u)
	logger.Log("T-800 successfully running")

	for update := range updates {
		if update.CallbackQuery != nil { // если инлайн ответ
			//fmt.Println(lastUserCommand)
			//fmt.Println(update.CallbackQuery.From.ID)
			switch lastUserCommand[int64(update.CallbackQuery.From.ID)] { // если последнее сообщение пользователя была команда
			case rxtypes.CmdUnsubscribe: // если команда подписки
				goodID, err := strconv.Atoi(update.CallbackQuery.Data)
				if err != nil {
					logger.LogError(err)
					go sendMessage(update.Message.Chat.ID, rxtypes.MsgInternalServerErrorOrNotFund, update)
					continue
				}
				responseMessage := Unsubscribe(goodID, update.CallbackQuery.From.ID)
				go sendCallbackMarkdownMessage(update.CallbackQuery.From.ID, responseMessage)
				continue
			case rxtypes.CmdManualCheck:
				goodID, err := strconv.Atoi(update.CallbackQuery.Data)
				if err != nil {
					logger.LogError(err)
					go sendMessage(update.Message.Chat.ID, rxtypes.MsgInternalServerErrorOrNotFund, update)
					continue
				}
				responseMessage := ManualCheckPrice(goodID, update.CallbackQuery.From.ID)
				go sendCallbackMarkdownMessage(update.CallbackQuery.From.ID, responseMessage)
				continue
			}
			//fmt.Println(update.CallbackQuery.Data, update.CallbackQuery.InlineMessageID)
			//_, _ = bot.AnswerCallbackQuery(tgapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
			//
			//_, _ = bot.Send(tgapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
		}
		if update.Message == nil {
			continue
		}
		chatID := update.Message.Chat.ID // для удобства

		if update.Message.IsCommand() { // если команда
			lastUserCommand[chatID] = update.Message.Command() // записываю ID пользователя и его последнюю команду

			switch update.Message.Command() { // перебираю команды
			case rxtypes.CmdSubscribe: // если команда подписки
				go sendMessage(chatID, rxtypes.SubscribeCommandHelp, update)
				continue
			case rxtypes.CmdUnsubscribe: // если команда отписки
				var subscribesList []SubscribesList
				query := `select * from get_subscribes_list(chat_id_in := $1);` // собираю все подписки пользователя по chat_id
				err = DB.Select(&subscribesList, query, chatID)
				if err != nil {
					logger.LogError(err)
					go sendMessage(chatID, rxtypes.MsgInternalServerErrorOrNotFund, update)
					continue

				}
				keyboard := MakeSubscribesKeyboard(subscribesList)                      // генерирую клавиатуру
				go sendInlineKeyboard(chatID, rxtypes.ChoseSubscribe, keyboard, update) // отправляю клавиатуру
				continue
			case rxtypes.CmdManualCheck: // дублирование с CmdUnsubscribe, стоит зарефакторить
				var subscribesList []SubscribesList
				query := `select * from get_subscribes_list(chat_id_in := $1);` // собираю все подписки пользователя по chat_id
				err = DB.Select(&subscribesList, query, chatID)
				if err != nil {
					logger.LogError(err)
					go sendMessage(chatID, rxtypes.MsgInternalServerErrorOrNotFund, update)
					continue

				}
				keyboard := MakeSubscribesKeyboard(subscribesList)                      // генерирую клавиатуру
				go sendInlineKeyboard(chatID, rxtypes.ChoseSubscribe, keyboard, update) // отправляю клавиатуру
				continue
			case rxtypes.CmdCancel:
				delete(lastUserCommand, chatID)
				go sendMessage(chatID, rxtypes.CommandCanceled, update)
			}
			continue
		}
		if update.InlineQuery == nil { // если обычный текст
			var dbResponse string
			query := `select * from new_message_handler(chat_id_in := $1, user_name_in := $2);`
			err := DB.Get(&dbResponse, query, chatID, update.Message.From.UserName)
			if err != nil {
				logger.LogError(err)
				go sendMessage(chatID, rxtypes.MsgInternalServerError, update)
				continue
			}
			if dbResponse == dbResponses.NewUser { // если новый пользователь, отправляю приветственное сообщение
				go sendMessage(chatID, strings.Replace(rxtypes.MsgForNewUser, "name", "@"+update.Message.From.UserName, 1), update)
				continue
			}

			switch lastUserCommand[chatID] { // если последнее сообщение пользователя была команда
			case rxtypes.CmdSubscribe: // если команда подписки
				responseMessage := Subscribe(update.Message.Text, chatID)
				go sendMarkdownMessage(chatID, responseMessage, update)
				continue

			}
		}

		go sendMessage(chatID, rxtypes.DoesNotUnderstand, update) // если не понял ничего
	}
}

func Subscribe(url string, chatID int64) string { // подписка на товар
	DB := DBModule.GetDB()
	price, title, err := checkPriceByURL(url) // получаю цену, ошибку
	if err != nil {
		logger.LogError(err)
		return rxtypes.MsgInternalServerErrorOrNotFund
	}

	query := `select * from subscribe(chat_id_in := $1, good_url_in := $2, title_in := $3, price_in := $4)`
	var DBResponse string
	err = DB.Get(&DBResponse, query, chatID, url, title, price)
	if err != nil {
		logger.LogError(err)
		return rxtypes.MsgInternalServerErrorOrNotFund
	}

	if DBResponse == dbResponses.SubscribeAlreadyExists {
		return rxtypes.AlreadySubscribed
	}

	return utils.SuccessfullySubscribedFormat(title, url, utils.PriceBeautify(strconv.Itoa(price)))
}

func ManualCheckPrice(goodID, chatID int) string {
	DB := DBModule.GetDB()
	query := `select * from get_manual_check_info(good_id_in := $1, chat_id_in := $2)`
	var goodInfo GoodInfo
	err := DB.Get(&goodInfo, query, goodID, chatID)
	if err != nil {
		logger.LogError(err)
		return rxtypes.MsgInternalServerErrorOrNotFund
	}
	go updateUserManualCheck(goodID, chatID, goodInfo.Price)
	if !canBeManualChecked(goodID, chatID) {
		price := utils.PriceBeautify(strconv.Itoa(goodInfo.Price))
		return utils.CooldownLimitFormatter(goodInfo.Title, goodInfo.URL, price)
	}
	parsedPrice, title, err := checkPriceByURL(goodInfo.URL) // получаю цену, ошибку
	if err != nil {
		logger.LogError(err)
		return rxtypes.MsgInternalServerErrorOrNotFund
	}
	if parsedPrice != goodInfo.Price {
		oldPrice, newPrice := utils.PriceBeautify(strconv.Itoa(goodInfo.Price)), utils.PriceBeautify(strconv.Itoa(parsedPrice))
		go updateLastManualCheck(chatID, goodID, parsedPrice)
		return utils.PriceChangeFormat(title, goodInfo.URL, oldPrice, newPrice)
	} else {
		return utils.PriceDoesNotChangeFormat(goodInfo.URL, title, utils.PriceBeautify(strconv.Itoa(goodInfo.Price)))
	}
}

func canBeManualChecked(goodID, chatID int) bool {
	for _, checks := range userManualCheck {
		if checks.ChatID == chatID {
			for _, good := range checks.GoodChecks {
				if good.GoodID == goodID && time.Time.Before(time.Now(), time.Time.Add(good.LastCheckTime, time.Minute*time.Duration(config.ManualCheckCooldownInMinutes))) {
					return false
				}
			}
		}
	}
	return true
}

func updateUserManualCheck(goodID, chatID, price int) {
	for _, checks := range userManualCheck {
		if checks.ChatID == chatID {
			for _, good := range checks.GoodChecks {
				if good.GoodID == goodID {
					good.LastCheckTime = time.Now()
					good.LastPrice = price
					return
				}
			}
		}
	}
	var userLocalMnlCheck UserManualCheck // если это первая ручная проверка
	userLocalMnlCheck.ChatID = chatID
	var goodManualCheck GoodManualChecked
	goodManualCheck.LastPrice = price
	goodManualCheck.GoodID = goodID
	goodManualCheck.LastCheckTime = time.Now()
	userLocalMnlCheck.GoodChecks = append(userLocalMnlCheck.GoodChecks, goodManualCheck)
	userManualCheck = append(userManualCheck, userLocalMnlCheck)
}

func Unsubscribe(goodID int, chatID int) string {
	DB := DBModule.GetDB()
	query := `select * from unsubscribe(good_id_in := $1, chat_id_in := $2)`
	var response GoodInfo
	err := DB.Get(&response, query, goodID, chatID)
	if err != nil {
		logger.LogError(err)
		return rxtypes.MsgInternalServerErrorOrNotFund
	}
	go deleteNonSubscribedGoods()
	return utils.SuccessfulUnsubscribedFormat(response.Title, response.URL, utils.PriceBeautify(strconv.Itoa(response.Price)))
}

func deleteNonSubscribedGoods() {
	DB := DBModule.GetDB()
	query := `select delete_non_subscribed_goods();`
	_, err := DB.Exec(query)
	logger.LogFatalIf(err)
}

func updateLastManualCheck(userID, goodID, price int) {
	DB := DBModule.GetDB()
	query := `select update_last_manual_check(chat_id_in := $1, good_id_in := $2, price_in := $3);`
	_, err := DB.Exec(query, userID, goodID, price)
	logger.LogErrorIf(err)
}
