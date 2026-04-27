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
  MQTTHost string
  MQTTPort string
}

func LoadEnv(){
	godotenv.Load()

	APPConfig = &Config{
		APPPort: getEnv("PORT", "3000"),
		DBHost : getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5432"),
		DBUser : getEnv("DB_USERNAME", "postrges"),
		DBPassword : getEnv("DB_PASSWORD", "admin"),
		DBName : getEnv("DB_NAME", "ecoprint_golang"),
		APPUrl : getEnv("APP_URL", "http://localhost:3000"),
		JWTSecret : getEnv("JWT_SECRET", "supersecret"),
		JWTExpireMinutes: getEnv("JWT_EXPIRY_MINUTES", "1800"),
		JWTRefreshToken: getEnv("REFRESH_TOKEN_EXPIRED", "24h"),
    MQTTHost: getEnv("MQTT_Host", "localhost"),
    MQTTPort : getEnv("MQTT_Port", "1883"),
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

func ConnectToDB(){
	cfg := APPConfig
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

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
