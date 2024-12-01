package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/saivarshith2000/MicroShop/currency/middleware"
)

func getEnv() (map[string]string, error) {
	var env = make(map[string]string)
	keys := []string{"FX_API_URL", "FX_API_TOKEN"}
	for _, k := range keys {
		value := os.Getenv(k)
		if value == "" {
			return nil, fmt.Errorf("environment variable %s not found", k)
		}
		env[k] = value
	}
	return env, nil
}

var ExchangeRateCache *ExchangeRateCacheType

func promHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel) // TODO: Get this from env

	env, err := getEnv()
	if err != nil {
		log.Fatal().Str("msg", fmt.Sprintf("Error loading environment variables - %s", err.Error()))
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.DefaultStructuredLogger())
	r.Use(gin.Recovery())
	r.Use(middleware.HttpRequestCounterMiddleware())

	prometheus.Register(middleware.HttpRequestsTotal)
	prometheus.Register(middleware.RequestLatency)

	initExchangeRateCache()
	err = updateExchangeRatesCache(env["FX_API_URL"], env["FX_API_TOKEN"])
	if err != nil {
		log.Fatal().Str("msg", fmt.Sprintf("Error fetching exchange rates - %s", err.Error()))
	}

	r.GET("/convert", getConvertedCurrency)
	r.GET("/metrics", promHandler())
	r.Run()
}
