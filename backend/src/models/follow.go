package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Follow struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	FollowerID int       `gorm:"column:follower_id" json:"follower_id"`
	FolloweeID int       `gorm:"column:followee_id" json:"followee_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func AddFollow(follow Follow, DB *sql.DB) (*Follow, error) {
	tx, err := DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return nil, err
	}

	stmt, err := tx.Prepare("INSERT INTO follows (follower_id, following_id) VALUES (?, ?)")
	if err != nil {
		log.Println("Error preparing statement:", err)
		tx.Rollback()
		return nil, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(follow.FollowerID, follow.FolloweeID)
	if err != nil {
		log.Println("Error executing statement:", err)
		tx.Rollback()
		return nil, err
	}

	relatedUserUsername, err := GetUsernameById(follow.FollowerID, DB)
	if err != nil {
		log.Println("Error getting username:", err)
		tx.Rollback()
		return nil, err
	}

	newNotification := Notification{
		UserID:        follow.FolloweeID,
		Type:          "Follow",
		RelatedUserID: follow.FollowerID,
		Content:       fmt.Sprintf("%s has started following you.", relatedUserUsername),
	}

	created, err := AddNotification(newNotification, tx)
	if err != nil || !created {
		log.Println("Error adding notification:", err)
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
		return nil, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	var createdFollow Follow
	row := DB.QueryRow("SELECT id, follower_id, following_id, created_at FROM follows WHERE id = ?", lastInsertID)
	if err := row.Scan(&createdFollow.ID, &createdFollow.FollowerID, &createdFollow.FolloweeID, &createdFollow.CreatedAt); err != nil {
		return nil, err
	}

	return &createdFollow, nil
}

func DeleteFollow(follower_id int, followee_id int, DB *sql.DB) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare("DELETE FROM follows WHERE follower_id = ? AND following_id = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(follower_id, followee_id)
	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func FollowExists(followerID int, followeeID int, DB *sql.DB) (bool, error) {
	query := `SELECT EXISTS (
        SELECT 1 
        FROM follows 
        WHERE follower_id = $1 AND following_id = $2
    )`

	var exists bool
	err := DB.QueryRow(query, followerID, followeeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking follow relationship: %v", err)
	}

	return exists, nil
}

func NumberOfFollowers(userID int, DB *sql.DB) (int, error) {
	query := `SELECT COUNT(*) FROM follows WHERE following_id = $1`

	var countOfFollowers int
	err := DB.QueryRow(query, userID).Scan(&countOfFollowers)
	if err != nil {
		return -1, fmt.Errorf("error fetching number of followers: %v", err)
	}
	return countOfFollowers, nil
}

func NumberOfFollowees(userID int, DB *sql.DB) (int, error) {
	query := `SELECT COUNT(*) FROM follows WHERE follower_id = $1`

	var countOfFollowees int
	err := DB.QueryRow(query, userID).Scan(&countOfFollowees)
	if err != nil {
		return -1, fmt.Errorf("error fetching number of followees: %v", err)
	}
	return countOfFollowees, nil
}
