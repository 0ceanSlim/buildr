package routes

import (
	"gfnwc/src/handlers"
	"gfnwc/src/types"
	"gfnwc/src/utils"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
    session, _ := handlers.User.Get(r, "session-name")

    publicKey, ok := session.Values["publicKey"].(string)
    if !ok || publicKey == "" {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    displayName, _ := session.Values["displayName"].(string)
    picture, _ := session.Values["picture"].(string)
    about, _ := session.Values["about"].(string)

    relays, ok := session.Values["relays"].(utils.RelayList)
    if !ok {
        log.Println("No relay list found in session for Index view")
        relays = utils.RelayList{}
    }

    // Convert RelayList to []string if necessary
    relayURLs := relays.ToStringSlice()

    // Fetch the last 10 kind 1 notes
    notes, err := utils.FetchLast10Kind1Notes(publicKey, relayURLs)
    if err != nil {
        log.Printf("Failed to fetch last 10 kind 1 notes: %v\n", err)
        notes = []types.NostrEvent{}
    }

    // Store notes in session if needed
    session.Values["notes"] = notes
    if err := session.Save(r, w); err != nil {
        log.Printf("Failed to save session: %v\n", err)
    }

    data := utils.PageData{
        Title:       "Dashboard",
        DisplayName: displayName,
        Picture:     picture,
        PublicKey:   publicKey,
        About:       about,
        Relays:      relays,
        Notes:       notes, // Pass the notes to the page data
    }

    utils.RenderTemplate(w, data, "index.html", false)
}
