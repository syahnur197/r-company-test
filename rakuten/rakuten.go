package rakuten

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/syahnur197/rakuten/storage"
)

type Handler struct {
	Storage storage.RakutenStore
}

func NewHandler(s storage.RakutenStore) Handler {
	return Handler{
		Storage: s,
	}
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ping": "pong"}`))
}

func (h *Handler) GetCurrencyRate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	date := strings.TrimPrefix(r.URL.Path, "/rates/")

	filter := storage.CurrencyFilter{}
	if date == "" {
		notFound(w)
		return
	} else if date == "latest" {
		filter.GetLatestDate = true
	} else {
		t, err := time.Parse("2006-01-02", date)
		if err != nil {
			notFound(w)
			return
		}

		filter.Date = t
	}

	rates, err := h.Storage.GetCurrencyRates(ctx, filter)
	if err != nil {
		internalError(w)
		return
	}

	if len(rates) == 0 {
		notFound(w)
		return
	}

	rateResponse := CurrencyRatesResponse{Base: "EUR", Rates: map[string]string{}}
	for _, rate := range rates {
		rateResponse.Rates[rate.Quote] = rate.Rate
	}

	ratesResponseJson, err := json.Marshal(rateResponse)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(ratesResponseJson)
}

func (h *Handler) GetAnalyzedCurrencyRate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	rates, err := h.Storage.GetAnalyzedCurrencyRates(ctx)
	if err != nil {
		internalError(w)
		return
	}

	if len(rates) == 0 {
		notFound(w)
		return
	}

	rateResponse := AnalyzedRatesResponse{Base: "EUR", RatesAnalyzed: map[string]AnalyzedRate{}}
	for _, rate := range rates {
		rateResponse.RatesAnalyzed[rate.Quote] = AnalyzedRate{
			Min: rate.Min,
			Max: rate.Max,
			Avg: rate.Avg,
		}
	}

	ratesResponseJson, err := json.Marshal(rateResponse)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(ratesResponseJson)
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

type ErrorResponse struct {
	Message string `json:"message"`
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

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)

	response := ErrorResponse{
		Message: "not found",
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	w.Write(responseJson)
}

func internalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)

	response := ErrorResponse{
		Message: "internal server error",
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	w.Write(responseJson)
}
