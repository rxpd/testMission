package main

import (
	"avitoTelegram/DBModule"
	"avitoTelegram/config"
	"avitoTelegram/models"
	"avitoTelegram/utils"
	"avitoTelegram/utils/dbResponses"
	"avitoTelegram/utils/logger"
	"avitoTelegram/utils/rxtypes"
	"errors"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/tebeka/selenium"
	"log"
	"strconv"
	"strings"

	"time"
)

var bot *tgbotapi.BotAPI
var wd selenium.WebDriver

func runCyborg(debugMode bool, token string) {
	DB := DBModule.GetDB()
	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = debugMode

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	lastUserCommand := map[int64]string{}
	updates, err := bot.GetUpdatesChan(u)
	logger.Log("T-800 successfully running")
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			lastUserCommand[update.Message.Chat.ID] = update.Message.Command()

			switch update.Message.Command() {
			case "subscribe":
				go sendMessage(update.Message.Chat.ID, rxtypes.SubscribeCommandHelp, update, bot)
				continue
			case "unsubscribe":
				go sendMessage(update.Message.Chat.ID, rxtypes.SubscribeCommandHelp, update, bot)
				continue
			}
		} else if update.InlineQuery == nil {
			var dbResponse string
			query := `select * from new_message_handler(chat_id_in := $1, user_name_in := $2);`
			err := DB.Get(&dbResponse, query, update.Message.Chat.ID, update.Message.From.UserName)
			if err != nil {
				logger.LogError(err)
				go sendMessage(update.Message.Chat.ID, rxtypes.MsgInternalServerError, update, bot)
				continue
			}
			if dbResponse == dbResponses.NewUser {
				go sendMessage(update.Message.Chat.ID, strings.Replace(rxtypes.MsgForNewUser, "name", "@"+update.Message.From.UserName, 1), update, bot)
				continue
			}

			switch lastUserCommand[update.Message.Chat.ID] {
			case rxtypes.CmdSubscribe:
				responseMessage := Subscribe(update.Message.Text, update.Message.Chat.ID)
				go sendMarkdownMessage(update.Message.Chat.ID, responseMessage, update, bot)
				continue
			case rxtypes.CmdUnsubscribe:
				responseMessage := Unsubscribe(update.Message.Text, update.Message.Chat.ID)
				go sendMessage(update.Message.Chat.ID, responseMessage, update, bot)
				continue
			}
		}
		go sendMessage(update.Message.Chat.ID, rxtypes.DoesNotUnderstand, update, bot)
	}
}

//func sendKeyboard(chatID int64, message string, keyboard tgbotapi.ReplyKeyboardMarkup, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
//	msg := tgbotapi.NewMessage(chatID, message)
//	msg.ReplyMarkup = makeAdsListKeyboard()
//	//msg.ReplyMarkup = keyboard
//	msg.ReplyToMessageID = update.Message.MessageID
//	_, err := bot.Send(msg)
//	logger.LogErrorIf(err)
//}

func sendMessage(chatID int64, message string, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyToMessageID = update.Message.MessageID
	_, err := bot.Send(msg)
	logger.LogErrorIf(err)
}

func sendMarkdownMessage(chatID int64, message string, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "Markdown"
	_, err := bot.Send(msg)
	logger.LogErrorIf(err)
}

func main() {
	runCyborg(false, config.TelegramToken)
}

func Subscribe(url string, chatID int64) string { // подписка на товар
	DB := DBModule.GetDB()
	price, title, err := checkPriceByURL(url) // получаю цену, ошибку
	if err != nil {
		logger.LogError(err)
		return rxtypes.MsgInternalServerErrorOrNotFund
	}

	query := `select * from subscribe(chat_id_in := $1, good_url_in := $2, price_in := $3)`
	var DBResponse string
	err = DB.Get(&DBResponse, query, chatID, url, price)
	if err != nil {
		logger.LogError(err)
		return rxtypes.MsgInternalServerErrorOrNotFund
	}

	if DBResponse == dbResponses.SubscribeAlreadyExists {
		return rxtypes.AlreadySubscribed
	}

	return utils.SuccessfullySubscribedFormat(title, url, utils.PriceBeautify(strconv.Itoa(price)))
}

func Unsubscribe(url string, chatID int64) string {
	DB := DBModule.GetDB()
	query := `select * from unsubscribe(chat_id_in := $1, url_in := $2)`
	_, err := DB.Exec(query, chatID, url)
	if err != nil {
		logger.LogError(err)
		return rxtypes.MsgInternalServerErrorOrNotFund
	}

	return rxtypes.Successfully
}

