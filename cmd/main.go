package main

import (
	"dos/cfg"
	. "dos/db"
	"dos/internal"
	"dos/logger"
	"net/http"
	"os"
)

//	func logging(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			start := time.Now()
//			next.ServeHTTP(w, r)
//			log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
//		})
//	}
//
// FIXME POST CREATE USER
// curl --header "Content-Type: application/json" --request POST --data '{"username":"xyz"}' http://localhost:8080/user
// FIXME DELETE
// curl --request DELETE http://localhost:8080/user/xyz
func main() {
	config := cfg.LoadConfig()

	dbState := NewDBClient(config.Dsn)
	defer dbState.DB.Close()

	srv := &internal.Server{DB: dbState}

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

	mux.HandleFunc("/user/{name}", func(w http.ResponseWriter, r *http.Request) {

		logger.L.Info("inside GET user name")
		if r.Method == http.MethodDelete {
			srv.DeleteUser(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/db/disconnect", srv.DbDisconnect)
	mux.HandleFunc("/db/connect", srv.DbConnect)
	mux.HandleFunc("/db/status", srv.DbStatus)

	logger.L.Info("application run and listen on", "port", config.AppPort)

	if err := http.ListenAndServe(":"+config.AppPort, internal.LogMW(mux)); err != nil {
		logger.L.Error("http server stopped", "error", err)
		os.Exit(1)
	}
}
