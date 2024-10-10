package main

import (
	"embed"
	"gfnwc/src/components"
	"gfnwc/src/handlers"
	"gfnwc/src/routes"
	"gfnwc/src/utils"

	"fmt"
	"net/http"
)

//go:embed web/*
var staticFiles embed.FS

func main() {
	// Load Configurations
	cfg, err := utils.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	mux := http.NewServeMux()
	// Login / Logout
	mux.HandleFunc("/login", routes.Login) // Login route
	mux.HandleFunc("/do-login", handlers.LoginHandler)
	mux.HandleFunc("/logout", handlers.LogoutHandler) // Logout process

	// Initialize Routes
	mux.HandleFunc("/", routes.Index)
	mux.HandleFunc("/relay-list", routes.RelayList)

	// Render component htmls
	mux.HandleFunc("/profile", components.ProfileHandler)

	// Function Handlers
	
	// Serve Static Files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/img/favicon.ico")
	})
	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.FS(staticFiles))))

	mux.HandleFunc("/wip-message", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<button class="px-4 py-2 mt-4 text-xs font-semibold text-white bg-red-500 rounded-md hover:bg-red-700">I'm Working on it ⚠️</button>`)
	})

	fmt.Printf("Server is running on http://localhost:%d\n", cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), mux)
}
