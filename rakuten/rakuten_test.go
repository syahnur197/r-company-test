package rakuten

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/syahnur197/rakuten/storage"
	"github.com/syahnur197/rakuten/storage/mock_storage"
	"testing"
	"time"
)

func TestHandler_GetAnalyzedCurrencyRate(t *testing.T) {
	ctrl := gomock.NewController(t)
	gAny := gomock.Any()

	testData := []storage.AnalyzedRate{
		{
			Base:  "EUR",
			Quote: "BND",
			Min:   "123",
			Max:   "456",
			Avg:   "222",
		},
	}

	storeValid := func(m *mock_storage.MockRakutenStore) {
		m.EXPECT().GetAnalyzedCurrencyRates(gAny).Return(testData, nil)
	}

	mockStore := mock_storage.NewMockRakutenStore(ctrl)
	storeValid(mockStore)

	h := NewHandler(mockStore)

	rates, err := h.GetAnalyzedCurrencyRate(context.Background())
	if err != nil {
		t.Fatal("unexpected err")
	}

	if rates.RatesAnalyzed["BND"].Min != "123" {
		t.Fatal("unexpected min value")
	}
	if rates.RatesAnalyzed["BND"].Max != "456" {
		t.Fatal("unexpected max value")
	}
	if rates.RatesAnalyzed["BND"].Avg != "222" {
		t.Fatal("unexpected avg value")
	}
}

func TestHandler_GetCurrencyRate(t *testing.T) {
	ctrl := gomock.NewController(t)
	gAny := gomock.Any()

	testData := []storage.Rate{
		{
			Base:  "EUR",
			Quote: "AUD",
			Rate:  "123",
			Date:  time.Now().Format("2006-01-02"),
		},
		{
			Base:  "EUR",
			Quote: "BND",
			Rate:  "123",
			Date:  time.Now().Format("2006-01-02"),
		},
	}

	storeValid := func(m *mock_storage.MockRakutenStore) {
		m.EXPECT().GetCurrencyRates(gAny, gAny).Return(testData, nil)
	}

	mockStore := mock_storage.NewMockRakutenStore(ctrl)
	storeValid(mockStore)

	h := NewHandler(mockStore)

	rates, err := h.GetCurrencyRate(context.Background(), &GetCurrencyRateRequest{})
	if err != nil {
		t.Fatal("unexpected err")
	}

	if rates.Rates["AUD"] != "123" {
		t.Fatal("unexpected value")
	}

	if rates.Rates["BND"] != "123" {
		t.Fatal("unexpected value")
	}
}
