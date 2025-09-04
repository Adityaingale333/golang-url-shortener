package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortenedURL string    `json:"shortened_url"`
	CreatedAt    time.Time `json:"created_at"`
}

/*
	d9736711 ->{
					ID : "d9736711",
					OriginalURL : "https://github.com/Aditya-1982",
					ShortenURL : "d9736711",
					CreatedAt : 4 sept 2025
	}
*/
var urlDB = make(map[string]URL)

func generateShortURL(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL)) // it converts originalURL to byte slice
	//fmt.Println("hasher: ", hasher)
	data := hasher.Sum(nil)
	//fmt.Println("hasher data: ", data)
	hash := hex.EncodeToString(data)
	//fmt.Println("hash: ", hash)
	return hash[:8] // return first 8 characters of the hash
}

func createURL(originalURL string) string {
	shortURL := generateShortURL(originalURL)
	id := shortURL // In this case, we use the short URL as the ID
	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  originalURL,
		ShortenedURL: shortURL,
		CreatedAt:    time.Now(),
	}
	return shortURL
}

func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func shortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data) // Decode the JSON request body into the data struct
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	shortURL := createURL(data.URL)

	response := struct {
		ShortenURL string `json:"shorten_url"`
	}{ShortenURL: shortURL}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):] // Extract the ID from the URL path
	url, err := getURL(id)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main() {
	fmt.Println("Starting URL shortener service...")

	// Register the handler function to handle all the requests to the root url ("/")
	http.HandleFunc("/", handler)

	http.HandleFunc("/shorten", shortURLHandler)

	http.HandleFunc("/redirect/", redirectHandler)

	// Start the HTTP server on port 3000
	fmt.Println("Server is running on port 3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
