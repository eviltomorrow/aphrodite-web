package model

import (
	"log"

	"github.com/eviltomorrow/aphrodite-web/db"
)

func init() {
	db.MySQLDSN = "root:root@tcp(localhost:3306)/aphrodite?charset=utf8mb4&parseTime=true&loc=Local"
	db.MySQLMaxOpen = 10
	db.MySQLMinOpen = 5
	db.BuildMySQL()
	log.Printf("model init function\r\n")
}
