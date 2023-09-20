package main

import (
	"log"
	"nnyd-back/config"
	"nnyd-back/db"
	"time"

	"github.com/golang-jwt/jwt"
)

func get_token(user_id string) string {
	claims := jwt.MapClaims{
		"user_id": user_id,
		"exp":     time.Now().Add(3000 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString([]byte(config.JST_SECRET_KEY))
	if err != nil {
		log.Println(err)
	}

	return accessToken
}

func main() {
	config.LoadConfig()

	db.Init()
	defer db.Close()
	db.AutoMigration()

	log.Println(get_token("3d5cfdb2-054d-49c4-9b7f-6d085dc9c701"))
	log.Println(get_token("3d5cfdb2-054d-49c4-9b7f-6d085dc9c702"))
	log.Println(get_token("3d5cfdb2-054d-49c4-9b7f-6d085dc9c703"))
}
