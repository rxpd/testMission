package main

import (
	"avitoTelegram/utils/logger"
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

func sendInlineKeyboard(chatID int64, message string, keyboard tgapi.InlineKeyboardMarkup, update tgapi.Update) {
	msg := tgapi.NewMessage(chatID, message)
	msg.ReplyMarkup = keyboard
	msg.ReplyToMessageID = update.Message.MessageID
	_, err := bot.Send(msg)
	logger.LogErrorIf(err)
}

func sendMessage(chatID int64, message string, update tgapi.Update) {
	msg := tgapi.NewMessage(chatID, message)
	msg.ReplyToMessageID = update.Message.MessageID
	_, err := bot.Send(msg)
	logger.LogErrorIf(err)
}

func sendCallbackMessage(chatID int, message string) /*, update tgapi.Update)*/ {
	msg := tgapi.NewMessage(int64(chatID), message)
	_, err := bot.Send(msg)
	logger.LogErrorIf(err)
}

func sendCallbackMarkdownMessage(chatID int, message string) /*, update tgapi.Update)*/ {
	msg := tgapi.NewMessage(int64(chatID), message)
	msg.ParseMode = "Markdown"
	_, err := bot.Send(msg)
	logger.LogErrorIf(err)
}

func sendMarkdownMessage(chatID int64, message string, update tgapi.Update) {
	msg := tgapi.NewMessage(chatID, message)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "Markdown"
	_, err := bot.Send(msg)
	logger.LogErrorIf(err)
}

func MakeSubscribesKeyboard(subscribesList []SubscribesList) tgapi.InlineKeyboardMarkup {
	var keyboard = tgapi.NewInlineKeyboardMarkup()
	row := tgapi.NewInlineKeyboardRow()
	for _, v := range subscribesList {
		goodID := strconv.Itoa(v.GoodID)
		row = append(row, tgapi.NewInlineKeyboardButtonData(v.Title, goodID))
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	return keyboard
}
