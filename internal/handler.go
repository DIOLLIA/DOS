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
	logger.L.Debug("[USER] GET")

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

	logger.L.Debug("[USER] POST")
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
	logger.L.Debug("[USER] DELETE")

	if err := DeleteUser(r.Context(), s.DB.DB); err != nil {
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
