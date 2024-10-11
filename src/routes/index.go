package routes

import (
	"gfnwc/src/handlers"
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

	// Fetch the relay list from the session
	relays, ok := session.Values["relays"].(utils.RelayList)
	if !ok {
		log.Println("No relay list found in session for Index view")
		// Optionally, you can initialize it to avoid nil issues in templates
		relays = utils.RelayList{}
	}

	data := utils.PageData{
		Title:       "Dashboard",
		DisplayName: displayName,
		Picture:     picture,
		PublicKey:   publicKey,
		About:       about,
		Relays:      relays,
	}

	utils.RenderTemplate(w, data, "index.html", false)
}
