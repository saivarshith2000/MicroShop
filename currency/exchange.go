package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type ExchangeRateCacheType struct {
	Rates      map[string]float64
	ValiedUpto time.Time
}

type ExchangeRateResponse struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
	Error string             `json:"error,omitempty"`
}

func initExchangeRateCache() {
	ExchangeRateCache = &ExchangeRateCacheType{}
	ExchangeRateCache.Rates = nil
	ExchangeRateCache.ValiedUpto = time.Now()
}

func updateExchangeRatesCache(apiURL, apiKey string) error {
	if ExchangeRateCache.ValiedUpto.Compare(time.Now()) <= 0 {
		log.Info().Msg("Updating exchange rate cache")
		response, err := FetchExchangeRates(apiURL, apiKey)
		if err != nil {
			return fmt.Errorf("failed to fetch exchange rates")
		}
		ExchangeRateCache.Rates = response.Rates
		ExchangeRateCache.ValiedUpto = time.Now().Add(time.Duration(time.Hour * 1))
		log.Info().Msg(fmt.Sprintf("successfully fetched exchange rates. Cache valid upto %s\n", ExchangeRateCache.ValiedUpto))
		return nil
	}
	return nil
}

/*
Fetch latest exchange rates from FXRatesAPI and stores them in a global var in memory :)
*/
func FetchExchangeRates(apiURL, apiKey string) (*ExchangeRateResponse, error) {
	url := fmt.Sprintf("%s?api=%s", apiURL, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange rates: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-200 response: %s", string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var response ExchangeRateResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	if response.Error != "" {
		return nil, fmt.Errorf("API error: %s", response.Error)
	}

	return &response, nil
}
