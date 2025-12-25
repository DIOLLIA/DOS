package db

import (
	"context"
	"database/sql"
	dos "dos/internal"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

func OpenDB(dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	runMigrations(db)
	return db
}

func runMigrations(db *sql.DB) {
	const schema = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		username TEXT NOT NULL
	);
	`
	if _, err := db.Exec(schema); err != nil {
		log.Fatal(err)
	}
	log.Println("DB migrations OK")
}

func GetUser(ctx context.Context, db *sql.DB) (*dos.User, error) {
	row := db.QueryRowContext(ctx, `SELECT username FROM users WHERE id = 1`)
	var u dos.User
	if err := row.Scan(&u.Name); err != nil {
		return nil, err
	}
	return &u, nil
}

func UpsertUser(ctx context.Context, db *sql.DB, u dos.User) error {
	_, err := db.ExecContext(ctx, `
	INSERT INTO users (id, username)
	VALUES (1, $1)
	ON CONFLICT (id)
	DO UPDATE SET username = EXCLUDED.username
	`, u.Name)
	return err
}

func DeleteUser(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `DELETE FROM users WHERE id = 1`)
	return err
}
