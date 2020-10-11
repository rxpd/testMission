package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	//"github.com/google/uuid"
	//"github.com/google/uuid"
)

type SubscribeParams struct {
	GoodURL string `json:"goodURL"`
	Email   string `json:"email"`
}

func (s SubscribeParams) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.GoodURL, validation.Required, is.URL),
		validation.Field(&s.Email, validation.Required, is.Email),
	)
}

type ManualCheckParams struct {
	GoodID uint `json:"goodID"`
	UserID uint `json:"userID"`
}

func (m ManualCheckParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.GoodID, validation.Required),
		validation.Field(&m.UserID, validation.Required),
	)
}

type SubscribeDBMessage struct {
	Message string `db:"message"`
	Token   string `db:"token_uuid"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

type GoodsForCheck struct {
	Url    string `db:"url"`
	Price  int    `db:"price"`
	Email  string `db:"email"`
	UserID string `db:"user_id"`
	GoodID string `db:"good_id"`
	//Title  string `db:"title"`
}

//type ManualChecks struct {
//	UserID string `db:"user_id"`
//	GoodID string `db:"good_id"`
//	Price  int    `db:"price"`
//}
