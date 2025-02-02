package models

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

type UserPost struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int       `gorm:"column:user_id" json:"user_id"`
	Media      []byte    `gorm:"type:blob" json:"media"`
	Caption    string    `gorm:"type:text" json:"caption"`
	UploadTime time.Time `gorm:"autoCreateTime" json:"upload_time"`
	Likes      []Like    `gorm:"foreignKey:PostID" json:"likes"`
	Comments   []Comment `gorm:"foreignKey:PostID" json:"comments"`
}

func AddUserPost(newPost UserPost, DB *sql.DB) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO user_posts (user_id, media, caption) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("Error preparing statement for user post:", err)
		tx.Rollback()
		return false, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(newPost.UserID, newPost.Media, newPost.Caption)
	if err != nil {
		log.Println("Error executing statement for user post:", err)
		tx.Rollback()
		return false, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error fetching last insert ID for user post:", err)
		tx.Rollback()
		return false, err
	}
	newPost.ID = int(postID)

	followers, err := GetUserFollowers(newPost.UserID, DB)
	if err != nil {
		log.Println("Error fetching followers for user ID", newPost.UserID, ":", err)
		tx.Rollback()
		return false, err
	}

	if len(followers) == 0 {
		if err := tx.Commit(); err != nil {
			log.Println("Error committing transaction:", err)
			return false, err
		}
		return true, nil
	}

	relatedUserUsername, err := GetUsernameById(newPost.UserID, DB)
	if err != nil {
		log.Println("Error fetching username for user ID", newPost.UserID, ":", err)
		tx.Rollback()
		return false, err
	}
	for _, follower := range followers {
		newNotification := Notification{
			UserID:        follower.ID,
			Type:          "UserPost",
			RelatedUserID: newPost.UserID,
			PostID:        newPost.ID,
			Content:       fmt.Sprintf("%s has added a new post.", relatedUserUsername),
		}

		created, err := AddNotification(newNotification, tx)
		if err != nil || !created {
			log.Println("Error adding notification for follower ID", follower.ID, ":", err)
			tx.Rollback()
			return false, err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction for user post:", err)
		return false, err
	}

	return true, nil
}

func DeleteUserPost(id string, DB *sql.DB) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare("DELETE FROM user_posts WHERE id = ?")
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

func GetUserPosts(DB *sql.DB) ([]UserPost, error) {
	rows, err := DB.Query("SELECT id, user_id, media, caption, upload_time FROM user_posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]UserPost, 0)

	for rows.Next() {
		singlePost := UserPost{}
		err = rows.Scan(&singlePost.ID, &singlePost.UserID, &singlePost.Media, &singlePost.Caption, &singlePost.UploadTime)
		if err != nil {
			return nil, err
		}

		postID := strconv.Itoa(singlePost.ID)

		singlePost.Likes, err = GetLikesByPost(postID, DB)
		if err != nil {
			return nil, err
		}

		singlePost.Comments, err = GetCommentsByPost(postID, DB)
		if err != nil {
			return nil, err
		}

		posts = append(posts, singlePost)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetUserPostsByUserId(user_id int, DB *sql.DB) ([]UserPost, error) {
	rows, err := DB.Query("SELECT id, user_id, media, caption, upload_time FROM user_posts WHERE user_id = ?", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]UserPost, 0)

	for rows.Next() {
		singlePost := UserPost{}
		err = rows.Scan(&singlePost.ID, &singlePost.UserID, &singlePost.Media, &singlePost.Caption, &singlePost.UploadTime)
		if err != nil {
			return nil, err
		}

		postID := strconv.Itoa(singlePost.ID)

		singlePost.Likes, err = GetLikesByPost(postID, DB)
		if err != nil {
			return nil, err
		}

		singlePost.Comments, err = GetCommentsByPost(postID, DB)
		if err != nil {
			return nil, err
		}

		posts = append(posts, singlePost)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetUserPostById(id string, DB *sql.DB) (UserPost, error) {

	row, err := DB.Prepare("SELECT id, user_id, media, caption, upload_time from user_posts WHERE id = ?")
	if err != nil {
		return UserPost{}, err
	}

	singlePost := UserPost{}
	sqlErr := row.QueryRow(id).Scan(&singlePost.ID, &singlePost.UserID, &singlePost.Media, &singlePost.Caption, &singlePost.UploadTime)
	postID := strconv.Itoa(singlePost.ID)
	singlePost.Likes, err = GetLikesByPost(postID, DB)
	if err != nil {
		return UserPost{}, err
	}
	singlePost.Comments, err = GetCommentsByPost(postID, DB)
	if err != nil {
		return UserPost{}, err
	}

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return UserPost{}, nil
		}
		return UserPost{}, sqlErr
	}

	return singlePost, err
}
