package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type db struct {
	db *sql.DB
}

func NewDB(cfg *dbConfig) *sql.DB {
	dbinfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.user, cfg.password, cfg.host, cfg.port, cfg.name, "disable")
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

func (db *db) insertPost(title, slug, html string) (int, error) {
	var id int
	sql := `INSERT INTO posts(title, slug, content) VALUES($1, $2, $3) RETURNING id`

	err := db.db.QueryRow(sql, title, slug, html).Scan(&id)
	if err != nil {
		fmt.Println("Error writing post to DB")
		return -1, err
	}

	return id, nil
}

func (db *db) insertTag(tag string, post_id int) (int, error) {
	var id int
	sql := `INSERT INTO tags(tag, post_id) VALUES($1, $2) RETURNING id`

	err := db.db.QueryRow(sql, tag, post_id).Scan(&id)
	if err != nil {
		fmt.Println("Error writing tag to DB")
		return -1, err
	}

	return id, nil
}
