package models

import (
	"database/sql"
	"fmt"
	"log"
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

func AddComment(newComment Comment, DB *sql.DB) (*Comment, error) {
	tx, err := DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return nil, err
	}

	stmt, err := tx.Prepare("INSERT INTO comments (user_id, post_id, content) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("Error preparing statement:", err)
		tx.Rollback()
		return nil, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(newComment.UserID, newComment.PostID, newComment.Content)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	post, err := GetUserPostById(strconv.Itoa(newComment.PostID), DB)
	if err != nil {
		log.Println("Error getting userpost: ", err)
		return nil, err
	}

	relatedUserUsername, err := GetUsernameById(newComment.UserID, DB)
	if err != nil {
		log.Println("Error getting username: ", err)
		return nil, err
	}

	if newComment.UserID != post.UserID {
		newNotification := Notification{
			UserID:        post.UserID,
			Type:          "Like",
			RelatedUserID: newComment.UserID,
			PostID:        newComment.PostID,
			Content:       fmt.Sprintf("%s has commented on post.", relatedUserUsername),
		}

		created, err := AddNotification(newNotification, tx)
		if err != nil || !created {
			log.Println("Error adding notification:", err)
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
		return nil, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	var createdComment Comment
	row := DB.QueryRow("SELECT id, user_id, post_id, content, created_at FROM comments WHERE id = ?", lastInsertID)
	if err := row.Scan(&createdComment.ID, &createdComment.UserID, &createdComment.PostID, &createdComment.Content, &createdComment.CreatedAt); err != nil {
		return nil, err
	}

	return &createdComment, nil
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
	fmt.Println(comments)

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return comments, err
}
