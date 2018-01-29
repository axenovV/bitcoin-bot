package database

import (
	"github.com/asdine/storm"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
)

type User struct {
	UserId int `storm:"id"`
	Nickname string `storm:"index"`
}

var dbStorm, err = storm.Open("my.db")


func SaveMessage(m *tb.Message) {
	user := User{UserId: m.Sender.ID, Nickname: m.Sender.Username}
	err := dbStorm.Save(&user)
	log.Print(err)
}

func GetUniqueUsers() ([]User, error) {
	var users []User
	return users, dbStorm.All(&users)
}
