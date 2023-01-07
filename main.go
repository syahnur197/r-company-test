package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/syahnur197/rakuten/rakuten"
	"github.com/syahnur197/rakuten/router"
	"github.com/syahnur197/rakuten/storage"
)

const (
	dbUser = "rakuten"
	dbPass = "rakuten"
	dbName = "rakuten"
)

var (
	dbPort = "5555"
	dbHost = "localhost"
)

func main() {
	if os.Getenv("DB_PORT") != "" {
		dbPort = os.Getenv("DB_PORT")
	}
	if os.Getenv("DB_HOST") != "" {
		dbHost = os.Getenv("DB_HOST")
	}

	// setting up db
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	s := storage.NewStorage(db)

	h := rakuten.NewHandler(s)

	// setup database schema
	log.Println("initialise schema")
	err = s.CreateCurrencyRatesTable()
	if err != nil {
		log.Fatal(err)
	}

	// fetch currency ratesList
	log.Println("fetching currency rates")
	ratesList, err := rakuten.FetchCurrencyRates()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("storing currency rates")
	for _, rate := range ratesList.Rates {
		rate.Base = "EUR"

		storageRate, err := rakuten.ConvertToStoreRate(rate)
		if err != nil {
			log.Fatal(err)
		}

		_, err = s.CreateCurrencyRate(context.Background(), storageRate)
		if err != nil {
			log.Fatal(err)
		}
	}

	// setting up mux
	log.Println("setting up mux")
	mux := http.NewServeMux()

	r := router.NewRouter(h)

	mux.HandleFunc("/ping", r.Ping)
	mux.HandleFunc("/rates/analyze", r.GetAnalyzedCurrencyRate)
	mux.HandleFunc("/rates/", r.GetCurrencyRate)

	log.Println("listening to port :4000")
	err = http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
