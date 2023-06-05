package configg

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func KoneksiData() *gorm.DB {
	db, err := gorm.Open("mysql", "root:@(localhost:3306)/auth_golang?parseTime=true")

	db.SingularTable(true)

	if err != nil {
		panic("database not connected")
	}

	return db
}
