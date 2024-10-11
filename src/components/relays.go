package components

import (
	"log"
	"net/http"

	"gfnwc/src/utils"
)

func RelaysComponent(w http.ResponseWriter, r *http.Request) {
	log.Println("RelaysComponent called")

	// Prepare the data to be passed to the template
	data := utils.PageData{
	}

	// Render the template, specifying the view path and the layout usage
	utils.RenderTemplate(w, data, "components/relays.html", false)
}
