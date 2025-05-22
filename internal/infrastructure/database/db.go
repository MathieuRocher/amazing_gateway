package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
    var err error

	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		/*os.Getenv("MARIADB_USER")*/ "cam1on",
		/*os.Getenv("MARIADB_PASSWORD")*/ "PouetPouet",
		"localhost",
		"3306",
		"amazing",
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("❌ Failed to connect to DB: %v", err))
	}
 
}
