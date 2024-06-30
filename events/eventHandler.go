package events

import (
	"bytes"
	"cryptodashboard/models"
	"encoding/json"
	"fmt"
	"net/http"
)

type EventHandler struct {
	LastUpdate  models.CurrencyUpdate
	DataChannel chan models.CurrencyUpdate
}

func (e *EventHandler) HandlerFunc(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	sendData := func(updateEvent models.CurrencyUpdate) {
		var buffer bytes.Buffer
		err := json.NewEncoder(&buffer).Encode(updateEvent.Currencies)
		if err != nil {
			fmt.Printf("Error encoding JSON: %v\n", err)
			return
		}
		fmt.Fprintf(w, "event: priceUpdate\ndata: %s\n\n", buffer.String())
		w.(http.Flusher).Flush()
	}

	sendData(e.LastUpdate)

	for updateEvent := range e.DataChannel {
		sendData(updateEvent)
	}
}
