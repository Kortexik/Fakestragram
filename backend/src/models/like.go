package models

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

type Like struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int       `gorm:"column:user_id" json:"user_id"`
	PostID    int       `gorm:"column:post_id" json:"post_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func AddLike(newLike Like, DB *sql.DB) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO likes (user_id, post_id) VALUES (?, ?)")

	if err != nil {
		log.Println("Error preparing statement:", err)
		tx.Rollback()
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newLike.UserID, newLike.PostID)

	if err != nil {
		log.Println("Error inserting like:", err)
		return false, err
	}

	post, err := GetUserPostById(strconv.Itoa(newLike.PostID), DB)
	if err != nil {
		log.Println("Error getting userpost: ", err)
		return false, err
	}

	relatedUserUsername, err := GetUsernameById(newLike.UserID, DB)
	if err != nil {
		log.Println("Error getting username: ", err)
		return false, err
	}

	if newLike.UserID != post.UserID {
		newNotification := Notification{
			UserID:        post.UserID,
			Type:          "Like",
			RelatedUserID: newLike.UserID,
			PostID:        newLike.PostID,
			Content:       fmt.Sprintf("%s has liked your post.", relatedUserUsername),
		}

		created, err := AddNotification(newNotification, tx)
		if err != nil || !created {
			log.Println("Error adding notification:", err)
			tx.Rollback()
			return false, err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
		return false, err
	}

	return true, nil
}

func DeleteLike(user_id int, post_id int, DB *sql.DB) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare("DELETE FROM likes WHERE user_id = ? AND post_id = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(user_id, post_id)
	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func GetLikes(count int, DB *sql.DB) ([]Like, error) {

	rows, err := DB.Query("SELECT id, user_id, post_id, created_at FROM likes LIMIT " + strconv.Itoa(count))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	likes := make([]Like, 0)

	for rows.Next() {
		singleLike := Like{}
		err = rows.Scan(&singleLike.ID, &singleLike.UserID, &singleLike.PostID, &singleLike.CreatedAt)

		if err != nil {
			return nil, err
		}

		likes = append(likes, singleLike)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return likes, err
}

func GetLikesByPost(post_id string, DB *sql.DB) ([]Like, error) {

	rows, err := DB.Query("SELECT id, user_id, post_id, created_at FROM likes WHERE post_id = " + post_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	likes := make([]Like, 0)

	for rows.Next() {
		singleLike := Like{}
		err = rows.Scan(&singleLike.ID, &singleLike.UserID, &singleLike.PostID, &singleLike.CreatedAt)

		if err != nil {
			return nil, err
		}

		likes = append(likes, singleLike)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return likes, err
}
