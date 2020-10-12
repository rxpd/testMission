package main

import "time"

type GoodsForCheck struct {
	URL    string `db:"url"`
	Price  int    `db:"price"`
	ChatID string `db:"chat_id"`
	GoodID string `db:"good_id"`
}

type SubscribesList struct {
	GoodID int    `db:"good_id"`
	Title  string `db:"title"`
}

type GoodInfo struct {
	Title string `db:"title_u"`
	URL   string `db:"url_u"`
	Price int    `db:"price_u"`
}

type ManualCheckPriceResponse struct {
	Message string `db:"message"`
	URL     string `db:"url_a"`
	//GoodID   int    `db:"good_id"`
	OldPrice int `db:"old_price"`
}

type UserManualCheck struct {
	ChatID     int
	GoodChecks []GoodManualChecked
}

type GoodManualChecked struct {
	GoodID        int
	LastCheckTime time.Time
	LastPrice     int
}
