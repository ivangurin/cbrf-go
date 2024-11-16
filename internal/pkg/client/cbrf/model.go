package cbrf_client

import "encoding/xml"

type getExchangeRateResponse struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valute  []struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Nominal  string `xml:"Nominal"`
		Name     string `xml:"Name"`
		Value    string `xml:"Value"`
	} `xml:"Valute"`
}

const (
	acceptHeader         = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"
	acceptEncodingHeader = "gzip, deflate, br"
	acceptLanguageHeader = "en-US,en;q=0.9,ru;q=0.8"
	cacheControlHeader   = "max-age=0"
)
