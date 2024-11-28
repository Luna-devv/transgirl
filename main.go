package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

var fileMap = make(map[int]string)
var mapMutex = &sync.RWMutex{}

var awsRegion = os.Getenv("AWS_REGION")
var awsBucket = os.Getenv("AWS_BUCKET")
var accessKey = os.Getenv("AWS_ACCESS_KEY_ID")
var secretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
var awsEndpoint = os.Getenv("AWS_ENDPOINT")
var awsPublicUrl = os.Getenv("AWS_PUBLIC_URL")
var prefix = os.Getenv("FILE_PREFIX")

func main() {
	rand.Seed(time.Now().UnixNano())

	go func() {
		_, err := fetchFileNames()
		if err != nil {
			log.Fatalf("Failed to fetch file names: %v", err)
		}
	}()

	http.HandleFunc("GET /", get)
	http.HandleFunc("GET /stats", stats)
	http.HandleFunc("POST /refresh", refresh)

	port := "8080"
	log.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func get(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not Found")
		return
	}

	mapMutex.RLock()
	defer mapMutex.RUnlock()

	if len(fileMap) == 0 {
		http.Error(w, "No files available", http.StatusInternalServerError)
		return
	}

	randomIndex := rand.Intn(len(fileMap))
	fileName := fileMap[randomIndex]

	response := map[string]string{
		"url": fmt.Sprintf("%s/%s", awsPublicUrl, fileName),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func stats(w http.ResponseWriter, r *http.Request) {
	mapMutex.RLock()
	defer mapMutex.RUnlock()

	response := map[string]int{
		"file_count": len(fileMap),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func refresh(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth != "Bearer "+secretKey {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	mapMutex.Lock()
	defer mapMutex.Unlock()

	count, err := fetchFileNames()
	if err != nil {
		http.Error(w, "Failed to fetch file names", http.StatusInternalServerError)
		return
	}

	response := map[string]int{
		"file_count": count,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
