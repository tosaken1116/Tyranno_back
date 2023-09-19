package config

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type ContextValueKey string

const (
	FIREBASE_ID = ContextValueKey("firebase_id")
	USER_ID     = ContextValueKey("user_id")
)

var (
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_HOST     string
	POSTGRES_DB       string
	ENV               string
	PORT              string
	APP_NAME          string
	JST_SECRET_KEY    string

	client *auth.Client
)

func LoadConfig() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("failed load env")
	}
	POSTGRES_USER = os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_HOST = os.Getenv("POSTGRES_HOST")
	POSTGRES_DB = os.Getenv("POSTGRES_DB")
	ENV = os.Getenv("ENV")
	PORT = os.Getenv("PORT")
	APP_NAME = os.Getenv("APP_NAME")
	JST_SECRET_KEY = os.Getenv("JST_SECRET_KEY")
}

func InitFirebaseAuthClient() {
	opt := option.WithCredentialsFile("nnyd.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal("failed load firebase setting json")
	}
	client, err = app.Auth(context.Background())
	if err != nil {
		log.Fatal("failed get client")
	}
}

func GetFirebaseAuthClient() *auth.Client {
	return client
}
