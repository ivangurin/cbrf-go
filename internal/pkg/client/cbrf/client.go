package cbrf_client

import (
	"compress/gzip"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"golang.org/x/text/encoding/charmap"
)

type IClient interface {
	GetExchangeRate(ctx context.Context, currencyID string, date time.Time) (float64, error)
}

type client struct {
	URL string
}

func NewClient() IClient {
	return &client{
		URL: "https://www.cbr.ru/scripts/XML_daily.asp?date_req=%s",
	}
}

func (c *client) GetExchangeRate(ctx context.Context, currencyID string, date time.Time) (float64, error) {
	url := fmt.Sprintf(c.URL, fmt.Sprintf("%02d/%02d/%d", date.Day(), date.Month(), date.Year()))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create http request: %w", err)
	}

	req = req.WithContext(ctx)

	req.Header.Add("accept", acceptHeader)
	req.Header.Add("accept-encoding", acceptEncodingHeader)
	req.Header.Add("accept-language", acceptLanguageHeader)
	req.Header.Add("cache-control", cacheControlHeader)
	req.Header.Add("user-agent", gofakeit.UserAgent())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return 0, err
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	xml := xml.NewDecoder(reader)
	xml.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}

	response := &getExchangeRateResponse{}
	err = xml.Decode(response)
	if err != nil {
		return 0, err
	}

	var res float64
	var found bool
	for _, resultRow := range response.Valute {
		if resultRow.CharCode != currencyID {
			continue
		}

		resultRow.Value = strings.Replace(resultRow.Value, ",", ".", 1)

		rate, err := strconv.ParseFloat(resultRow.Value, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse rate value: %w", err)
		}

		nominal, err := strconv.ParseInt(resultRow.Nominal, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse rate nominal: %w", err)
		}

		if nominal == 0 {
			return 0, fmt.Errorf("nominal is zero for currency: %s", currencyID)
		}

		found = true
		res = rate / float64(nominal)
	}

	if !found {
		return 0, fmt.Errorf("currency %s not found", currencyID)
	}

	return res, nil
}
