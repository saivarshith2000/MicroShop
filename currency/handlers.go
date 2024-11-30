package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ConvertCurrencyRequest struct {
	Value float64 `form:"value" binding:"required"`
	From  string  `form:"from" binding:"required"`
	To    string  `form:"to" binding:"required"`
}

type ConversionResponse struct {
	Value float64 `json:"value"`
	From  string  `json:"from"`
	To    string  `json:"to"`
}

/*
Converts money from one currency to another
Params: /convert?value=50&original=USD&target=AUD (or)
Params: /convert?value=50&target=AUD
*/
func getConvertedCurrency(c *gin.Context) {
	var req ConvertCurrencyRequest
	if err := c.BindQuery(&req); err != nil {
		c.JSON(400, "Invalid request parameters")
		return
	}
	original_rate, ok := ExchangeRateCache.Rates[req.From]
	if !ok {
		c.JSON(400, "Invalid original currency")
		return
	}
	target_rate, ok := ExchangeRateCache.Rates[req.To]
	if !ok {
		c.JSON(400, "Invalid target currency")
		return
	}
	if original_rate == target_rate {
		c.JSON(400, "From and target currencies must be different")
		return
	}
	convertedValue := (target_rate / original_rate) * req.Value
	log.Info().
		Str("from", req.From).
		Str("to", req.To).
		Str("originalValue", fmt.Sprintf("%v", req.Value)).
		Str("convertedValue", fmt.Sprintf("%v", convertedValue)).
		Msg("Conversion result")
	c.JSON(200, ConversionResponse{
		Value: convertedValue,
		From:  req.From,
		To:    req.To,
	})
}
