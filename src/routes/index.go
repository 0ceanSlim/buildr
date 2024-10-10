package routes

import (
	"gfnwc/src/handlers"
	"gfnwc/src/utils"
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

	data := utils.PageData{
		Title:       "Dashboard",
		DisplayName: displayName,
		Picture:     picture,
		PublicKey:   publicKey,
		About:       about,
	}

	utils.RenderTemplate(w, data, "index.html", false)
}
