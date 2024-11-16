package cbrf

import (
	"context"
	"time"

	"github.com/ivangurin/cbrf-go/internal/pkg/service_provider"
)

var service = service_provider.GetServiceProvider().GetCbrfService()

func GetExchangeRate(ctx context.Context, currencyID string, date time.Time) (float64, error) {
	rate, err := service.GetExchangeRate(ctx, currencyID, date)
	if err != nil {
		return 0, err
	}

	return rate, nil
}

func Convert(ctx context.Context, from string, to string, value float64, date time.Time) (float64, error) {
	value, err := service.Convert(ctx, from, to, value, date)
	if err != nil {
		return 0, err
	}

	return value, nil
}
