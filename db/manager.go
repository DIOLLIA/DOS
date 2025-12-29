package db

import (
	"context"
	"database/sql"
	"dos/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"sync/atomic"
)

type Client struct {
	DB        *sql.DB
	Connected atomic.Bool
}

type User struct {
	Name string `json:"username"`
}
type Entry struct {
	Id    int    `json:"id"`
	Value string `json:"value"`
}

func NewDBClient(dsn string) *Client {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		logger.L.Error("[DB] cannot connect to database")
		panic(err)
	}

	if err := db.Ping(); err != nil {
		logger.L.Error("[DB] database not respond")
		panic(err)
	}

	runMigrations(db)

	state := &Client{DB: db}
	state.Connected.Store(true)
	return state
}

func (s *Client) IsConnected() bool {
	return s.Connected.Load()
}

func (s *Client) Connect() {
	s.Connected.Store(true)
	logger.L.Info("[DB] CONNECTED")
}

func (s *Client) Disconnect() {
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
		logger.L.Error("[DB] error on db migration", "error", err.Error())
		panic(err)
	}
	logger.L.Info("[DB] migrations OK")
}

func GetUser(ctx context.Context, db *sql.DB) (*User, error) {
	row := db.QueryRowContext(ctx, `SELECT username FROM users WHERE id = 1`)
	var u User
	if err := row.Scan(&u.Name); err != nil {
		logger.L.Error("[DB] get user query error" + err.Error())

		return nil, err
	}

	return &u, nil
}

func GetEntries(ctx context.Context, db *sql.DB) ([]Entry, error) {
	rows, err := db.QueryContext(ctx, `SELECT id, entry FROM entries`)
	if err != nil {
		logger.L.Error("[DB] get entries query error" + err.Error())
	}
	var entries []Entry
	defer rows.Close()

	for rows.Next() {
		var entry Entry
		if err := rows.Scan(&entry.Id, &entry.Value); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func PutUser(ctx context.Context, db *sql.DB, u User) error {
	_, err := db.ExecContext(ctx, `
	INSERT INTO users (id, username)
	VALUES (1, $1)
	ON CONFLICT (id)
	DO UPDATE SET username = EXCLUDED.username
	`, u.Name)

	if err != nil {
		logger.L.Error("[DB] put user query error" + err.Error())
	}
	return err
}

func PutEntry(ctx context.Context, db *sql.DB, entry string) (int64, error) {
	logger.L.Info("[DB] PUT entry: " + entry)

	var id int64
	err := db.QueryRowContext(ctx, `
	INSERT INTO entries (entry)
	VALUES ($1)
RETURNING id`, entry).Scan(&id)

	if err != nil {
		logger.L.Error("[DB] put entry query error" + err.Error())
	}
	return id, err
}

func DeleteUser(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `DELETE FROM users WHERE id = 1`)
	if err != nil {
		logger.L.Error("[DB] delete user query error" + err.Error())
	}
	return err
}

func DeleteEntry(ctx context.Context, db *sql.DB, id string) error {
	_, err := db.ExecContext(ctx, `DELETE FROM entries WHERE id = $1`, id)
	if err != nil {
		logger.L.Error("[DB] delete entry query error" + err.Error())
	}
	return err
}
