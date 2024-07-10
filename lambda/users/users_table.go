package main

import (
	"credential-store/database"
	"database/sql"
)

type IUsersTable interface {
	PutUser(*User) error
	GetUser(string) (*GetUsersResponse, error)
	GetAllUsers() ([]*GetUsersResponse, error)
	UpdateUser(*User) error
	DeleteUser(*User) error
}

type UsersTable struct {
	db *sql.DB
}

func NewUsersTable() (*UsersTable, error) {
	db, err := database.NewCredentialStoreDB()

	if err != nil {
		return nil, err
	}

	return &UsersTable{db: db.GetDatabase()}, nil
}

func (d *UsersTable) GetUser(userName string) (*GetUsersResponse, error) {
	row := d.db.QueryRow("SELECT username, email FROM users WHERE username=$1;", userName)
	user := &GetUsersResponse{}
	err := row.Scan(&user.Username, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (c *UsersTable) GetAllUsers() ([]*GetUsersResponse, error) {
	rows, err := c.db.Query("SELECT username, email FROM users;")
	if err != nil {
		return nil, err
	}

	users := []*GetUsersResponse{}
	for rows.Next() {
		user := GetUsersResponse{}
		err := rows.Scan(
			&user.Username,
			&user.Email)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (c *UsersTable) PutUser(user *User) error {
	_, err := c.db.Exec(`
	INSERT INTO users(username, password, email, is_admin, created_at)
	VALUES($1, $2, $3, $4, $5);`, user.Username, user.Password, user.Email, user.IsAdmin, user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (c *UsersTable) UpdateUser(*User) error {
	return nil
}

func (c *UsersTable) DeleteUser(*User) error {
	return nil
}

func (c *UsersTable) GetDatabase() *sql.DB {
	return c.db
}
