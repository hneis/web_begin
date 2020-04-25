package models

import (
	"database/sql"
)

type PostItem struct {
	ID      string
	Title   string
	Text    string
	Created string
	Author  string
}

type PostItemSlice []PostItem

func (post *PostItem) Insert(db *sql.DB) error {
	_, err := db.Exec(
		"INSERT INTO POSTS (ID, Title, Text, Created, Author) VALUES( ?, ?, ?, ?, ?)",
		post.ID, post.Title, post.Text, post.Created, post.Author)
	return err
}

func (post *PostItem) Delete(db *sql.DB) error {
	_, err := db.Exec(
		"DELETE FROM POSTS WHERE ID = ?",
		post.ID,
	)

	return err
}

func (post *PostItem) Update(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE POSTS SET TITLE= ?, TEXT = ? WHERE ID = ?",
		post.Title, post.Text, post.ID,
	)

	return err
}

func (post *PostItem) Get(db *sql.DB) error {
	row := db.QueryRow("SELECT ID, TITLE, TEXT,  CREATED, AUTHOR From POSTS WHERE ID = ?", post.ID)
	if err := row.Scan(&post.ID, &post.Title, &post.Text, &post.Created, &post.Author); err != nil {
		return err
	}
	return nil
}

func GetAllPostItems(db *sql.DB) (PostItemSlice, error) {
	rows, err := db.Query("SELECT ID, TITLE, TEXT,  CREATED, AUTHOR From POSTS")
	if err != nil {
		return nil, err
	}

	posts := make(PostItemSlice, 0, 8)
	for rows.Next() {
		post := PostItem{}
		if err = rows.Scan(&post.ID, &post.Title, &post.Text, &post.Created, &post.Author); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
