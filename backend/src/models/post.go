package models

import (
	"database/sql"
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
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO user_posts (user_id, media, caption) VALUES (?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newPost.UserID, newPost.Media, newPost.Caption)

	if err != nil {
		return false, err
	}

	tx.Commit()

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

func GetUserPosts(count int, DB *sql.DB) ([]UserPost, error) {
	rows, err := DB.Query("SELECT id, user_id, media, caption, upload_time FROM user_posts LIMIT ?", count)
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
	sqlErr := row.QueryRow(id).Scan(singlePost.ID, &singlePost.UserID, &singlePost.Media, &singlePost.Caption, &singlePost.UploadTime)
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
