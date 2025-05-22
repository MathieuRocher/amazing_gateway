package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Charge le .env si les variables ne sont pas déjà définies
	_ = godotenv.Load()

	user := os.Getenv("MARIADB_USER")
	pass := os.Getenv("MARIADB_PASSWORD")
	dbname := os.Getenv("MARIADB_DATABASE")

	host := os.Getenv("MARIADB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("MARIADB_PORT")
	if port == "" {
		port = "3306"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, pass, host, port, dbname,
	)

	var err error
	retries := 10
	for i := 0; i < retries; i++ {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("✅ Connected to database")
			return
		}
		log.Printf("❌ DB connection failed (attempt %d/%d): %v", i+1, retries, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("❌ Failed to connect to DB after %d attempts: %v", retries, err)
}
