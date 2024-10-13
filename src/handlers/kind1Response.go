package handlers

import (
	"encoding/json"
	"gfnwc/src/utils"
	"net/http"
)

// Kind1ResponseHandler displays the status of sending the message.
func Kind1ResponseHandler(w http.ResponseWriter, r *http.Request) {
	var relayStatuses map[string]string
	err := json.NewDecoder(r.Body).Decode(&relayStatuses)
	if err != nil {
		http.Error(w, "Failed to decode response data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Create slices for successful and failed relays.
	var successRelays, failedRelays []string
	for relay, status := range relayStatuses {
		if status == "Success" {
			successRelays = append(successRelays, relay)
		} else {
			failedRelays = append(failedRelays, relay)
		}
	}

	// Prepare additional data for rendering.
	data := utils.PageData{
		Title:         "Message Status",
		PublicKey:     "UserPublicKey",          // Replace with actual public key if available.
		Message:       "Your message was sent!", // Customize this as needed.
		SuccessRelays: successRelays,
		FailedRelays:  failedRelays,
		Relays: utils.RelayList{
			Read:  successRelays, // If needed, you can use success relays as read relays.
			Write: successRelays, // Similar mapping if needed.
		},
	}

	// Render the response template with the additional data.
	utils.RenderTemplate(w, data, "components/kind1-response.html", false)
}
