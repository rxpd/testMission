package main

import (
	"avitoTelegram/config"
	cmd "avitoTelegram/utils/commands"
	"avitoTelegram/utils/logger"
	msgFormater "avitoTelegram/utils/messageFormater"
	"fmt"
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
	"strings"
)

var bot *tgapi.BotAPI

func main() {
	initializeDB(config.DBConf)

	initSelenium()
	go startPriceChecker()

	runCyborg(false, config.TelegramToken)
}

func runCyborg(debugMode bool, token string) {
	// TG SETUP
	var err error
	bot, err = tgapi.NewBotAPI(token)
	logger.LogFatalIf(err)
	bot.Debug = debugMode
	u := tgapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	logger.Log("T-800 successfully running")

	// HANDLER
	lastUserCommand := map[int64]string{}
	for update := range updates {

		// ИНЛАЙН ОТВЕТЫ
		if update.CallbackQuery != nil {
			switch lastUserCommand[int64(update.CallbackQuery.From.ID)] { // если последнее сообщение пользователя была команда
			case cmd.Unsubscribe: // если команда подписки
				goodID, err := strconv.Atoi(update.CallbackQuery.Data)
				if err != nil {
					logger.LogError(err)
					go sendMarkdownMessage(update.Message.Chat.ID, msgFormater.InternalServerError, update)
					continue
				}
				responseMessage := Unsubscribe(goodID, update.CallbackQuery.From.ID)
				go sendCallbackMarkdownMessage(update.CallbackQuery.From.ID, responseMessage)
				continue
			case cmd.ManualCheck:
				goodID, err := strconv.Atoi(update.CallbackQuery.Data)
				if err != nil {
					logger.LogError(err)
					go sendMarkdownMessage(update.Message.Chat.ID, msgFormater.InternalServerError, update)
					continue
				}
				responseMessage := ManualCheckPrice(goodID, update.CallbackQuery.From.ID)
				go sendCallbackMarkdownMessage(update.CallbackQuery.From.ID, responseMessage)
				continue
			}
		}
		//if update.Message == nil { TODO: нужно ли это?
		//	continue
		//}
		chatID := update.Message.Chat.ID // для удобства

		// КОМАНДЫ
		if update.Message.IsCommand() { // если команда
			lastUserCommand[chatID] = update.Message.Command() // записываю ID пользователя и его последнюю команду

			switch update.Message.Command() { // перебираю команды
			case cmd.Start:
				isNewUser, err := IsNewUser(chatID, update.Message.From.UserName)
				if err != nil {
					logger.LogError(err)
					go sendMarkdownMessage(chatID, msgFormater.InternalServerError, update)
					continue
				}
				if isNewUser {
					go sendMarkdownMessage(chatID, strings.Replace(msgFormater.NewUserGreeting, "name", "@"+update.Message.From.UserName, 1), update)
					continue
				} else {
					go sendMarkdownMessage(chatID, msgFormater.Help, update)
					continue
				}
			case cmd.Subscribe: // если команда подписки
				go sendMessage(chatID, msgFormater.SubscribeCommandHelp, update)
				continue
			case cmd.Unsubscribe, cmd.ManualCheck: // если команда отписки
				subscribesList, err := GetSubscribesList(chatID)
				if err != nil {
					go sendMarkdownMessage(chatID, msgFormater.InternalServerError, update)
					continue
				}
				keyboard := MakeSubscribesKeyboard(subscribesList)                          // генерирую клавиатуру
				go sendInlineKeyboard(chatID, msgFormater.ChoseSubscribe, keyboard, update) // отправляю клавиатуру
			//case cmd.ManualCheck: // дублирование с Unsubscribe, стоит зарефакторить
			//	var subscribesList []SubscribesList
			//	query := `select * from select_subscribes_list(chat_id_in := $1);` // собираю все подписки пользователя по chat_id
			//	err = DB.Select(&subscribesList, query, chatID)
			//	if err != nil {
			//		logger.LogError(err)
			//		go sendMarkdownMessage(chatID, msgFormater.InternalServerError, update)
			//		continue
			//	}
			//	keyboard := MakeSubscribesKeyboard(subscribesList)                          // генерирую клавиатуру
			//	go sendInlineKeyboard(chatID, msgFormater.ChoseSubscribe, keyboard, update) // отправляю клавиатуру
			case cmd.Help:
				go sendMessage(chatID, msgFormater.Help, update)
			case cmd.Cancel:
				delete(lastUserCommand, chatID)
				go sendMessage(chatID, msgFormater.CommandCanceled, update)
			}
			continue
		}

		// ОБЫЧНЫЙ ТЕКСТ
		if update.InlineQuery == nil { // если обычный текст
			switch lastUserCommand[chatID] { // если последнее сообщение пользователя была команда
			case cmd.Subscribe: // если команда подписки
				responseMessage := Subscribe(update.Message.Text, chatID)
				go sendMarkdownMessage(chatID, responseMessage, update)
				continue
			}
			isNewUser, err := IsNewUser(chatID, update.Message.From.UserName)
			if err != nil {
				logger.LogError(err)
				go sendMarkdownMessage(chatID, msgFormater.InternalServerError, update)
				continue
			}
			if isNewUser {
				go sendMarkdownMessage(chatID, strings.Replace(msgFormater.NewUserGreeting, "name", "@"+update.Message.From.UserName, 1), update)
				continue
			}
		}

		go sendMessage(chatID, msgFormater.DoesNotUnderstand, update) // если не понял ничего
	}
}

func initializeDB(config config.DBConfig) {
	fmt.Println()
	var err error
	connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s  sslmode=%s", config.Address, config.Port, config.Username, config.Password, config.DBName, config.SSLMode)

	if DB, err = sqlx.Connect(config.Driver, connectionUrl); err != nil {
		log.Fatal(fmt.Sprintf("\nCannot connect to \"%s\" database\n", config.DBName))
	}

	fmt.Printf("\nSuccessful connection \"%s\" database\n", config.DBName)
}
