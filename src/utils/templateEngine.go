package utils

import (
	"buildr/src/types"
	"html/template"
	"net/http"
	"path/filepath"
)

type PageData struct {
	Title         string
	Theme         string
	PublicKey     string
	DisplayName   string
	Picture       string
	About         string
	Relays        RelayList
	Message       string   // For displaying the message content
	SuccessRelays []string // List of relays where the message was successfully sent
	FailedRelays  []string // List of relays where the message failed
	Notes       []types.NostrEvent // Add the notes field to store kind 1 events
}

// Define the base directories for views and templates
const (
	viewsDir     = "web/views/"
	templatesDir = "web/views/templates/"
)

// Define the common layout templates filenames
var templateFiles = []string{
	"layout.html",
	"header.html",
	"footer.html",
}

// Initialize the common templates with full paths
var layout = PrependDir(templatesDir, templateFiles)

var loginLayout = PrependDir(templatesDir, []string{"login-layout.html", "footer.html"})

func RenderTemplate(w http.ResponseWriter, data PageData, view string, useLoginLayout bool) {
    viewTemplate := filepath.Join(viewsDir, view)
    componentPattern := filepath.Join(viewsDir, "components", "*.html")
    componentTemplates, err := filepath.Glob(componentPattern)
    if err != nil {
        http.Error(w, "Error loading component templates: "+err.Error(), http.StatusInternalServerError)
        return
    }

    var templates []string
    if useLoginLayout {
        templates = append(loginLayout, viewTemplate)
    } else {
        templates = append(layout, viewTemplate)
    }
    templates = append(templates, componentTemplates...)

    tmpl, err := template.New("").Funcs(template.FuncMap{
        "formatTimestamp":   formatTimestamp,
        "renderNoteContent": renderNoteContent, // Register the content rendering function
    }).ParseFiles(templates...)
    if err != nil {
        http.Error(w, "Error parsing templates: "+err.Error(), http.StatusInternalServerError)
        return
    }

    layoutName := "layout"
    if useLoginLayout {
        layoutName = "login-layout"
    }
    err = tmpl.ExecuteTemplate(w, layoutName, data)
    if err != nil {
        http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
    }
}
