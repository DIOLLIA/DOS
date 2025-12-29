package main

import (
	"dos/cfg"
	. "dos/db"
	"dos/internal"
	"dos/logger"
	"net/http"
	"os"
)

func main() {
	config := cfg.LoadConfig()

	dbClient := NewDBClient(config.Dsn)
	defer dbClient.DB.Close()

	srv := &internal.Server{DB: dbClient}

	mux := http.NewServeMux()

	mux.HandleFunc("/entries", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			srv.GetEntries(w, r)
		case http.MethodPost:
			srv.PostEntry(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			srv.GetUser(w, r)
		case http.MethodPost:
			srv.PostUser(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/user/{name}", func(w http.ResponseWriter, r *http.Request) {
		username := r.PathValue("name") // dummy
		logger.L.Debug("path parameter handled", "username", username)
		if r.Method == http.MethodDelete {
			srv.DeleteUser(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/entries/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id") // dummy

		if r.Method == http.MethodDelete {
			srv.DeleteEntry(w, r, id)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/db/disconnect", srv.DbDisconnect)
	mux.HandleFunc("/db/connect", srv.DbConnect)
	mux.HandleFunc("/db/status", srv.DbStatus)

	logger.L.Info("application run and listen on", "port", config.AppPort)

	if err := http.ListenAndServe(":"+config.AppPort, internal.LogMW(internal.CorsMW(mux, config))); err != nil {
		logger.L.Error("http server stopped", "error", err)
		os.Exit(1)
	}
}
