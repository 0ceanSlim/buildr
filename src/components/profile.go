// src/components/profile.go
package components

import (
	"net/http"

	"gfnwc/src/handlers"
	"gfnwc/src/utils"
)

// ProfileHandler serves the profile view
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
    session, _ := handlers.User.Get(r, "session-name")

    // Retrieve the necessary user data from the session
    publicKey, _ := session.Values["publicKey"].(string)
    displayName, _ := session.Values["displayName"].(string)
    picture, _ := session.Values["picture"].(string)
    about, _ := session.Values["about"].(string)
    relays, _ := session.Values["relays"].(utils.RelayList)

    // Prepare the page data
    data := utils.PageData{
        Title:       "User Profile",
        PublicKey:   publicKey,
        DisplayName: displayName,
        Picture:     picture,
        About:       about,
        Relays:      relays,
    }

    // Render the profile view using the common layout
    utils.RenderTemplate(w, data, "components/profile.html", false)
}
