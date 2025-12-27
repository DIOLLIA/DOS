package internal

import (
	"database/sql"
	. "dos/db"
	"dos/logger"
	"encoding/json"
	"errors"
	"net/http"
)

type Server struct {
	DB *DBClient
}

func (s *Server) isConnected(w http.ResponseWriter, r *http.Request) bool {
	if !s.DB.IsConnected() {
		logger.L.Error("db was disconnected")
		http.Error(w, "database connection failed", http.StatusInternalServerError)
		return false
	}
	return true
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	if !s.isConnected(w, r) {
		return
	}
	logger.L.Debug("[USER] GET invoked")

	user, err := GetUser(r.Context(), s.DB.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.L.Info("no user found", "status", http.StatusNotFound)
			//w.WriteHeader(http.StatusNotFound) //todo if no body needed, delete http.notFound(w,r) and uncomment this line
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
	if !s.isConnected(w, r) {
		return
	}

	logger.L.Debug("[USER] POST invoked")
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := PutUser(r.Context(), s.DB.DB, u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if !s.isConnected(w, r) {
		return
	}
	logger.L.Debug("[USER] DELETE invoked")

	if err := DeleteUser(r.Context(), s.DB.DB); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) GetEntries(w http.ResponseWriter, r *http.Request) {
	if !s.isConnected(w, r) {
		return
	}
	logger.L.Debug("[ENTRIES] GET")

	entries, err := GetEntries(r.Context(), s.DB.DB)
	if err != nil {
		logger.L.Error(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			logger.L.Info("no entries found", "status", http.StatusNotFound)
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	logger.L.Debug("[ENTRIES] GET result", "entries", entries)
	json.NewEncoder(w).Encode(entries)
}

func (s *Server) PostEntry(w http.ResponseWriter, r *http.Request) {
	if !s.isConnected(w, r) {
		return
	}

	logger.L.Debug("[ENTRIES] POST invoked")
	var entryStr Entry
	var entry string
	if err := json.NewDecoder(r.Body).Decode(&entryStr); err != nil {
		logger.L.Error(err.Error())
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := PutEntry(r.Context(), s.DB.DB, entryStr.AdjNoun); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entry)
}

func (s *Server) DeleteEntry(w http.ResponseWriter, r *http.Request, entry string) {
	if !s.isConnected(w, r) {
		return
	}
	logger.L.Debug(" [ENTRIES] DELETE invoked")

	if err := DeleteEntry(r.Context(), s.DB.DB, entry); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) DbDisconnect(w http.ResponseWriter, r *http.Request) {
	s.DB.Disconnect()
	w.WriteHeader(http.StatusOK)
}

func (s *Server) DbConnect(w http.ResponseWriter, r *http.Request) {
	s.DB.Connect()
	w.WriteHeader(http.StatusOK)
}

func (s *Server) DbStatus(w http.ResponseWriter, r *http.Request) {
	if s.DB.IsConnected() {
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Error(w, "db disconnected", http.StatusInternalServerError)
}
