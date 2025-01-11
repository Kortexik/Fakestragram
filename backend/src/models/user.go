package models

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id        int    `gorm:"primaryKey" json:"id"`
	Username  string `gorm:"type:varchar(25)" json:"username"`
	Password  string `gorm:"type:varchar(35)" json:"password"`
	FirstName string `gorm:"type:varchar(100)" json:"first_name"`
	LastName  string `gorm:"type:varchar(100)" json:"last_name"`
	Email     string `gorm:"type:varchar(100);uniqueIndex" json:"email"`
}

func GetUsers(count int, DB *sql.DB) ([]User, error) {

	rows, err := DB.Query("SELECT id, username, first_name, last_name, email from users LIMIT " + strconv.Itoa(count))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]User, 0)

	for rows.Next() {
		singleUser := User{}
		err = rows.Scan(&singleUser.Id, &singleUser.Username, &singleUser.FirstName, &singleUser.LastName, &singleUser.Email)

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
	row, err := DB.Prepare("SELECT id, username, first_name, last_name, email FROM users WHERE id = ?")

	if err != nil {
		return User{}, err
	}

	user := User{}

	sqlErr := row.QueryRow(id).Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName, &user.Email)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return User{}, nil
		}
		return User{}, sqlErr
	}

	return user, nil
}

func AddUser(newUser User, DB *sql.DB) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO users (username, first_name, last_name, email) VALUES (?, ?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newUser.Username, newUser.FirstName, newUser.LastName, newUser.Email)

	if err != nil {
		return false, err
	}

	tx.Commit()

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

func UpdateUser(id string, newUserData User, DB *sql.DB) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare("UPDATE users SET username = ?, first_name = ?, last_name = ?, email = ? WHERE id = ?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(newUserData.Username, newUserData.FirstName, newUserData.LastName, newUserData.Email, id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsAffected == 0 {
		return false, fmt.Errorf("user with ID %s not found", id)
	}

	tx.Commit()

	return true, nil
}
