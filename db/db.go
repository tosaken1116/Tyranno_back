package db

import (
	"fmt"
	"log"

	"nnyd-back/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

// Init is initialize db from main function
func Init() (err error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		config.POSTGRES_HOST,
		config.POSTGRES_USER,
		config.POSTGRES_PASSWORD,
		config.POSTGRES_DB,
		"5432",
	)

	db, err = gorm.Open(
		postgres.Open(dsn), &gorm.Config{},
	)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

// GetDB is called in models
func GetDB() *gorm.DB {
	return db
}

// Close is closing db
func Close() (err error) {
	db, _ := db.DB()
	if err := db.Close(); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func AutoMigration() (err error) {
	if err = db.AutoMigrate(&Users{}, &Posts{}, &Favorites{}, &Follows{}); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
