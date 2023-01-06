package rakuten

import (
	"context"
	"encoding/xml"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"time"

	"github.com/syahnur197/rakuten/storage"
)

type Handler struct {
	Storage storage.RakutenStore
}

func NewHandler(s storage.RakutenStore) *Handler {
	return &Handler{
		Storage: s,
	}
}

type GetCurrencyRateRequest struct {
	GetLatestDate bool
	Date          time.Time
}

func (h *Handler) GetCurrencyRate(ctx context.Context, req *GetCurrencyRateRequest) (*CurrencyRatesResponse, error) {
	filter := storage.CurrencyFilter{}

	if req.GetLatestDate {
		filter.GetLatestDate = true
	} else if !req.Date.IsZero() {
		filter.Date = req.Date
	}

	rates, err := h.Storage.GetCurrencyRates(ctx, filter)
	if err != nil {
		return nil, err
	}

	rateResponse := CurrencyRatesResponse{Base: "EUR", Rates: map[string]string{}}
	for _, rate := range rates {
		rateResponse.Rates[rate.Quote] = rate.Rate
	}

	return &rateResponse, nil
}

func (h *Handler) GetAnalyzedCurrencyRate(ctx context.Context) (*AnalyzedRatesResponse, error) {
	rates, err := h.Storage.GetAnalyzedCurrencyRates(ctx)
	if err != nil {
		return nil, err
	}
	if len(rates) == 0 {
		return nil, errors.Wrap(err, "zero rates")
	}

	rateResponse := AnalyzedRatesResponse{Base: "EUR", RatesAnalyzed: map[string]AnalyzedRate{}}
	for _, rate := range rates {
		rateResponse.RatesAnalyzed[rate.Quote] = AnalyzedRate{
			Min: rate.Min,
			Max: rate.Max,
			Avg: rate.Avg,
		}
	}

	return &rateResponse, nil
}

func FetchCurrencyRates() (Rates, error) {
	url := "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
	v := Rates{}

	resp, err := http.Get(url)
	if err != nil {
		return v, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return v, err
	}

	if err := xml.Unmarshal(body, &v); err != nil {
		return v, err
	}

	return v, nil
}

type CurrencyRatesResponse struct {
	Base  string            `json:"base"`
	Rates map[string]string `json:"rates"`
}

type AnalyzedRate struct {
	Min string `json:"min"`
	Max string `json:"max"`
	Avg string `json:"avg"`
}

type AnalyzedRatesResponse struct {
	Base          string                  `json:"base"`
	RatesAnalyzed map[string]AnalyzedRate `json:"rates_analyze"`
}
