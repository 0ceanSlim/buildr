package handlers

import (
	"gfnwc/src/types"
	"gfnwc/src/utils"
	"log"
	"net/http"
)

func FetchNotes(w http.ResponseWriter, r *http.Request) {
    session, _ := User.Get(r, "session-name")

    publicKey, ok := session.Values["publicKey"].(string)
    if !ok || publicKey == "" {
        http.Error(w, "Missing public key", http.StatusUnauthorized)
        return
    }

    relays, ok := session.Values["relays"].(utils.RelayList)
    if !ok {
        log.Println("No relay list found in session")
        relays = utils.RelayList{}
    }

    relayURLs := relays.ToStringSlice()
    notes, err := utils.FetchLast10Kind1Notes(publicKey, relayURLs)
    if err != nil {
        log.Printf("Failed to fetch last 10 kind 1 notes: %v\n", err)
        notes = []types.NostrEvent{}
    }

    data := utils.PageData{
        Notes: notes, // Only the notes are needed here
    }

    utils.RenderTemplate(w, data, "components/notes.html", false) // Only render the notes template
}
