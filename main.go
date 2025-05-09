package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
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
var corsOrigin = os.Getenv("ALLOW_CORS_ORIGIN")

func main() {
	_, err := fetchFileNames()
	if err != nil {
		log.Fatalf("Failed to fetch file names: %v", err)
	}

	http.HandleFunc("OPTIONS /", cors)
	http.HandleFunc("GET /", get)
	http.HandleFunc("GET /stats", stats)
	http.HandleFunc("POST /refresh", refresh)

	port := "8080"
	log.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func cors(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", corsOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintf(w, "")
}

func get(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
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

	w.Header().Set("Access-Control-Allow-Origin", corsOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")
	w.Header().Set("Access-Control-Allow-Headers", "*")
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

	w.Header().Set("Access-Control-Allow-Origin", corsOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")
	w.Header().Set("Access-Control-Allow-Headers", "*")
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
