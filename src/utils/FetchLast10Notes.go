package utils

import (
	"encoding/json"
	"log"
	"sort"
	"sync"

	"gfnwc/src/types"

	"github.com/gorilla/websocket"
)

const MaxNotes = 10 // Define max number of notes to fetch

func FetchLast10Kind1Notes(publicKey string, relays []string) ([]types.NostrEvent, error) {
	var notes []types.NostrEvent
	uniqueNotes := make(map[string]types.NostrEvent) // Map to store unique notes by their IDs
	var notesMu sync.Mutex                           // Mutex to protect access to uniqueNotes
	var wg sync.WaitGroup                            // WaitGroup to manage concurrent requests

	maxNotes := int64(MaxNotes)
	results := make(chan types.NostrEvent)

	// Create a function to process each relay connection
	for _, url := range relays {
		wg.Add(1)
		go func(relayURL string) {
			defer wg.Done()
			log.Printf("Connecting to WebSocket: %s\n", relayURL)
			conn, _, err := websocket.DefaultDialer.Dial(relayURL, nil)
			if err != nil {
				log.Printf("Failed to connect to WebSocket: %v\n", err)
				return
			}
			defer conn.Close()

			filter := types.SubscriptionFilter{
				Authors: []string{publicKey},
				Kinds:   []int{1},
				Limit:   &maxNotes,
			}

			subRequest := []interface{}{
				"REQ",
				"sub2",
				filter,
			}

			requestJSON, err := json.Marshal(subRequest)
			if err != nil {
				log.Printf("Failed to marshal subscription request: %v\n", err)
				return
			}

			if err := conn.WriteMessage(websocket.TextMessage, requestJSON); err != nil {
				log.Printf("Failed to send subscription request: %v\n", err)
				return
			}

			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					log.Printf("Error reading WebSocket message: %v\n", err)
					return
				}

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

					// Add the note to the map if it's unique
					notesMu.Lock()
					if _, exists := uniqueNotes[event.ID]; !exists {
						uniqueNotes[event.ID] = event
						results <- event
					}
					notesMu.Unlock()
				} else if response[0] == "EOSE" {
					log.Println("End of subscription signal received")
					break
				}
			}
		}(url)
	}

	// Collect up to MaxNotes unique notes
	go func() {
		wg.Wait()
		close(results)
	}()

	for note := range results {
		notes = append(notes, note)
		if len(notes) >= MaxNotes {
			break
		}
	}

	// Sort notes by timestamp in descending order
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].CreatedAt > notes[j].CreatedAt
	})

	return notes, nil
}
