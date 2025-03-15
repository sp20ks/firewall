package auth

import (
	"encoding/json"
	"log"
	"net/http"
)

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleGetJwtToken(a *Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UserRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Username == "" || req.Password == "" {
			log.Println("Both 'username' and 'password' must be provided")
			http.Error(w, "Both 'username' and 'password' must be provided", http.StatusBadRequest)
			return
		}

		tokenString, err := a.CreateToken(req.Username)
		if err != nil {
			log.Printf("Error while creating token: %v", err)
			http.Error(w, "Failed to create token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	}
}

func VerifyJwtToken(a *Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			log.Println("Missing required 'token' parameter")
			http.Error(w, "Missing required 'token' parameter", http.StatusBadRequest)
			return
		}

		err := a.VerifyToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		log.Printf("successfully handle token verification %s", token)
		w.WriteHeader(http.StatusOK)
	}
}