//func GetSubscribes(response http.ResponseWriter, request *http.Request) { // тут должна быть авторизация
//	DB := DBModule.GetDB()
//	email := utils.GetRequestParams(request.URL.Query())["email"]
//	query := `select get_subscribes(email_in := $1);`
//	var data json.RawMessage
//	err := DB.Get(&data, query, email)
//	if err != nil {
//		responses.Error(response, http.StatusInternalServerError, err)
//		return
//	}
//	responses.JSON(response, http.StatusOK, data)
//}
//
//func Unsubscribe(response http.ResponseWriter, request *http.Request) {
//	DB := DBModule.GetDB()
//	var params models.ManualCheckParams
//	err := json.NewDecoder(request.Body).Decode(&params) // заполняю структуру с параметрами запроса
//	if err != nil {
//		responses.Error(response, http.StatusPreconditionFailed, err)
//		return
//	}
//	err = params.Validate()
//	if err != nil { // если параметры не прошли валидацию, отправлю 412 код
//		responses.Error(response, http.StatusPreconditionFailed, err)
//		return
//	}
//	query := `select * from unsubscribe(user_id_in := $1, good_id_in := $2);`
//	_, err = DB.Exec(query, params.UserID, params.GoodID)
//	if err != nil {
//		responses.Error(response, http.StatusInternalServerError, err)
//		return
//	}
//}
//
//func ManualCheckPrice(response http.ResponseWriter, request *http.Request) {
//	DB := DBModule.GetDB()
//	var params models.ManualCheckParams
//	err := json.NewDecoder(request.Body).Decode(&params) // заполняю структуру с параметрами запроса
//	if err != nil {
//		responses.Error(response, http.StatusPreconditionFailed, err)
//		return
//	}
//	err = params.Validate()
//	if err != nil { // если параметры не прошли валидацию, отправлю 412 код
//		responses.Error(response, http.StatusPreconditionFailed, err)
//		return
//	}
//
//	query := `select manual_check_validation(user_id_in := $1, good_id_in := $2, cooldown_in_minutes := $3);`
//	var dbResponse string
//	err = DB.Get(&dbResponse, query, params.UserID, params.GoodID, config.ManualCheckCooldownInMinutes)
//	if err != nil {
//		responses.Error(response, http.StatusInternalServerError, err)
//		return
//	}
//	var responseMessage models.ResponseMessage
//	if !strings.Contains(dbResponse, "https://www.avito") {
//		responseMessage.Message = dbResponse
//		responses.JSON(response, http.StatusForbidden, responseMessage)
//		return
//	}
//
//	price, _, err := checkPriceByURL(dbResponse) // получаю цену, ошибку
//	if err != nil {
//		responses.Error(response, http.StatusPreconditionFailed, err)
//		return
//	}
//	responseMessage.Message = strconv.Itoa(price)
//	responses.JSON(response, http.StatusOK, responseMessage)
//	go updateLastCheck(DB, int(params.UserID), int(params.GoodID), price) // обновляю в горутине, чтобы сократить время запроса
//}
//
//func updateLastCheck(DB *sqlx.DB, userID, goodID, price int) {
//	query := `select update_last_check(user_id_in := $1, good_id_in := $2, price_in := $3);`
//	_, err := DB.Exec(query, userID, goodID, price)
//	utils.LogErrorIf(err)
//}
//
func checkPriceByURL(url string) (int, string, error) {

	elem, err := wd.FindElement(selenium.ByTagName, "body")
	if err != nil {
		return 0, "", err
	}
	defer elem.SendKeys("CONTROL + W") // закрытие вкладки
	err = elem.SendKeys("CONTROL + T") // открытие вкладки
	if err != nil {
		return 0, "", err
	}
	tabs, err := wd.WindowHandles() // получения списка вкладок
	if err != nil {
		return 0, "", err
	}
	err = wd.SwitchWindow(tabs[len(tabs)-1]) // переход на последнюю открытую вкладку
	if err != nil {
		return 0, "", err
	}
	if err := wd.Get(url); err != nil {
		return 0, "", err
	}
	a, err := wd.Title()
	if err != nil {
		return 0, "", err
	}
	if !checkPageExists(a) {
		return 0, "", errors.New("page does not exists")
	}
	elem, err = wd.FindElement(selenium.ByCSSSelector, ".js-item-price")
	if err != nil {
		return 0, "", err
	}
	priceStr, err := elem.GetAttribute("content")
	if err != nil {
		return 0, "", err
	}

	price, err := strconv.Atoi(priceStr)
	if err != nil {

		return 0, "", err
	}
	elem, err = wd.FindElement(selenium.ByCSSSelector, ".title-info-title-text")
	if err != nil {
		return 0, "", err
	}
	title, err := elem.Text()
	if err != nil {
		return 0, "", err
	}
	_ = elem.SendKeys("CONTROL + W") // закрытие вкладки
	return price, title, nil
}

