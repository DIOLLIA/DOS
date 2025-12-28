package db

import (
	"context"
	"database/sql"
	"dos/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"sync/atomic"
)

type DBClient struct {
	DB        *sql.DB
	Connected atomic.Bool
}

type User struct {
	Name string `json:"username"`
}
type Entry struct {
	AdjNoun string `json:"entry"`
}

func NewDBClient(dsn string) *DBClient {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		logger.L.Error("cannot connect to database")
		panic(err)
	}

	if err := db.Ping(); err != nil {
		logger.L.Error("database not respond")
		panic(err)
	}

	runMigrations(db)

	state := &DBClient{DB: db}
	state.Connected.Store(true)
	return state
}

func (s *DBClient) IsConnected() bool {
	return s.Connected.Load()
}

func (s *DBClient) Connect() {
	s.Connected.Store(true)
	logger.L.Info("[DB] CONNECTED")
}

func (s *DBClient) Disconnect() {
	s.Connected.Store(false)
	logger.L.Info("[DB] DISCONNECTED")

}

func runMigrations(db *sql.DB) {
	const schema = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		username TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS entries (
		id SERIAL PRIMARY KEY,
		entry TEXT NOT NULL
	);
	`

	if _, err := db.Exec(schema); err != nil {
		log.Fatal(err)
	}
	log.Println("[DB] migrations OK")
}

func GetUser(ctx context.Context, db *sql.DB) (*User, error) {
	row := db.QueryRowContext(ctx, `SELECT username FROM users WHERE id = 1`)
	var u User
	if err := row.Scan(&u.Name); err != nil {
		return nil, err
	}
	return &u, nil
}

func GetEntries(ctx context.Context, db *sql.DB) (*[]string, error) {
	rows, err := db.QueryContext(ctx, `SELECT entry FROM entries`)
	if err != nil {
		logger.L.Error("query error" + err.Error())
	}
	var entryString string
	var entries []string
	for rows.Next() {
		if err := rows.Scan(&entryString); err != nil {
			return nil, err
		}
		entries = append(entries, entryString)
	}
	return &entries, nil
}

func PutUser(ctx context.Context, db *sql.DB, u User) error {
	_, err := db.ExecContext(ctx, `
	INSERT INTO users (id, username)
	VALUES (1, $1)
	ON CONFLICT (id)
	DO UPDATE SET username = EXCLUDED.username
	`, u.Name)
	return err
}

func PutEntry(ctx context.Context, db *sql.DB, entry string) error {
	logger.L.Info("[DB] PUT entry: " + entry)
	_, err := db.ExecContext(ctx, `
	INSERT INTO entries (entry)
	VALUES ($1)
	ON CONFLICT (id)
	DO UPDATE SET entry = EXCLUDED.entry
	`, entry)
	return err
}

func DeleteUser(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `DELETE FROM users WHERE id = 1`)
	return err
}

func DeleteEntry(ctx context.Context, db *sql.DB, name string) error {
	logger.L.Debug("Entry to be deleted", "entry", name)
	_, err := db.ExecContext(ctx, `DELETE FROM entries WHERE entry = $1`, name)
	return err
}
