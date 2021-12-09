package db

import (
	"log"

	"gopkg.in/ini.v1"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	cfg, err := ini.Load("config/config.ini")
	if err != nil {
		log.Fatal(err)
	}

	USER := cfg.Section("db").Key("user").String()
	PASS := cfg.Section("db").Key("password").String()
	DBNAME := cfg.Section("db").Key("database_name").String()
	CONNECT := USER + ":" + PASS + "@/" + DBNAME + "?parseTime=true"
	db, err := gorm.Open(mysql.Open(CONNECT), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	DB = db
}