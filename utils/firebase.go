package utils

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var (
	client *auth.Client
)

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
