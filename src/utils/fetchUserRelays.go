package utils

import (
	"encoding/json"
	"log"
	"time"

	"buildr/src/types"

	"github.com/gorilla/websocket"
)

type RelayList struct {
	Read  []string
	Write []string
	Both  []string
}

// ToStringSlice combines Read, Write, and Both into a single []string.
func (r RelayList) ToStringSlice() []string {
	var urls []string
	urls = append(urls, r.Read...)
	urls = append(urls, r.Write...)
	urls = append(urls, r.Both...)
	return urls
}

func FetchUserRelays(publicKey string, relays []string) (*RelayList, error) {
	for _, url := range relays {
		log.Printf("Connecting to WebSocket: %s\n", url)
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Printf("Failed to connect to WebSocket: %v\n", err)
			continue
		}
		defer conn.Close()

		filter := types.SubscriptionFilter{
			Authors: []string{publicKey},
			Kinds:   []int{10002}, // Kind 10002 corresponds to relay list (NIP-65)
		}

		subRequest := []interface{}{
			"REQ",
			"sub1",
			filter,
		}

		requestJSON, err := json.Marshal(subRequest)
		if err != nil {
			log.Printf("Failed to marshal subscription request: %v\n", err)
			return nil, err
		}

		log.Printf("Sending subscription request: %s\n", requestJSON)

		if err := conn.WriteMessage(websocket.TextMessage, requestJSON); err != nil {
			log.Printf("Failed to send subscription request: %v\n", err)
			return nil, err
		}

		// WebSocket message or timeout handling
		msgChan := make(chan []byte)
		errChan := make(chan error)

		go func() {
			_, message, err := conn.ReadMessage()
			if err != nil {
				errChan <- err
			} else {
				msgChan <- message
			}
		}()

		select {
		case message := <-msgChan:
			log.Printf("Received WebSocket message: %s\n", message)
			var response []interface{}
			if err := json.Unmarshal(message, &response); err != nil {
				log.Printf("Failed to unmarshal response: %v\n", err)
				continue
			}

			if response[0] == "EVENT" {
				eventData, err := json.Marshal(response[2])
				if err != nil {
					log.Printf("Failed to marshal event data: %v\n", err)
					continue
				}

				var event types.NostrEvent
				if err := json.Unmarshal(eventData, &event); err != nil {
					log.Printf("Failed to parse event data: %v\n", err)
					continue
				}

				log.Printf("Received Nostr event: %+v\n", event)

				relayList := &RelayList{}

				for _, tag := range event.Tags {
					if len(tag) > 1 && tag[0] == "r" {
						relayURL := tag[1]
						if len(tag) == 3 {
							switch tag[2] {
							case "read":
								relayList.Read = append(relayList.Read, relayURL)
							case "write":
								relayList.Write = append(relayList.Write, relayURL)
							}
						} else {
							relayList.Both = append(relayList.Both, relayURL)
						}
					}
				}
				return relayList, nil
			}
		case err := <-errChan:
			log.Printf("Error reading WebSocket message: %v\n", err)
			continue
		case <-time.After(WebSocketTimeout):
			log.Printf("WebSocket response timeout from %s\n", url)
			continue
		}
	}
	return nil, nil
}
