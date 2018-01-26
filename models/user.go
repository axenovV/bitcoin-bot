package models

import "github.com/jinzhu/gorm"

type User struct {

	gorm.Model

	Portfolios []Portfolio

	Num          int     `gorm:"AUTO_INCREMENT"`

}