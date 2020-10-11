package controllers

import (
	"avito/DBModule"
	"avito/config"
	"avito/models"
	"avito/responses"
	"avito/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/tebeka/selenium"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var wd selenium.WebDriver

func GetSubscribes(response http.ResponseWriter, request *http.Request) { // тут должна быть авторизация
	DB := DBModule.GetDB()
	email := utils.GetRequestParams(request.URL.Query())["email"]
	query := `select get_subscribes(email_in := $1);`
	var data json.RawMessage
	err := DB.Get(&data, query, email)
	if err != nil {
		responses.Error(response, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(response, http.StatusOK, data)
}

func Unsubscribe(response http.ResponseWriter, request *http.Request) {
	DB := DBModule.GetDB()
	var params models.ManualCheckParams
	err := json.NewDecoder(request.Body).Decode(&params) // заполняю структуру с параметрами запроса
	if err != nil {
		responses.Error(response, http.StatusPreconditionFailed, err)
		return
	}
	err = params.Validate()
	if err != nil { // если параметры не прошли валидацию, отправлю 412 код
		responses.Error(response, http.StatusPreconditionFailed, err)
		return
	}
	query := `select * from unsubscribe(user_id_in := $1, good_id_in := $2);`
	_, err = DB.Exec(query, params.UserID, params.GoodID)
	if err != nil {
		responses.Error(response, http.StatusInternalServerError, err)
		return
	}
}

func ManualCheckPrice(response http.ResponseWriter, request *http.Request) {
	DB := DBModule.GetDB()
	var params models.ManualCheckParams
	err := json.NewDecoder(request.Body).Decode(&params) // заполняю структуру с параметрами запроса
	if err != nil {
		responses.Error(response, http.StatusPreconditionFailed, err)
		return
	}
	err = params.Validate()
	if err != nil { // если параметры не прошли валидацию, отправлю 412 код
		responses.Error(response, http.StatusPreconditionFailed, err)
		return
	}

	query := `select manual_check_validation(user_id_in := $1, good_id_in := $2, cooldown_in_minutes := $3);`
	var dbResponse string
	err = DB.Get(&dbResponse, query, params.UserID, params.GoodID, config.ManualCheckCooldownInMinutes)
	if err != nil {
		responses.Error(response, http.StatusInternalServerError, err)
		return
	}
	var responseMessage models.ResponseMessage
	if !strings.Contains(dbResponse, "https://www.avito") {
		responseMessage.Message = dbResponse
		responses.JSON(response, http.StatusForbidden, responseMessage)
		return
	}

	price, _, err := checkPriceByURL(dbResponse) // получаю цену, ошибку
	if err != nil {
		responses.Error(response, http.StatusPreconditionFailed, err)
		return
	}
	responseMessage.Message = strconv.Itoa(price)
	responses.JSON(response, http.StatusOK, responseMessage)
	go updateLastCheck(DB, int(params.UserID), int(params.GoodID), price) // обновляю в горутине, чтобы сократить время запроса
}

func updateLastCheck(DB *sqlx.DB, userID, goodID, price int) {
	query := `select update_last_check(user_id_in := $1, good_id_in := $2, price_in := $3);`
	_, err := DB.Exec(query, userID, goodID, price)
	utils.LogErrorIf(err)
}

func Subscribe(response http.ResponseWriter, request *http.Request) { // подписка на товар
	DB := DBModule.GetDB()
	var params models.SubscribeParams
	err := json.NewDecoder(request.Body).Decode(&params) // заполняю структуру с параметрами запроса
	if err != nil {
		responses.Status(response, http.StatusPreconditionFailed)
		return
	}
	err = params.Validate()
	if err != nil { // если параметры не прошли валидацию, отправлю 412 код
		responses.Error(response, http.StatusPreconditionFailed, err)
		return
	}

	price, _, err := checkPriceByURL(params.GoodURL) // получаю цену, ошибку
	if err != nil {
		responses.Error(response, http.StatusPreconditionFailed, err)
		return
	}

	tx, err := DB.Beginx()
	if err != nil {
		responses.Error(response, http.StatusInternalServerError, err)
		return
	}

	query := `select subscribe(email_in := $1, good_url_in := $2, price_in := $3);`
	var DBResponse string
	err = tx.Get(&DBResponse, query, // вызываю хранимую функцию
		params.Email,
		params.GoodURL,
		price)
	if err != nil {
		responses.Error(response, http.StatusInternalServerError, err)
		_ = tx.Rollback()
		return
	}

	var token string
	if DBResponse == "new user registered" {
		query := `select email_confirmation_token_create(email_in := $1);`
		err = tx.Get(&token, query, params.Email) // вызываю хранимую функцию
		if err != nil {
			responses.Error(response, http.StatusInternalServerError, err)
			_ = tx.Rollback()
			return
		}
		err := SendConfirmationToken(params.Email, token)
		if err != nil {
			responses.Error(response, http.StatusInternalServerError, err)
			_ = tx.Rollback()
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		responses.Error(response, http.StatusInternalServerError, err)
		return
	}

	response.WriteHeader(http.StatusCreated)
}

func EmailConfirmation(response http.ResponseWriter, request *http.Request) {
	DB := DBModule.GetDB()
	token := utils.GetRequestParams(request.URL.Query())["token"]

	var dbResult string
	query := `select user_confirm_email(token_in := $1);`
	err := DB.Get(&dbResult, query, token)
	if err != nil {
		responses.Error(response, http.StatusInternalServerError, err)
		return
	}
	if dbResult != "" {
		responses.Status(response, http.StatusNotFound)
		return
	}
}

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

func checkPageExists(title string) bool {
	if title == "Ошибка 404. Страница не найдена — Объявления на сайте Авито" {
		return false
	}
	return true
}

func startPriceChecker() {
	DB := DBModule.GetDB()
	t := time.NewTicker(time.Second * time.Duration(config.ParseIntervalInSeconds))
	for {
		var urlsForParse []models.GoodsForCheck
		query := `select * from urls_for_parse_select();`
		err := DB.Select(&urlsForParse, query)
		utils.LogErrorIf(err)

		for _, good := range urlsForParse {
			checkedPrice, title, err := checkPriceByURL(good.Url)
			if err != nil {
				<-t.C
			}
			utils.LogErrorIf(err)
			if checkedPrice != good.Price {
				var emailsForNotify []string
				query := `select * from get_emails_for_notify(good_id_in := $1, price_in := $2);`
				err := DB.Select(&emailsForNotify, query, good.GoodID, checkedPrice)
				if len(emailsForNotify) == 0 {
					continue
				}
				utils.LogErrorIf(err)
				//fmt.Println(urlsForParse)
				checkedPriceFormatted := utils.PriceBeautify(strconv.Itoa(checkedPrice))
				oldPriceFormatted := utils.PriceBeautify(strconv.Itoa(good.Price))
				//fmt.Println(emailsForNotify, checkedPriceFormatted, oldPriceFormatted, good.Url, title)
				err = SendNotifyEmail(emailsForNotify, checkedPriceFormatted, oldPriceFormatted, good.Url, title)

				utils.LogErrorIf(err)
				query = `select update_price(good_id_in := $1, price_in := $2);`
				_, err = DB.Exec(query, good.GoodID, checkedPrice)
				utils.LogErrorIf(err)
			}
		}
		<-t.C
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
	utils.LogFatalIf(err)
	//defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	utils.LogFatalIf(err)
	fmt.Printf("\nSelenium successfully initialized\n")
	go startPriceChecker()
}
