package internal

import (
	"database/sql"
	"dos/db"
	"encoding/json"
	"errors"
	"net/http"
)

type Server struct {
	Database *sql.DB
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUser(r.Context(), s.Database)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (s *Server) PostUser(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if u.Name == "" {
		http.Error(w, "username required", http.StatusBadRequest)
		return
	}

	if err := db.UpsertUser(r.Context(), s.Database, u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if err := db.DeleteUser(r.Context(), s.Database); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
