package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBName     string
	DBPassword string
	DBParams   string

	S3Bucket string
	S3Secret string
	S3ID     string
	S3Url    string
	S3Region string

	JWTSecret   string
	BCRYPT_Salt int
	JWTExp      int
}

func Get() (*Config, error) {

	var Conf *Config
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	JWTExp, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION"))
	if err != nil {
		log.Println("Error parsing JWT_EXPIRATION, Setting JWT_EXPIRATION to 60")
		JWTExp = 60
	}

	salt, err := strconv.Atoi(os.Getenv("BCRYPT_SALT"))
	if err != nil {
		log.Println("Error parsing BCRYPT_SALT, Setting SALT to 8")
		salt = 8
	}

	S3Bucket := os.Getenv("S3_BUCKET_NAME")
	if S3Bucket == "" {
		log.Println("S3_BUCKET_NAME is empty")
	}

	S3Region := os.Getenv("S3_REGION")
	if S3Region == "" {
		log.Println("S3_REGION is empty")
	}

	s3Url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com", S3Bucket, os.Getenv("S3_REGION"))

	Conf = &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USERNAME"),
		DBName:     os.Getenv("DB_NAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBParams:   os.Getenv("DB_PARAMS"),

		S3Secret: os.Getenv("S3_SECRET_KEY"),
		S3ID:     os.Getenv("S3_ID"),
		S3Bucket: S3Bucket,
		S3Url:    s3Url,
		S3Region: S3Region,

		JWTSecret:   os.Getenv("JWT_SECRET"),
		BCRYPT_Salt: salt,
		JWTExp:      JWTExp,
	}

	return Conf, nil
}
