package main

import (
	. "dos/db"
	. "dos/internal"
	"log"
	"net/http"
	"os"
	"time"
)

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db := OpenDB(dsn)
	defer db.Close()

	srv := &Server{Database: db}

	mux := http.NewServeMux()
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

	mux.HandleFunc("/user/name", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			srv.DeleteUser(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", logging(mux)))
}
