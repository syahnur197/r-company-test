package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

type Rate struct {
	Base  string `db:"base"`
	Quote string `db:"quote"`
	Rate  string `db:"rate"`
	Date  string `db:"published_date"`
}

type AnalyzedRate struct {
	Base  string `db:"base"`
	Quote string `db:"quote"`
	Min   string `db:"min"`
	Max   string `db:"max"`
	Avg   string `db:"avg"`
}

type RakutenStore interface {
	CreateCurrencyRatesTable() error

	CreateCurrencyRate(ctx context.Context, rate Rate) (string, error)
	GetCurrencyRates(ctx context.Context, filter CurrencyFilter) ([]Rate, error)
	GetAnalyzedCurrencyRates(ctx context.Context) ([]AnalyzedRate, error)
}

type CurrencyFilter struct {
	Date          time.Time
	GetLatestDate bool
}

var (
	_ RakutenStore = (*Storage)(nil)
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		db: db,
	}
}
