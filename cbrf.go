package cbrf

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
)

const (
	CurrencyRUB = "RUB"
	CurrencyUSD = "USD"
	CurrencyEUR = "EUR"
	CurrencyCNY = "CNY"
	CurrencyHKD = "HKD"
)

type resultType struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valute  []struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Nominal  string `xml:"Nominal"`
		Name     string `xml:"Name"`
		Value    string `xml:"Value"`
	} `xml:"Valute"`
}

type cacheKeyType struct {
	CurrencyId string
	Date       string
}

type cacheResultType struct {
	Rate float64
}

var urlTemplate string = "https://www.cbr.ru/scripts/XML_daily.asp?date_req=%s"

var cache map[cacheKeyType]*cacheResultType

func GetExchangeRate(currencyId string, date time.Time) (float64, error) {

	if cache == nil {
		cache = map[cacheKeyType]*cacheResultType{}
	}

	reqDate := fmt.Sprintf("%02d/%02d/%d", date.Day(), date.Month(), date.Year())

	cacheKey := cacheKeyType{CurrencyId: currencyId, Date: reqDate}

	cacheResult, exists := cache[cacheKey]

	if !exists {

		url := fmt.Sprintf(urlTemplate, reqDate)

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return 0, err
		}

		req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Header.Add("accept-encoding", "gzip, deflate, br")
		req.Header.Add("accept-language", "en-US,en;q=0.9,ru;q=0.8")
		req.Header.Add("cache-control", "max-age=0")
		req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")

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

		result := &resultType{}

		err = xml.Decode(result)
		if err != nil {
			return 0, err
		}

		for _, resultRow := range result.Valute {

			if resultRow.CharCode != currencyId {
				continue
			}

			resultRow.Value = strings.Replace(resultRow.Value, ",", ".", 1)

			rate, err := strconv.ParseFloat(resultRow.Value, 64)
			if err != nil {
				return 0, err
			}

			nominal, err := strconv.ParseInt(resultRow.Nominal, 10, 64)
			if err != nil {
				return 0, err
			}

			cacheResult = &cacheResultType{Rate: rate / float64(nominal)}

			cache[cacheKey] = cacheResult

		}

	}

	return cacheResult.Rate, nil

}

func Convert(from string, to string, value float64, date time.Time) (float64, error) {

	if from == to {
		return value, nil
	}

	if value == 0 {
		return 0, nil
	}

	result := value

	if from != CurrencyRUB {

		exchangeRate, err := GetExchangeRate(from, date)
		if err != nil {
			return 0, err
		}

		result = result * exchangeRate

	}

	if to != CurrencyRUB {

		exchangeRate, err := GetExchangeRate(to, date)
		if err != nil {
			return 0, err
		}

		result = result / exchangeRate

	}

	return (math.Floor(result*100) / 100), nil

}
