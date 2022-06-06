package repository

import (
	"chat/models"
	"database/sql"
)

const (
	// Defind SQL commnad
	InsertUser      = "INSERT INTO user(id, name) values (?,?)"
	DeleteUser      = "DELETE FROM user WHERE id = ?"
	FindUserByIDSQL = "SELECT id,name FROM user WHERE id = ? LIMIT 1"
	GetAllUsersList = "SELETE id,name FROM user"
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
	stmt, err := repo.Db.Prepare(InsertUser)
	_, err = stmt.Exec(user.GetID(), user.GetName())

	HandleError(err, "Values to User table")
}

// Remove user
func (repo *UserRepository) RemoveUser(user models.User) {
	stmt, err := repo.Db.Prepare(DeleteUser)
	_, err = stmt.Exec(user.GetID())

	HandleError(err, "Delete User table")
}

// Find user by ID
func (repo *UserRepository) FindUserByID(ID string) models.User {
	row := repo.Db.QueryRow(FindUserByIDSQL, ID)

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
	rows, err := repo.Db.Query(GetAllUsersList)

	// Handle error
	HandleError(err, "List all user")

	// Closr rows
	var users []models.User
	defer rows.Close()

	// Scan users and append it into list
	for rows.Next() {
		var user User
		rows.Scan(&user.ID, &user.Name)
		users = append(users, &user)
	}

	return users
}
