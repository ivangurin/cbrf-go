package cbrf_service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ivangurin/cbrf-go/internal/pkg/suite_provider"
	"github.com/ivangurin/cbrf-go/pkg/model"
	"github.com/stretchr/testify/require"
)

func TestGetExchangeRate(t *testing.T) {
	t.Parallel()

	type testCase struct {
		Name       string
		CurrencyID string
		Date       time.Time
		Rate       float64
		Error      error
	}

	testCases := []testCase{
		{
			Name:       "Get rate - Error case",
			CurrencyID: model.CurrencyUSD,
			Date:       time.Now().UTC(),
			Rate:       100,
			Error:      fmt.Errorf("some error"),
		},
		{
			Name:       "Get rate for RUB",
			CurrencyID: model.CurrencyRUB,
			Date:       time.Now().UTC(),
			Rate:       1,
		},
		{
			Name:       "Get rate for USD",
			CurrencyID: model.CurrencyUSD,
			Date:       time.Now().UTC(),
			Rate:       100,
		},
		{
			Name:       "Get rate for EUR",
			CurrencyID: model.CurrencyUSD,
			Date:       time.Now().UTC(),
			Rate:       110,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			sp := suite_provider.NewSuiteProvider()

			if tc.CurrencyID != model.CurrencyRUB {
				sp.GetCbrfClientMock().EXPECT().
					GetExchangeRate(ctx, tc.CurrencyID, tc.Date).
					Return(tc.Rate, tc.Error).
					Once()
			}

			rate, err := sp.GetCbrfService().GetExchangeRate(ctx, tc.CurrencyID, tc.Date)
			if tc.Error != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.Rate, rate)

			sp.GetCbrfClientMock().AssertExpectations(t)
		})
	}
}

func TestConvert(t *testing.T) {
	t.Parallel()

	type testCase struct {
		Name           string
		FromCurrencyID string
		FromRate       float64
		ToCurrencyID   string
		ToRate         float64
		Date           time.Time
		Value          float64
		ValueDest      float64
		Error          error
	}

	testCases := []testCase{
		{
			Name:           "Convert zero value",
			FromCurrencyID: model.CurrencyRUB,
			ToCurrencyID:   model.CurrencyUSD,
			Value:          0,
			ValueDest:      0,
			Date:           time.Now().UTC(),
		},
		{
			Name:           "Convert RUB->RUB",
			FromCurrencyID: model.CurrencyRUB,
			ToCurrencyID:   model.CurrencyRUB,
			Value:          100,
			ValueDest:      100,
			Date:           time.Now().UTC(),
		},
		{
			Name:           "Convert RUB->USD",
			FromCurrencyID: model.CurrencyRUB,
			ToCurrencyID:   model.CurrencyUSD,
			ToRate:         100,
			Value:          1,
			ValueDest:      .01,
			Date:           time.Now().UTC(),
		},
		{
			Name:           "Convert USD->RUB",
			FromCurrencyID: model.CurrencyUSD,
			FromRate:       100,
			ToCurrencyID:   model.CurrencyRUB,
			Value:          1,
			ValueDest:      100,
			Date:           time.Now().UTC(),
		},
		{
			Name:           "Convert USD->USD",
			FromCurrencyID: model.CurrencyUSD,
			ToCurrencyID:   model.CurrencyUSD,
			Value:          1,
			ValueDest:      1,
			Date:           time.Now().UTC(),
		},
		{
			Name:           "Convert USD->EUR",
			FromCurrencyID: model.CurrencyUSD,
			FromRate:       100,
			ToCurrencyID:   model.CurrencyEUR,
			ToRate:         110,
			Value:          100,
			ValueDest:      90.90,
			Date:           time.Now().UTC(),
		},
		{
			Name:           "Convert - Error case",
			FromCurrencyID: model.CurrencyRUB,
			ToCurrencyID:   model.CurrencyUSD,
			Value:          100,
			Date:           time.Now().UTC(),
			Error:          fmt.Errorf("some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			sp := suite_provider.NewSuiteProvider()

			if tc.FromCurrencyID != tc.ToCurrencyID && tc.Value != 0 {
				if tc.FromCurrencyID != model.CurrencyRUB {
					sp.GetCbrfClientMock().EXPECT().
						GetExchangeRate(ctx, tc.FromCurrencyID, tc.Date).
						Return(tc.FromRate, tc.Error)
				}

				if tc.ToCurrencyID != model.CurrencyRUB {
					sp.GetCbrfClientMock().EXPECT().
						GetExchangeRate(ctx, tc.ToCurrencyID, tc.Date).
						Return(tc.ToRate, tc.Error)
				}
			}

			value, err := sp.GetCbrfService().Convert(ctx, tc.FromCurrencyID, tc.ToCurrencyID, tc.Value, tc.Date)
			if tc.Error != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.ValueDest, value)

			sp.GetCbrfClientMock().AssertExpectations(t)
		})
	}
}
