package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-phings/umbrella"
	_ "github.com/lib/pq"
)

const dbDSN = "host=localhost user=uuser password=upass port=54321 dbname=udb sslmode=disable"
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

	// secret stuff
	http.Handle("/secret_stuff/", u.GetHTTPHandlerWrapper(secretStuff(), umbrella.HandlerConfig{}))

	log.Fatal(http.ListenAndServe(":8001", nil))
}

func secretStuff() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := umbrella.GetUserIDFromRequest(r)
		switch userID {
		case 1: 
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("SecretStuffOnlyForAdmin"))
		case 0:
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("YouHaveToBeLoggedIn"))
		default:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("SecretStuffForOtherUser"))
		}
	})
}
