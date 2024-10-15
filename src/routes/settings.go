package routes

import (
	"log"
	"net/http"

	"buildr/src/handlers"
	"buildr/src/utils"
)

func Settings(w http.ResponseWriter, r *http.Request) {
	//log.Println("Settings handler called")

	session, err := handlers.User.Get(r, "session-name")
	if err != nil {
		log.Printf("Error getting session: %v\n", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	publicKey, ok := session.Values["publicKey"].(string)
	if !ok || publicKey == "" {
		log.Println("No publicKey found in session")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Fetch the relay list from the session
	relays, ok := session.Values["relays"].(utils.RelayList)
	if !ok {
		log.Println("No relay list found in session for Settings view")
		relays = utils.RelayList{} // Initialize it to avoid nil issues in templates if required
	}

	// Prepare the data to be passed to the template
	data := utils.PageData{
		Title:     "Settings",
		PublicKey: publicKey,
		Relays:    relays,
	}

	// Render the template
	utils.RenderTemplate(w, data, "settings.html", false)
}
