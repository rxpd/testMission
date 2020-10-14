package main

import (
	"avitoTelegram/config"
	"avitoTelegram/utils"
	"avitoTelegram/utils/logger"
	msgFormater "avitoTelegram/utils/messageFormater"
	"errors"
	"fmt"
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/tebeka/selenium"
	"strconv"
	"time"
)

var wd selenium.WebDriver

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
	t := time.NewTicker(time.Second * time.Duration(config.ParseIntervalInSeconds))
	for {
		var urlsForParse []GoodsForCheck
		query := `select * from select_urls_for_parse();`
		err := DB.Select(&urlsForParse, query)
		logger.LogErrorIf(err)

		for _, good := range urlsForParse {
			checkedPrice, title, err := checkPriceByURL(good.URL)
			if err != nil {
				<-t.C
			}
			logger.LogErrorIf(err)
			if checkedPrice != good.Price {
				var usersForNotify []int64
				query := `select * from select_users_for_notify(good_id_in := $1, price_in := $2);`
				err := DB.Select(&usersForNotify, query, good.GoodID, checkedPrice)
				if len(usersForNotify) == 0 {
					continue
				}
				logger.LogErrorIf(err)
				//fmt.Println(urlsForParse)
				checkedPriceFormatted := utils.PriceBeautify(strconv.Itoa(checkedPrice))
				oldPriceFormatted := utils.PriceBeautify(strconv.Itoa(good.Price))
				//fmt.Println(usersForNotify, checkedPriceFormatted, oldPriceFormatted, good.URL, title)
				//err = SendNotifyEmail(usersForNotify, checkedPriceFormatted, oldPriceFormatted, good.URL, title)
				message := msgFormater.PriceChangedFormat(title, good.URL, oldPriceFormatted, checkedPriceFormatted)
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
		msg := tgapi.NewMessage(u, message)
		msg.ParseMode = "Markdown"
		_, err := bot.Send(msg)
		logger.LogErrorIf(err)
	}
}

func initSelenium() {
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
}
