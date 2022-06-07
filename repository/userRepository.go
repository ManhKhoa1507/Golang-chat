package repository

import (
	"chat/models"
	"database/sql"
	"fmt"
	"log"
)

const (
	// Defind SQL commnad
	InsertUser      = "INSERT INTO user (id, name) values (?,?)"
	DeleteUser      = "DELETE FROM user WHERE id = ?"
	FindUserByIDSQL = "SELECT id,name FROM user WHERE id = ? LIMIT 1"
	GetAllUsersList = "SELECT id, name FROM user"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Get user ID
func (user *User) GetID() string {
	return user.ID
}

// Get user name
func (user *User) GetName() string {
	return user.Name
}

type UserRepository struct {
	Db *sql.DB
}

// Add user value to user table
func (repo *UserRepository) AddUser(user models.User) {
	stmt, err := repo.Db.Prepare("INSERT INTO user (id, name) values (?,?)")
	_, err = stmt.Exec(user.GetID(), user.GetName())

	if err != nil {
		fmt.Println("Error insert value to user table")
	}
}

// Remove user
func (repo *UserRepository) RemoveUser(user models.User) {
	stmt, err := repo.Db.Prepare("DELETE FROM user WHERE id = ?")
	_, err = stmt.Exec(user.GetID())

	HandleError(err, "Delete User table")
}

// Find user by ID
func (repo *UserRepository) FindUserByID(ID string) models.User {
	row := repo.Db.QueryRow("SELECT id, name FROM user WHERE id = ? LIMIT 1", ID)

	var user User

	if err := row.Scan(&user.ID, &user.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}

	return &user
}

// Get all user in database
func (repo *UserRepository) GetAllUsers() []models.User {
	rows, err := repo.Db.Query("SELECT id, name FROM user")
	// rows, err := repo.Db.Prepare(GetAllUsersList)
	// rows.Exec()
	// // Handle error
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	// Closr rows
	var users []models.User

	// Scan users and append it into list
	for rows.Next() {
		var user User
		rows.Scan(&user.ID, &user.Name)
		users = append(users, &user)
	}

	return users
}
