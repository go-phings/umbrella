package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-phings/umbrella"
)

const dbDSN = "host=localhost user=protouser password=protopass port=54320 dbname=protodb sslmode=disable"
const tblPrefix = "p_"

func main() {
	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		log.Fatal("Error connecting to db")
	}

	// create umbrella controller
	u := *umbrella.NewUmbrella(db, tblPrefix, &umbrella.JWTConfig{
		Key:               "someSecretKey",
		Issuer:            "someIssuer",
		ExpirationMinutes: 15,
	}, &umbrella.UmbrellaConfig{
		TagName:           "ui",
	})

	// create database tables
	err = u.CreateDBTables()
	if err != nil {
		log.Fatalf("error creating database tables: %s", err.Error())
	}

	// create admin user
	key, err := u.CreateUser("admin@example.com", "admin", map[string]string{
		"Name": "admin",
	})
	if err != nil {
		log.Fatalf("error with creating admin: %s", err.Error())
	}
	err = u.ConfirmEmail(key)
	if err != nil {
		log.Fatalf("error with confirming admin email: %s", err.Error())
	}

	// /umbrella/{login,logout,register,confirm}
	http.Handle("/umbrella/", u.GetHTTPHandler("/umbrella/"))

	log.Fatal(http.ListenAndServe(":8001", nil))
}
