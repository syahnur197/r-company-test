package router

import (
	"encoding/json"
	"github.com/syahnur197/rakuten/rakuten"
	"net/http"
	"strings"
	"time"
)

type Router struct {
	H *rakuten.Handler
}

func NewRouter(h *rakuten.Handler) *Router {
	return &Router{H: h}
}

func (router *Router) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ping": "pong"}`))
}

func (router *Router) GetCurrencyRate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	date := strings.TrimPrefix(r.URL.Path, "/rates/")

	req := &rakuten.GetCurrencyRateRequest{}

	if date == "" {
		notFound(w)
		return
	} else if date == "latest" {
		req.GetLatestDate = true
	} else {
		t, err := time.Parse("2006-01-02", date)
		if err != nil {
			notFound(w)
			return
		}

		req.Date = t
	}

	rates, err := router.H.GetCurrencyRate(ctx, req)

	ratesResponseJson, err := json.Marshal(rates)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(ratesResponseJson)
}

func (router *Router) GetAnalyzedCurrencyRate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	rates, err := router.H.GetAnalyzedCurrencyRate(ctx)
	if err != nil {
		internalError(w)
		return
	}

	ratesResponseJson, err := json.Marshal(rates)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(ratesResponseJson)
}

type ErrorResponse struct {
	Message string `json:"message"`
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
