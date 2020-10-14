package main

import (
	"avitoTelegram/config"
	"avitoTelegram/utils"
	"avitoTelegram/utils/DBResponses"
	"avitoTelegram/utils/logger"
	msgFormater "avitoTelegram/utils/messageFormater"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

var DB *sqlx.DB
var userManualCheck []UserManualCheck

func GetSubscribesList(chatID int64) ([]SubscribesList, error) {
	var subscribesList []SubscribesList
	query := `select * from select_subscribes_list(chat_id_in := $1);` // собираю все подписки пользователя по chat_id
	err := DB.Select(&subscribesList, query, chatID)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	return subscribesList, nil
}
func IsNewUser(chatID int64, username string) (bool, error) {
	var dbResponse string
	query := `select * from new_message_handler(chat_id_in := $1, username_in := $2);`
	err := DB.Get(&dbResponse, query, chatID, username)
	if err != nil {
		logger.LogError(err)
		return false, err
	}
	if dbResponse == DBResponses.NewUser { // если новый пользователь, отправляю приветственное сообщение
		return true, nil
	}
	return false, nil
}
func Subscribe(url string, chatID int64) string { // подписка на товар
	price, title, err := checkPriceByURL(url) // получаю цену, ошибку
	if err != nil {
		logger.LogError(err)
		return msgFormater.GoodNotFound
	}

	query := `select * from subscribe(chat_id_in := $1, good_url_in := $2, title_in := $3, price_in := $4)`
	var DBResponse string
	err = DB.Get(&DBResponse, query, chatID, url, title, price)
	if err != nil {
		logger.LogError(err)
		return msgFormater.InternalServerError
	}

	if DBResponse == DBResponses.SubscribeAlreadyExists {
		return msgFormater.AlreadySubscribed
	}

	return msgFormater.SuccessfullySubscribedFormat(title, url, utils.PriceBeautify(strconv.Itoa(price)))
}

func ManualCheckPrice(goodID, chatID int) string { // TODO: сделать, чтобы в случае изменения цены после ручной проверки цена поменялась, то уведомить подписанных пользователей
	query := `select * from get_manual_check_info(good_id_in := $1, chat_id_in := $2)`
	var goodInfo GoodInfo
	err := DB.Get(&goodInfo, query, goodID, chatID)
	if err != nil {
		logger.LogError(err)
		return msgFormater.InternalServerError
	}
	go updateUserManualCheck(goodID, chatID, goodInfo.Price)
	if !canBeManualChecked(goodID, chatID) {
		price := utils.PriceBeautify(strconv.Itoa(goodInfo.Price))
		return msgFormater.CooldownLimitFormat(goodInfo.Title, goodInfo.URL, price)
	}
	parsedPrice, title, err := checkPriceByURL(goodInfo.URL) // получаю цену, ошибку
	if err != nil {
		logger.LogError(err)
		return msgFormater.GoodNotFound
	}
	if parsedPrice != goodInfo.Price { // ЕСЛИ ЦЕНА ИЗМЕНИЛАСЬ УВЕДОМИ ВСЕХ ПОЛЬЗОВАТЕЛЕЙ
		oldPrice, newPrice := utils.PriceBeautify(strconv.Itoa(goodInfo.Price)), utils.PriceBeautify(strconv.Itoa(parsedPrice))
		go updateLastManualCheck(chatID, goodID, parsedPrice)
		return msgFormater.PriceChangedFormat(title, goodInfo.URL, oldPrice, newPrice)
	} else {
		return msgFormater.PriceDoesNotChangedFormat(goodInfo.URL, title, utils.PriceBeautify(strconv.Itoa(goodInfo.Price)))
	}
}

func canBeManualChecked(goodID, chatID int) bool {
	for _, checks := range userManualCheck {
		if checks.ChatID == chatID {
			for _, good := range checks.GoodChecks {
				nowWithCooldown := time.Time.Add(good.LastCheckTime, time.Minute*time.Duration(config.ManualCheckCooldownInMinutes))
				if good.GoodID == goodID && time.Time.Before(time.Now(), nowWithCooldown) {
					return false
				} else {
					return true
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

	var userLocalMnlCheck = UserManualCheck{
		ChatID:     chatID,
		GoodChecks: nil,
	}

	var goodManualCheck = GoodManualChecked{
		GoodID:        goodID,
		LastCheckTime: time.Now(),
		LastPrice:     price,
	}
	userLocalMnlCheck.GoodChecks = append(userLocalMnlCheck.GoodChecks, goodManualCheck)
	userManualCheck = append(userManualCheck, userLocalMnlCheck)
}

func Unsubscribe(goodID int, chatID int) string {
	query := `select * from unsubscribe(good_id_in := $1, chat_id_in := $2)`
	var response GoodInfo
	err := DB.Get(&response, query, goodID, chatID)
	if err != nil {
		logger.LogError(err)
		return msgFormater.InternalServerError
	}
	go deleteNonSubscribedGoods()
	return msgFormater.SuccessfulUnsubscribedFormat(response.Title, response.URL, utils.PriceBeautify(strconv.Itoa(response.Price)))
}

func deleteNonSubscribedGoods() {
	query := `select delete_non_subscribed_goods();`
	_, err := DB.Exec(query)
	logger.LogFatalIf(err)
}

func updateLastManualCheck(userID, goodID, price int) {
	query := `select update_last_manual_check(chat_id_in := $1, good_id_in := $2, price_in := $3);`
	_, err := DB.Exec(query, userID, goodID, price)
	logger.LogErrorIf(err)
}
