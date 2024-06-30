package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sp0x/cryptodashboard/config"
	"github.com/sp0x/cryptodashboard/datasources/coingecko"
	"github.com/sp0x/cryptodashboard/events"
	"github.com/sp0x/cryptodashboard/models"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := config.Load("config.yaml"); err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	signalChannel := make(chan struct{})
	updateChannel := make(chan models.CurrencyUpdate)
	sourceChannel := coingecko.StartPolling(ctx,
		time.Duration(config.Get.UpdateIntervalInSec)*time.Second,
		config.Get.CoinGeckoAPIKey,
		config.Get.TargetCurrency,
		config.Get.CryptoCurrencyIds)

	eventHandler := events.EventHandler{
		DataChannel: updateChannel,
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(updateChannel)
				return
			case updateEvent := <-sourceChannel:
				eventHandler.LastUpdate = updateEvent
				select {
				case signalChannel <- struct{}{}:
				default:
				}

				updateChannel <- updateEvent
			}
		}
	}()

	staticFileServer := http.FileServer(http.Dir("./static"))

	http.Handle("/", staticFileServer)
	http.Handle("/feed", http.HandlerFunc(eventHandler.HandlerFunc))

	<-signalChannel
	bindPoint := fmt.Sprintf("%s:%d", config.Get.Host, config.Get.Port)
	fmt.Printf("Starting server on http://%s\n", bindPoint)
	err := http.ListenAndServe(bindPoint, nil)

	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
