package databases

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

var Db *gorm.DB

func init() {
	var err error
	Db, err = gorm.Open("sqlite3", "./data.db")
	if err != nil {
		log.Panicln("err", err.Error())
	}
}
