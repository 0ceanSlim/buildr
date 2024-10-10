// src/components/profile.go
package components

import (
	"net/http"

	"gfnwc/src/utils"
)

// ProfileHandler serves the profile view
func DropdownComponent(w http.ResponseWriter, r *http.Request) {

    // Prepare the page data
    data := utils.PageData{
        Title:       "dropdown menu",

    }

    // Render the profile view using the common layout
    utils.RenderTemplate(w, data, "components/dropdown.html", false)
}
