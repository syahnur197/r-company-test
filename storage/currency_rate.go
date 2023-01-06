package storage

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
)

const (
	createCurrencyRateSql = `
		INSERT INTO currency_rate (
			base,
		  	quote, 
			rate, 
			published_date
		) VALUES (
			:base,
		  	:quote, 
			:rate, 
			:published_date
		) RETURNING id;
	`

	getCurrencyRateSql = `
		SELECT 
		    base, 
		    quote, 
		    TRIM(TRAILING '0' FROM CAST(rate AS TEXT)) as rate, 
		    published_date 
		FROM currency_rate
	`

	getLatestCurrencyRateDateSql = `
		SELECT MAX(published_date) FROM currency_rate LIMIT 1
	`

	getAnalyzedCurrencyRateSql = `
		SELECT 
			base, 
			quote, 
			TRIM(TRAILING '0' FROM CAST(MIN(rate) AS TEXT)) as min, 
			TRIM(TRAILING '0' FROM CAST(MAX(rate) AS TEXT)) as max, 
			AVG(rate) as avg
		FROM currency_rate
		GROUP BY base, quote
	`
)

func (s *Storage) CreateCurrencyRate(ctx context.Context, rate Rate) (string, error) {
	var id string
	nstmt, err := s.db.PrepareNamedContext(ctx, createCurrencyRateSql)
	if err != nil {
		return "", errors.Wrap(err, "failed to prepared name context")
	}
	defer nstmt.Close()
	if err := nstmt.QueryRowContext(ctx, rate).Scan(&id); err != nil {
		return "", errors.Wrap(err, "failed to create currency rate")
	}
	return id, nil
}

func (s *Storage) GetCurrencyRates(ctx context.Context, filter CurrencyFilter) ([]Rate, error) {
	var rates []Rate

	query := ""

	if filter.Date.IsZero() && filter.GetLatestDate {
		query = fmt.Sprintf(`%s WHERE published_date = (%s)`, getCurrencyRateSql, getLatestCurrencyRateDateSql)
	}

	params := map[string]interface{}{}

	if !filter.Date.IsZero() {
		query = fmt.Sprintf(`%s WHERE published_date = :published_date`, getCurrencyRateSql)
		params["published_date"] = filter.Date
	}

	nstmt, err := s.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to prepare statement for retrieving currency rates")
	}
	defer nstmt.Close()
	if err = nstmt.SelectContext(ctx, &rates, params); err != nil {
		return nil, errors.Wrap(err, "failed to retrieve currency rates")
	}
	return rates, nil
}

func (s *Storage) GetAnalyzedCurrencyRates(ctx context.Context) ([]AnalyzedRate, error) {
	var rates []AnalyzedRate

	nstmt, err := s.db.PrepareNamedContext(ctx, getAnalyzedCurrencyRateSql)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to prepare statement for retrieving analyzed currency rates")
	}
	defer nstmt.Close()
	if err = nstmt.SelectContext(ctx, &rates, map[string]interface{}{}); err != nil {
		return nil, errors.Wrap(err, "failed to retrieve analyzed currency rates")
	}
	return rates, nil
}
