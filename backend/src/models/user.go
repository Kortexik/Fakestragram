package models

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         int     `gorm:"primaryKey" json:"id"`
	Username   string  `gorm:"type:varchar(35)" json:"username"`
	Password   []byte  `gorm:"type:binary(60)" json:"password"`
	FirstName  string  `gorm:"type:varchar(30)" json:"first_name"`
	LastName   string  `gorm:"type:varchar(40)" json:"last_name"`
	Email      string  `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Bio        *string `gorm:"type:varchar(150)" json:"bio"`
	ProfilePic []byte  `gorm:"type:blob" json:"profile_pic"`
	CreatedAt  string  `gorm:"autoCreateTime" json:"created_at"`
}

func HashPassword(password []byte) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func GetUsers(count int, DB *sql.DB) ([]User, error) {

	rows, err := DB.Query("SELECT id, username, password, first_name, last_name, email, bio, profile_pic, created_at from users LIMIT " + strconv.Itoa(count))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]User, 0)

	for rows.Next() {
		singleUser := User{}
		err = rows.Scan(&singleUser.ID, &singleUser.Username, &singleUser.Password, &singleUser.FirstName, &singleUser.LastName, &singleUser.Email, &singleUser.Bio, &singleUser.ProfilePic, &singleUser.CreatedAt)

		if err != nil {
			return nil, err
		}

		users = append(users, singleUser)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return users, err
}

func GetUserById(id string, DB *sql.DB) (User, error) {
	row, err := DB.Prepare("SELECT id, username, password, first_name, last_name, email, bio, profile_pic, created_at from users WHERE id = ?")

	if err != nil {
		return User{}, err
	}

	user := User{}

	sqlErr := row.QueryRow(id).Scan(&user.ID, &user.Username, &user.Password, &user.FirstName, &user.LastName, &user.Email, &user.Bio, &user.ProfilePic, &user.CreatedAt)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return User{}, nil
		}
		return User{}, sqlErr
	}

	return user, nil
}

func GetUserProfileDataByUsername(username string, DB *sql.DB) (User, error) {
	row, err := DB.Prepare("SELECT id, username, first_name, last_name, bio, profile_pic, created_at from users WHERE username = ?")

	if err != nil {
		return User{}, err
	}

	user := User{}

	sqlErr := row.QueryRow(username).Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Bio, &user.ProfilePic, &user.CreatedAt)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return User{}, nil
		}
		return User{}, sqlErr
	}
	return user, nil
}

func GetUsernameById(id int, DB *sql.DB) (string, error) {
	row, err := DB.Prepare("SELECT username from users WHERE id = ?")
	if err != nil {
		return "", err
	}

	var username string
	sqlErr := row.QueryRow(id).Scan(&username)
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return "", nil
		}
		return "", sqlErr
	}

	return username, nil
}

func AddUser(newUser User, DB *sql.DB) (bool, error) {
	// Hash the user's password
	hashedPassword, err := HashPassword(newUser.Password)
	if err != nil {
		return false, fmt.Errorf("failed to hash password: %v", err)
	}

	// Begin a transaction
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO users (username, password, first_name, last_name, email) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	// Execute the statement with the hashed password
	_, err = stmt.Exec(newUser.Username, hashedPassword, newUser.FirstName, newUser.LastName, newUser.Email)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteUser(id string, DB *sql.DB) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare("DELETE FROM users WHERE id = ?")
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

func GetUserIdByUsername(username string, DB *sql.DB) (int, error) {
	row, err := DB.Prepare("SELECT id from users WHERE username = ?")

	if err != nil {
		return 0, err
	}

	var userID int

	sqlErr := row.QueryRow(username).Scan(&userID)
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return 0, nil
		}
		return 0, sqlErr
	}
	return userID, nil
}

func GetUserFollowers(userID int, DB *sql.DB) ([]User, error) {
	rows, err := DB.Query("SELECT id, username, password, first_name, last_name, email, bio, profile_pic, created_at FROM users WHERE id IN (SELECT follower_id FROM follows WHERE following_id = ?)", userID)
	if err != nil {
		log.Println("Error querying followers for user ID", userID, ":", err)
		return nil, err
	}
	defer rows.Close()

	followers := make([]User, 0)

	for rows.Next() {
		singleFollower := User{}
		err = rows.Scan(&singleFollower.ID, &singleFollower.Username, &singleFollower.Password, &singleFollower.FirstName, &singleFollower.LastName, &singleFollower.Email, &singleFollower.Bio, &singleFollower.ProfilePic, &singleFollower.CreatedAt)
		if err != nil {
			log.Println("Error scanning follower row for user ID", userID, ":", err)
			return nil, err
		}
		followers = append(followers, singleFollower)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over follower rows for user ID", userID, ":", err)
		return nil, err
	}

	return followers, nil
}

func UpdateProfilePic(editData User, DB *sql.DB) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE users SET profile_pic = ? WHERE id = ?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(editData.ProfilePic, editData.ID)
	if err != nil {
		return false, err
	}

	tx.Commit()
	return true, nil
}

func UpdateBio(editData User, DB *sql.DB) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE users SET bio = ? WHERE id = ?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(*editData.Bio, editData.ID)
	if err != nil {
		return false, err
	}

	tx.Commit()
	return true, nil
}
