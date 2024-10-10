package handlers

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"

	"gfnwc/src/utils"

	"github.com/gorilla/sessions"
)

var User = sessions.NewCookieStore([]byte("your-secret-key"))

func init() {
	// Register the RelayList type with gob
	gob.Register(utils.RelayList{})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("LoginHandler called")

	if err := r.ParseForm(); err != nil {
		log.Printf("Failed to parse form: %v\n", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	publicKey := r.FormValue("publicKey")
	if publicKey == "" {
		log.Println("Missing publicKey in form data")
		http.Error(w, "Missing publicKey", http.StatusBadRequest)
		return
	}

	log.Printf("Received publicKey: %s\n", publicKey)

	// Log the public key to a file
	logPublicKey(publicKey)

	// Fetch user relay list from an initial relay
	initialRelays := []string{
		"wss://purplepag.es", "wss://relay.damus.io", "wss://nos.lol", "wss://relay.primal.net", "wss://relay.nostr.band", "wss://offchain.pub", // Add any initial relay URLs here
	}
	userRelays, err := utils.FetchUserRelays(publicKey, initialRelays)
	if err != nil {
		log.Printf("Failed to fetch user relays: %v\n", err)
		http.Error(w, "Failed to fetch user relays", http.StatusInternalServerError)
		return
	}
	log.Printf("Fetched user relays: %+v\n", userRelays)

	// Combine all relays (read, write, both) into a single slice
	allRelays := append(userRelays.Read, userRelays.Write...)
	allRelays = append(allRelays, userRelays.Both...)

	// Fetch user metadata from the combined relay list
	userContent, err := utils.FetchUserMetadata(publicKey, allRelays)
	if err != nil {
		log.Printf("Failed to fetch user metadata: %v\n", err)
		http.Error(w, "Failed to fetch user metadata", http.StatusInternalServerError)
		return
	}
	log.Printf("Fetched user metadata: %+v\n", userContent)

	// Store the public key, user data, and relays in the session
	session, _ := User.Get(r, "session-name")
	session.Values["publicKey"] = publicKey
	session.Values["displayName"] = userContent.DisplayName
	session.Values["picture"] = userContent.Picture
	session.Values["about"] = userContent.About
	session.Values["relays"] = userRelays // Store the relay list categorized by read, write, and both
	if err := session.Save(r, w); err != nil {
		log.Printf("Failed to save session: %v\n", err)
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	log.Println("Session saved successfully")

	// Redirect to the root ("/")
	http.Redirect(w, r, "/", http.StatusSeeOther)
	log.Println("Redirecting to /")
}

// logPublicKey logs the public key to a text file in the logs directory, avoiding duplicates
func logPublicKey(publicKey string) {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		log.Printf("Error creating logs directory: %v\n", err)
		return
	}

	logFilePath := "logs/login_log.txt"

	// Read existing keys to check for duplicates
	existingKeys, err := readExistingKeys(logFilePath)
	if err != nil {
		log.Printf("Error reading existing keys: %v\n", err)
		return
	}

	// Check if the public key already exists
	if _, exists := existingKeys[publicKey]; exists {
		log.Printf("Public key already logged: %s\n", publicKey)
		return
	}

	// Open the log file in the logs directory
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("%s\n", publicKey)); err != nil {
		log.Printf("Error writing to log file: %v\n", err)
	}
}

// readExistingKeys reads the existing public keys from the log file
func readExistingKeys(filePath string) (map[string]struct{}, error) {
	keys := make(map[string]struct{})

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist, return empty map
			return keys, nil
		}
		return nil, err
	}
	defer file.Close()

	var key string
	for {
		_, err := fmt.Fscanln(file, &key)
		if err != nil {
			break
		}
		keys[key] = struct{}{}
	}

	return keys, nil
}