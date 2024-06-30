package coingecko

import (
	"context"
	"cryptodashboard/models"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

func StartPolling(ctx context.Context,
	updateInterval time.Duration,
	apiKey string,
	targetCurrency string,
	cryptoCurrencyIds []string) chan models.CurrencyUpdate {
	client := NewCoingeckoClient(nil, apiKey, CoingeckoDefaultBaseUrl)
	dataChannel := make(chan models.CurrencyUpdate)
	simplePriceUrl := GetSimplePriceEndpoint(targetCurrency, cryptoCurrencyIds)

	go func() {
		ticker := time.NewTicker(updateInterval)
		defer ticker.Stop()
		defer close(dataChannel)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				bytes, err := client.Get(simplePriceUrl)
				if err != nil {
					fmt.Printf("Error fetching data: %v\n", err)
					continue
				}
				var prices Price
				err = json.Unmarshal(bytes, &prices)
				if err != nil {
					fmt.Printf("Error deserializing data: %v\n", err)
				}

				updateEvent := models.CurrencyUpdate{
					Currencies: make([]models.CurrencyValue, len(prices)),
				}

				var currencyNames []string = make([]string, len(prices))
				index := 0
				for k := range prices {
					currencyNames[index] = k
					index++
				}
				sort.Strings(currencyNames)
				index = 0
				for i, currency := range currencyNames {
					updateEvent.Currencies[i] = models.CurrencyValue{
						Name:   currency,
						Value:  prices[currency][targetCurrency],
						Symbol: resolveSymbol(currency),
					}
				}

				dataChannel <- updateEvent
			}
		}
	}()

	return dataChannel
}

func resolveSymbol(currency string) string {
	switch currency {
	case "bitcoin":
		return "btc"
	case "ethereum":
		return "eth"
	case "cardano":
		return "ada"
	default:
		return strings.ToLower(currency)
	}
}
