package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	HOST := os.Getenv("HOST")
	DBUSER := os.Getenv("DBUSER")
	DBPASSWORD := os.Getenv("DBPASSWORD")
	DBNAME := os.Getenv("DBNAME")
	PORT := os.Getenv("PORT")
	SSLMODE := os.Getenv("SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", HOST, DBUSER, DBPASSWORD, DBNAME, PORT, SSLMODE)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	err = db.AutoMigrate(
		&Node{},
		&Message{},
		&QuickReply{},
		&KeywordDecision{},
		&TagDecision{},
		&Random{},
		&Tag{},
		&UserSession{},
		&FirstStep{},
		&LineBotChannelSetting{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database schema: %v", err)
	}

	return db
}
