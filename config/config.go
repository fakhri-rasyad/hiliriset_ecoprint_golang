package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	APPConfig *Config
)

type Config struct{
	APPPort string
	DBHost string
	DBPort string
	DBUser string
	DBPassword string
	DBName string
	APPUrl string
	JWTSecret string
	JWTExpireMinutes string
	JWTRefreshToken string
}

func LoadEnv(){
	godotenv.Load()

	APPConfig = &Config{
		APPPort: getEnv("PORT", "3000"),
		DBHost : getEnv("DB_HOST", "localhost"),
		DBUser : getEnv("DB_USERNAME", "postrges"),
		DBPassword : getEnv("DB_PASSWORD", "admin"),
		DBName : getEnv("DB_NAME", "ecoprint_golang"),
		APPUrl : getEnv("APP_URL", "http://localhost:3000"),
		JWTSecret : getEnv("JWT_SECRET", "supersecret"),
		JWTExpireMinutes: getEnv("JWT_EXPIRY_MINUTES", "1800"),
		JWTRefreshToken: getEnv("REFRESH_TOKEN_EXPIRED", "24h"),
	}


}


func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)

	if exists{
		return value
	} else {
		return fallback
	}
}

func connectToDB(){
	cfg := APPConfig
	dsn := fmt.Sprintf("host=%s post=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.APPPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatal("Failed to get database", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

}