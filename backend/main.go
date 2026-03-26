package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Exemple de fonction Handler qui renvoie du JSON
func handleGetNews(w http.ResponseWriter, r *http.Request) {
	// Obligatoire pour autoriser React (qui tourne sur un autre port) à faire des requêtes
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json") // Les réponses doivent être en JSON

	// Données fictives pour l'exemple
	response := map[string]string{"message": "Voici la liste des articles !"}
	
	// Encodage en JSON et envoi
	json.NewEncoder(w).Encode(response)
}

func main() {
	// native router creation
	mux := http.NewServeMux()
	
	// routes definition
	// mux.HandleFunc("GET /api/news", handleGetNews)
	// mux.HandleFunc("POST /api/comments", handlePostComment)
	// mux.HandleFunc("POST /api/vote", handleVote)

	// launching of the backround task to fetch regularly the news from the externa news API
	go func() {
		ticker := time.NewTicker(24 * time.Hour) // every day
		for range ticker.C {
			fmt.Println("Updating the database with the news API...")
			// TODO: call the service which fetch the news
		}
	}()

	// starting the server
	port := ":8080"
	fmt.Println("Server started on port ", port)
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal("Error at the start of the server :", err)
	}
}