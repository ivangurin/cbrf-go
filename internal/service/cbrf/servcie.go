package cbrf_service

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/ivangurin/cbrf-go/internal/model"
	cbrf_client "github.com/ivangurin/cbrf-go/internal/pkg/client/cbrf"
)

type IService interface {
	GetExchangeRate(ctx context.Context, currencyID string, date time.Time) (float64, error)
	Convert(ctx context.Context, from string, to string, value float64, date time.Time) (float64, error)
}

type service struct {
	cbrfClient cbrf_client.IClient
	cache      map[cacheKeyType]*cacheResultType
	mu         sync.Mutex
}

func NewService(
	cbrfClient cbrf_client.IClient,
) IService {
	return &service{
		cbrfClient: cbrfClient,
		cache:      map[cacheKeyType]*cacheResultType{},
		mu:         sync.Mutex{},
	}
}

func (s *service) GetExchangeRate(ctx context.Context, currencyID string, date time.Time) (float64, error) {
	if currencyID == model.CurrencyRUB {
		return 1, nil
	}

	cacheKey := cacheKeyType{CurrencyID: currencyID, Date: fmt.Sprintf("%d/%d/%d", date.Year(), date.Month(), date.Day())}
	cacheResult, exists := s.cache[cacheKey]
	if exists {
		return cacheResult.Rate, nil
	}

	rate, err := s.cbrfClient.GetExchangeRate(ctx, currencyID, date)
	if err != nil {
		return 0, fmt.Errorf("failed to get exchange rate: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache[cacheKey] = &cacheResultType{Rate: rate}

	return rate, nil
}

func (s *service) Convert(ctx context.Context, from string, to string, value float64, date time.Time) (float64, error) {
	if from == to {
		return value, nil
	}
	if value == 0 {
		return 0, nil
	}

	res := value

	if from != model.CurrencyRUB {
		exchangeRate, err := s.GetExchangeRate(ctx, from, date)
		if err != nil {
			return 0, err
		}

		res = res * exchangeRate
	}

	if to != model.CurrencyRUB {
		exchangeRate, err := s.GetExchangeRate(ctx, to, date)
		if err != nil {
			return 0, err
		}

		res = res / exchangeRate
	}

	return (math.Floor(res*100) / 100), nil
}
