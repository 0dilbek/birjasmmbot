package db

import (
	"database/sql"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"sort"

	_ "github.com/lib/pq"
)

func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	log.Println("Database connected")
	return db, nil
}

func Migrate(db *sql.DB) error {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(filename), "..", "..", "migrations")

	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil || len(files) == 0 {
		files, err = filepath.Glob("migrations/*.sql")
		if err != nil {
			return err
		}
	}

	sort.Strings(files)
	for _, f := range files {
		content, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}
		if _, err := db.Exec(string(content)); err != nil {
			return err
		}
		log.Printf("Applied migration: %s", filepath.Base(f))
	}
	return nil
}
