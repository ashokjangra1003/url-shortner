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
	OrignalURL   string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

/*
	fdfdf233 --> {

	ID: fdfdf233 ,
	OrignalURL : "www.abc.com",
	ShortURL : "fdfdf233"
	CreationDate:  time.Now()
	}

*/

var urlDB = make(map[string]URL)

func genrateShortURL(OrignalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OrignalURL))
	fmt.Println("hasher", hasher)
	data := hasher.Sum(nil)
	fmt.Println("hasher data", data)
	hash := hex.EncodeToString(data)
	fmt.Println("hfinal string", hash[:9])

	return hash[:9]
}

func createUrl(orignalURL string) string {

	shortURL := genrateShortURL(orignalURL)
	id := shortURL
	urlDB[id] = URL{
		ID:           id,
		OrignalURL:   orignalURL,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}
	return shortURL
}

func getURL(id string) (URL, error) {
	url, ok := urlDB[id]

	if !ok {
		return URL{}, errors.New("URl not found")
	}
	return url, nil

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Get method ")

}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "invaild reuest body", http.StatusBadRequest)
		return
	}
	shortURL_ := createUrl(data.URL)
	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL_}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil {
		http.Error(w, "invaild request", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.OrignalURL, http.StatusFound)
}

func main() {
	// fmt.Println("Starting Url shortner... .")
	// OrignalURL := "https://github.com/ashokjangra1003"
	// genrateShortURL(OrignalURL)
	http.HandleFunc("/", handler)
	http.HandleFunc("/shorten", ShortURLHandler)
	http.HandleFunc("/redirect/", redirectURLHandler)

	//start htpp host
	fmt.Println("Server start running port number 30000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("error on starting server ", err)
	}

}
