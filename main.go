package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type TokenStore struct {
	issued map[string]time.Time
}

var store *TokenStore

func issueToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := "tok_" + time.Now().Format("20060102150405")
	expiresAt := time.Now().Add(24 * time.Hour).UTC()
	store.issued[token] = expiresAt

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"token":      token,
		"expires_at": expiresAt.Format(time.RFC3339),
	})
}

func revokeToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "revoked"})
}

func main() {
	http.HandleFunc("/tokens/issue", issueToken)
	http.HandleFunc("/tokens/revoke", revokeToken)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	log.Println("token-service listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
