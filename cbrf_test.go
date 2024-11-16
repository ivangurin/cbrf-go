package cbrf

import (
	"context"
	"testing"
	"time"

	"github.com/ivangurin/cbrf-go/pkg/model"
	"github.com/stretchr/testify/require"
)

// Получение курса валюты на дату
func TestGetExchangeRate(t *testing.T) {
	ctx := context.Background()
	date := time.Date(2023, time.Month(5), 1, 0, 0, 0, 0, time.UTC)
	result := 80.5093

	exchangeRate, err := GetExchangeRate(ctx, model.CurrencyUSD, date)
	require.NoError(t, err)
	require.Equal(t, result, exchangeRate)
}

// Получение курса валюты на дату
func TestConvert(t *testing.T) {
	ctx := context.Background()
	date := time.Date(2023, time.Month(5), 1, 0, 0, 0, 0, time.UTC)
	result := 91.10

	value, err := Convert(ctx, model.CurrencyUSD, model.CurrencyEUR, float64(100), date)
	require.NoError(t, err)
	require.Equal(t, result, value)
}
