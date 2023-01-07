package storage

func (s *Storage) CreateCurrencyRatesTable() error {
	sql := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	DROP TABLE IF EXISTS currency_rate;
	CREATE TABLE IF NOT EXISTS currency_rate (
    	"id" UUID DEFAULT uuid_generate_v1() PRIMARY KEY,
    	"base" VARCHAR(3) NOT NULL,	
    	"quote" VARCHAR(3) NOT NULL,
    	"rate" NUMERIC(20,10) NOT NULL,
    	"published_date" DATE NOT NULL 
	);`

	_, err := s.db.Exec(sql)
	return err
}
