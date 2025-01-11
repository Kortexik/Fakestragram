package models

import (
	"database/sql"
	"strconv"
	"time"
)

type Comment struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int       `gorm:"column:user_id" json:"user_id"`
	PostID    int       `gorm:"column:post_id" json:"post_id"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func AddComment(newComment Comment, DB *sql.DB) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO comments (user_id, post_id, content) VALUES (?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newComment.UserID, newComment.PostID, newComment.Content)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func DeleteComment(id string, DB *sql.DB) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare("DELETE FROM comments WHERE id = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func GetComments(count int, DB *sql.DB) ([]Comment, error) {

	rows, err := DB.Query("SELECT id, user_id, post_id, content, created_at FROM comments LIMIT " + strconv.Itoa(count))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := make([]Comment, 0)

	for rows.Next() {
		singleComment := Comment{}
		err = rows.Scan(&singleComment.ID, &singleComment.UserID, &singleComment.PostID, &singleComment.Content, &singleComment.CreatedAt)

		if err != nil {
			return nil, err
		}

		comments = append(comments, singleComment)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return comments, err
}

func GetCommentsByPost(post_id string, DB *sql.DB) ([]Comment, error) {

	rows, err := DB.Query("SELECT id, user_id, post_id, content, created_at FROM comments WHERE post_id = " + post_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := make([]Comment, 0)

	for rows.Next() {
		singleComment := Comment{}
		err = rows.Scan(&singleComment.ID, &singleComment.UserID, &singleComment.PostID, &singleComment.Content, &singleComment.CreatedAt)

		if err != nil {
			return nil, err
		}

		comments = append(comments, singleComment)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return comments, err
}
