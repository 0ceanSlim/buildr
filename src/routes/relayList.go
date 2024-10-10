package routes

import (
	"log"
	"net/http"

	"gfnwc/src/handlers"
	"gfnwc/src/utils"
)

func RelayList(w http.ResponseWriter, r *http.Request) {
	log.Println("RelayListHandler called")

	session, _ := handlers.User.Get(r, "session-name")

	publicKey, ok := session.Values["publicKey"].(string)
	if !ok || publicKey == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Fetch the relay list from the session
	relays, ok := session.Values["relays"].(utils.RelayList)
	if !ok {
		log.Println("No relay list found in session")
		http.Error(w, "Relay list not found", http.StatusInternalServerError)
		return
	}

	// Prepare the data to be passed to the template
	data := utils.PageData{
		Title:     "User Relays",
		PublicKey: publicKey,
		Relays:    relays,
	}

	// Render the template
	utils.RenderTemplate(w, data, "relayList.html", false)
}
