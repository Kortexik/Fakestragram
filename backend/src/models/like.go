package models

import (
	"database/sql"
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
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO likes (user_id, post_id) VALUES (?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newLike.UserID, newLike.PostID)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func DeleteLike(id string, DB *sql.DB) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare("DELETE FROM likes WHERE id = ?")
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
