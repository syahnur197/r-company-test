package router

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/syahnur197/rakuten/rakuten"
)

type Router struct {
	H *rakuten.Handler
}

func NewRouter(h *rakuten.Handler) *Router {
	return &Router{H: h}
}

func (rtr *Router) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ping": "pong"}`))
}

func (rtr *Router) GetCurrencyRate(w http.ResponseWriter, r *http.Request) {
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
		// validate date format
		t, err := time.Parse("2006-01-02", date)
		if err != nil {
			badRequest(w, "invalid date format, must be YYYY-MM-DD")
			return
		}

		req.Date = t
	}

	rates, err := rtr.H.GetCurrencyRate(ctx, req)
	if err != nil {
		log.Println("failed to obtained currency rates")
		internalError(w)
		return
	}

	ratesResponseJson, err := json.Marshal(rates)
	if err != nil {
		log.Println("failed to marshal rates")
		internalError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(ratesResponseJson)
}

func (rtr *Router) GetAnalyzedCurrencyRate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	rates, err := rtr.H.GetAnalyzedCurrencyRate(ctx)
	if err != nil {
		log.Println("failed to obtained analyzed currency rates")
		internalError(w)
		return
	}

	ratesResponseJson, err := json.Marshal(rates)
	if err != nil {
		log.Println("failed to marshal analyzed rates")
		internalError(w)
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
		// shouldn't happen
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
		// shouldn't happen
		panic(err)
	}
	w.Write(responseJson)
}

func badRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)

	if message == "" {
		message = "bad request"
	}

	response := ErrorResponse{
		Message: message,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		// shouldn't happen
		panic(err)
	}
	w.Write(responseJson)
}
