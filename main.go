package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/syahnur197/rakuten/rakuten"
	"github.com/syahnur197/rakuten/router"
	"github.com/syahnur197/rakuten/storage"
	"log"
	"net/http"
)

const (
	host     = "localhost"
	port     = 5555
	user     = "rakuten"
	password = "rakuten"
	dbname   = "rakuten"
)

func main() {
	// setting up db
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		panic(err)
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
		_, err := s.CreateCurrencyRate(context.Background(), rakuten.ConvertToStoreRate(rate))
		if err != nil {
			log.Fatal(err)
		}
	}

	// setting up mux
	mux := http.NewServeMux()

	router := router.NewRouter(h)

	mux.HandleFunc("/ping", router.Ping)
	mux.HandleFunc("/rates/analyze", router.GetAnalyzedCurrencyRate)
	mux.HandleFunc("/rates/", router.GetCurrencyRate)

	log.Println("listening to port :4000")
	err = http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