//
func checkPageExists(title string) bool {
	if title == "Ошибка 404. Страница не найдена — Объявления на сайте Авито" {
		return false
	}
	return true
}

//
func startPriceChecker() {
	DB := DBModule.GetDB()
	t := time.NewTicker(time.Second * time.Duration(config.ParseIntervalInSeconds))
	for {
		var urlsForParse []models.GoodsForCheck
		query := `select * from urls_for_parse_select();`
		err := DB.Select(&urlsForParse, query)
		logger.LogErrorIf(err)

		for _, good := range urlsForParse {
			checkedPrice, title, err := checkPriceByURL(good.Url)
			if err != nil {
				<-t.C
			}
			logger.LogErrorIf(err)
			if checkedPrice != good.Price {
				var usersForNotify []int64
				query := `select * from get_users_for_notify(good_id_in := $1, price_in := $2);`
				err := DB.Select(&usersForNotify, query, good.GoodID, checkedPrice)
				if len(usersForNotify) == 0 {
					continue
				}
				logger.LogErrorIf(err)
				//fmt.Println(urlsForParse)
				checkedPriceFormatted := utils.PriceBeautify(strconv.Itoa(checkedPrice))
				oldPriceFormatted := utils.PriceBeautify(strconv.Itoa(good.Price))
				//fmt.Println(usersForNotify, checkedPriceFormatted, oldPriceFormatted, good.Url, title)
				//err = SendNotifyEmail(usersForNotify, checkedPriceFormatted, oldPriceFormatted, good.Url, title)
				message := utils.PriceChangeFormat(title, good.Url, oldPriceFormatted, checkedPriceFormatted)
				priceUpdateMessage(usersForNotify, message)
				query = `select update_price(good_id_in := $1, price_in := $2);`
				_, err = DB.Exec(query, good.GoodID, checkedPrice)
				logger.LogErrorIf(err)
			}
		}
		<-t.C
	}
}

func priceUpdateMessage(usersForNotify []int64, message string) {
	for _, u := range usersForNotify {
		msg := tgbotapi.NewMessage(u, message)
		msg.ParseMode = "Markdown"
		_, err := bot.Send(msg)
		logger.LogErrorIf(err)
	}
}

func init() {
	// Start a Selenium WebDriver server instance (if one is not already
	// running).
	const (
		// These paths will be different on your system.
		seleniumPath    = "selenium/selenium-server-standalone-3.141.59.jar"
		geckoDriverPath = "selenium/geckodriver"
		port            = 8008
	)
	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
		selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		//selenium.Output(os.Stderr),            // Output debug information to STDERR.
	}
	selenium.SetDebug(false)
	/*service*/ _, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	logger.LogFatalIf(err)
	//defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	logger.LogFatalIf(err)
	fmt.Printf("\nSelenium successfully initialized\n")
	go startPriceChecker()
}

//func makeAdsListKeyboard( /*ads map[int]string*/) []tgbotapi.InlineKeyboardButton {
//	//var adsListKeyboard = tgbotapi.NewReplyKeyboard(
//	//	tgbotapi.NewKeyboardButtonRow(),
//	//)
//	//for i, v := range ads {
//	//	s := strconv.Itoa(i)
//	//	adsListKeyboard.Keyboard[0] =  tgbotapi.NewInlineKeyboardButtonData(v, s)
//	//}
//	var adsListKeyboard = tgbotapi.NewInlineKeyboardRow(
//		tgbotapi.NewInlineKeyboardButtonData("", ""))
//	return adsListKeyboard
//	test := tgbotapi.NewReplyKeyboard()
//	test.Keyboard = append(tgbotapi.NewReplyKeyboard().Keyboard, adsListKeyboard)
//}
//
//var menuKeyboard = tgbotapi.NewReplyKeyboard(
//	tgbotapi.NewKeyboardButtonRow(
//		tgbotapi.NewKeyboardButton("Подписаться"),
//		tgbotapi.NewKeyboardButton("Отписаться"),
//	),
//	tgbotapi.NewKeyboardButtonRow(
//		tgbotapi.NewKeyboardButton("Проверить цену"),
//	),
//)
