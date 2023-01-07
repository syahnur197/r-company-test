package storage

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = "6666"
	user     = "rakuten"
	password = "rakuten"
	dbname   = "rakuten"
)

func SetupTestDb() *sqlx.DB {
	// setting up db
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func TestStorage_CreateCurrencyRate(t *testing.T) {
	db := SetupTestDb()
	defer db.Close()

	s := NewStorage(db)

	// create schema
	err := s.CreateCurrencyRatesTable()
	if err != nil {
		log.Fatal(err)
	}

	date, err := time.Parse("2006-01-02", "2023-01-05")
	if err != nil {
		log.Fatal(err)
	}

	_, err = s.CreateCurrencyRate(context.Background(), Rate{
		Base:  "EUR",
		Quote: "SGD",
		Rate:  "100",
		Date:  date,
	})
	if err != nil {
		log.Fatal(err)
	}

	rates, err := s.GetCurrencyRates(context.Background(), CurrencyFilter{
		Date:          time.Time{},
		GetLatestDate: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	if len(rates) < 1 {
		log.Fatal("rates count is less than 1")
	}
}

func TestStorage_GetCurrencyRates(t *testing.T) {
	db := SetupTestDb()
	defer db.Close()

	s := NewStorage(db)

	// create schema
	err := s.CreateCurrencyRatesTable()
	if err != nil {
		log.Fatal(err)
	}

	date1, err := time.Parse("2006-01-02", "2023-01-05")
	if err != nil {
		log.Fatal(err)
	}
	date2, err := time.Parse("2006-01-02", "2023-01-01")
	if err != nil {
		log.Fatal(err)
	}

	_, err = s.CreateCurrencyRate(context.Background(), Rate{
		Base:  "EUR",
		Quote: "SGD",
		Rate:  "100",
		Date:  date1,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = s.CreateCurrencyRate(context.Background(), Rate{
		Base:  "EUR",
		Quote: "SGD",
		Rate:  "200",
		Date:  date2,
	})
	if err != nil {
		log.Fatal(err)
	}

	rates, err := s.GetCurrencyRates(context.Background(), CurrencyFilter{
		Date:          time.Time{},
		GetLatestDate: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	if len(rates) < 1 {
		log.Fatal("rates count is less than 1")
	}

	if rates[0].Date.Format("2006-01-02") != date1.Format("2006-01-02") && rates[0].Rate != "100" {
		log.Fatal("invalid latest rate")
	}

	rates, err = s.GetCurrencyRates(context.Background(), CurrencyFilter{
		Date:          date2,
		GetLatestDate: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	if len(rates) < 1 {
		log.Fatal("rates count is less than 1")
	}

	if rates[0].Date.Format("2006-01-02") != date2.Format("2006-01-02") && rates[0].Rate != "200" {
		log.Fatal("invalid latest rate")
	}
}

func TestStorage_GetAnalyzedCurrencyRates(t *testing.T) {
	db := SetupTestDb()
	defer db.Close()

	s := NewStorage(db)

	// create schema
	err := s.CreateCurrencyRatesTable()
	if err != nil {
		log.Fatal(err)
	}

	date1, err := time.Parse("2006-01-02", "2023-01-05")
	if err != nil {
		log.Fatal(err)
	}
	date2, err := time.Parse("2006-01-02", "2023-01-01")
	if err != nil {
		log.Fatal(err)
	}

	_, err = s.CreateCurrencyRate(context.Background(), Rate{
		Base:  "EUR",
		Quote: "SGD",
		Rate:  "100",
		Date:  date1,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = s.CreateCurrencyRate(context.Background(), Rate{
		Base:  "EUR",
		Quote: "SGD",
		Rate:  "200",
		Date:  date2,
	})
	if err != nil {
		log.Fatal(err)
	}

	rates, err := s.GetAnalyzedCurrencyRates(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if len(rates) < 1 {
		log.Fatal("rates count is less than 1")
	}

	if rates[0].Min != "100" {
		log.Fatal("unexpected min")
	}
	if rates[0].Max != "200" {
		log.Fatal("unexpected max")
	}
	if rates[0].Avg != "150" {
		log.Fatal("unexpected avg")
	}
}
