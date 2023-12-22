package db

import (
	"fmt"

	"github.com/tahamazari/outpatient_server/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var database *gorm.DB
var e error

func DatabaseInit() {
	host := "localhost"
	user := "postgres"
	password := "1234"
	dbName := "outpatient_db"
	port := 5432

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbName, port)
	database, e = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err := database.AutoMigrate(
		&models.Employee{},
		&models.Dependent{},
		&models.BillingClaim{},
		&models.MedicalBill{},
	); err != nil {
		panic(err)
	}

	if e != nil {
		panic(e)
	}
}

func DB() *gorm.DB {
	return database
}
