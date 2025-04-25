package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func cleanProfanity(text string) string {
    profaneWords := map[string]bool{
        "kerfuffle": true,
        "sharbert":  true,
        "fornax":    true,
    }
	// Split the text into tokens (words)
	tokens := strings.Split(text, " ")

    for i, token := range tokens {
        // Check if the token is a profane word (case-insensitive)
        if profaneWords[strings.ToLower(token)] {
            tokens[i] = "****" // Replace the word
        }
    }

    // Reconstruct the text with spaces and punctuation preserved
    return strings.Join(tokens, " ")
}


func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong"})
		return
	}

	if len(params.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Chirp is too long"})
		return
	}

	cleanedBody := cleanProfanity(params.Body)


	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"cleaned_body": cleanedBody})
}
