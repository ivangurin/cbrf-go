package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ivangurin/cbrf-go"
)

func main() {

	now := time.Now()

	// Получение курса валюты на дату
	exchangeRate, err := cbrf.GetExchangeRate(cbrf.CurrencyUSD, now)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Текущий курс USD: %.2f рублей\n", exchangeRate)

	// Конвертация 100 USD в EUR на дату
	valueUSD := float64(100)

	valueEUR, err := cbrf.Convert(cbrf.CurrencyUSD, cbrf.CurrencyEUR, valueUSD, now)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Стоимость %.2f USD равна %.2f EUR\n", valueUSD, valueEUR)

}
