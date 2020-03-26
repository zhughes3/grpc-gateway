package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	_ "github.com/lib/pq"
	"github.com/zhughes3/grpc/models"
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

func (db *db) updatePost(oldSlug, title, newSlug, html string) (*models.Post, error) {
	sql := `UPDATE posts
			SET title = $1, slug = $2, content = $3
			WHERE slug = $4
			RETURNING id;`
	var createdAt, updatedAt time.Time
	var id uint64

	err := db.db.QueryRow(sql, title, newSlug, html, oldSlug).Scan(&id)
	if err != nil {
		fmt.Println("Error updating Post")
		return nil, err
	}

	created, _ := convertTime(createdAt)
	updated, _ := convertTime(updatedAt)

	return &models.Post{
		Id:        id,
		Title:     title,
		Content:   html,
		Slug:      newSlug,
		CreatedAt: created,
		UpdatedAt: updated,
	}, nil
}

// do select first then update if returned or create new one
func (db *db) updateTags(postID uint64, tags []string) ([]string, error) {
	tsql := `SELECT tag FROM tags WHERE post_id = $1`
	rows, err := db.db.Query(tsql, postID)
	if err != nil {
		fmt.Println("Error selecting from tags")
		return nil, err
	}

	// tag --> bool (true value will be deleted from DB at the end of the function
	currentTags := make(map[string]bool)
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			fmt.Println(nil, err)
		}
		currentTags[t] = true
	}
	for _, t := range tags {
		if _, ok := currentTags[t]; ok {
			currentTags[t] = false
			continue
		}
		db.insertTag(t, postID)
	}

	for k, delete := range currentTags {
		if delete {
			var tid string
			sql := `DELETE from tags WHERE post_id = $1 and tag = $2 RETURNING id`
			err := db.db.QueryRow(sql, postID, k).Scan(&tid)
			if err != nil {
				fmt.Println("Error updating tag")
				return nil, err
			}
		}
	}

	return tags, nil
}

func (db *db) readPost(slug string) (*models.Post, error) {
	var post models.Post
	var createdAt, updatedAt time.Time
	sql := `SELECT * FROM posts WHERE slug = $1;`
	err := db.db.QueryRow(sql, slug).Scan(&post.Id, &post.Title, &post.Slug, &post.Content, &createdAt, &updatedAt)
	if err != nil {
		fmt.Println("Error reading post from db")
		return nil, err
	}

	if post.CreatedAt, err = convertTime(createdAt); err != nil {
		return nil, err
	}

	if post.UpdatedAt, err = convertTime(updatedAt); err != nil {
		return nil, err
	}

	sql = `SELECT id, tag FROM tags WHERE post_id = $1;`

	rows, err := db.db.Query(sql, post.GetId())
	if err != nil {
		fmt.Println("Error reading tags from db")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tag string
		var id int
		err := rows.Scan(&id, &tag)
		if err != nil {
			fmt.Println("Error reading individual row from tags query")
			return nil, err
		}
		post.Tags = append(post.Tags, tag)
	}

	return &post, nil
}

func (db *db) deletePost(slug string) (int, error) {
	var id int
	sql := `DELETE FROM posts WHERE slug = $1 RETURNING id;`
	err := db.db.QueryRow(sql, slug).Scan(&id)
	if err != nil {
		fmt.Println("Error deleting post from DB")
		return -1, err
	}
	return id, nil
}

func (db *db) insertPost(title, slug, html string) (*models.Post, error) {
	var createdAt, updatedAt time.Time
	var id uint64
	sql := `INSERT INTO posts(title, slug, content) VALUES($1, $2, $3) RETURNING id, created_at, updated_at;`

	err := db.db.QueryRow(sql, title, slug, html).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		fmt.Println("Error writing post to DB")
		return nil, err
	}

	created, _ := convertTime(createdAt)
	updated, _ := convertTime(updatedAt)

	return &models.Post{
		Id:        id,
		CreatedAt: created,
		UpdatedAt: updated,
	}, nil
}

func (db *db) insertTag(tag string, post_id uint64) (int, error) {
	var id int
	sql := `INSERT INTO tags(tag, post_id) VALUES($1, $2) RETURNING id`

	err := db.db.QueryRow(sql, tag, post_id).Scan(&id)
	if err != nil {
		fmt.Println("Error writing tag to DB")
		return -1, err
	}

	return id, nil
}

func convertTime(t time.Time) (*timestamp.Timestamp, error) {
	time, err := ptypes.TimestampProto(t)
	if err != nil {
		fmt.Println("Error converting time.time to timestamp.Timestamp")
		return nil, err
	}
	return time, nil
}
