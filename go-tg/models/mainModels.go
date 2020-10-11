package models

type GoodsForCheck struct {
	Url    string `db:"url"`
	Price  int    `db:"price"`
	ChatID string `db:"chat_id"`
	GoodID string `db:"good_id"`
}
