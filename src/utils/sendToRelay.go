package utils

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/nbd-wtf/go-nostr"
)

// SendToRelay sends the signed Nostr event to the specified WebSocket relay
func SendToRelay(relayURL string, event nostr.Event) error {
	// Open a WebSocket connection to the relay
	conn, _, err := websocket.DefaultDialer.Dial(relayURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to relay: %v", err)
	}
	defer conn.Close()


	// Prepare the message as per Nostr protocol: ["EVENT", <event JSON>]
	message := []interface{}{"EVENT", event}
	err = conn.WriteJSON(message)
	if err != nil {
		return fmt.Errorf("failed to send event to relay: %v", err)
	}

	// Optionally, wait for an acknowledgment from the relay (if required)
	_, reply, err := conn.ReadMessage()
	if err != nil {
		log.Printf("Failed to read response from relay: %v", err)
	} else {
		log.Printf("Response from relay: %s", reply)
	}

	return nil
}
