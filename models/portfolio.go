package models

type Portfolio struct {

	UserId int `json:"user_id",gorm:"index"`
	Name string  `json:"name",gorm:"type:varchar(100)"`
	Currencies []CurrencyInPortfolio `json:"currencies"`
}

type CurrencyInPortfolio struct {

	Code string `json:"code",gorm:"type:varchar(30)"`
	Count int `json:"count"`

}

