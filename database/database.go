package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

var ORM, Errs = GormInit()

// GormInit init gorm ORM.
func GormInit() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "./portfolio.db")
	if err != nil {
		log.Panic(err)
	}
	return db, err
}