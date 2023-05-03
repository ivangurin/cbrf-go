package cbrf

import (
	"testing"
	"time"
)

func TestGetExchangeRate(t *testing.T) {

	date := time.Date(2023, time.Month(05), 1, 0, 0, 0, 0, time.UTC)
	result := 80.5093

	// Получение курса валюты на дату
	exchangeRate, err := GetExchangeRate(CurrencyUSD, date)

	if err != nil {
		t.Error(err)
	}

	if exchangeRate != result {
		t.Errorf("Ожидалось %f, но получили %f", result, exchangeRate)
	}

}

func TestConvert(t *testing.T) {

	date := time.Date(2023, time.Month(05), 1, 0, 0, 0, 0, time.UTC)
	result := 91.10

	// Получение курса валюты на дату
	value, err := Convert(CurrencyUSD, CurrencyEUR, float64(100), date)

	if err != nil {
		t.Error(err)
	}

	if value != result {
		t.Errorf("Ожидалось %f, но получили %f", result, value)
	}

}
