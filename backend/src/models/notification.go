package models

import (
	"database/sql"
	"log"
	"time"
)

type Notification struct {
	ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int       `gorm:"column:user_id" json:"user_id"`
	Type          string    `gorm:"type:varchar(50)" json:"type"`
	RelatedUserID int       `gorm:"column:related_user_id" json:"related_user_id"`
	PostID        int       `gorm:"column:post_id" json:"post_id"`
	Content       string    `gorm:"type:text" json:"content"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	Seen          bool      `gorm:"column:seen" json:"seen"`
}

func GetUnseenNotifications(userID int, DB *sql.DB) ([]Notification, error) {
	rows, err := DB.Query("SELECT id, user_id, type, related_user_id, post_id, content, created_at, seen FROM notifications WHERE user_id = ? AND seen = 0", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := make([]Notification, 0)
	for rows.Next() {
		var notification Notification
		err = rows.Scan(&notification.ID, &notification.UserID, &notification.Type, &notification.RelatedUserID, &notification.PostID, &notification.Content, &notification.CreatedAt, &notification.Seen)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func AddNotification(newNotification Notification, tx *sql.Tx) (bool, error) {
	stmt, err := tx.Prepare("INSERT INTO notifications (user_id, type, related_user_id, post_id, content) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Error preparing statement for notification:", err)
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newNotification.UserID, newNotification.Type, newNotification.RelatedUserID, newNotification.PostID, newNotification.Content)
	if err != nil {
		log.Println("Error executing notification statement for user ID", newNotification.UserID, ":", err)
		return false, err
	}

	return true, nil
}

func MarkNotificationsAsSeen(userID int, DB *sql.DB) error {
	_, err := DB.Exec("UPDATE notifications SET seen = 1 WHERE user_id = ? AND seen = 0", userID)
	return err
}
