package main

import (
	"database/sql"

	//We are using the pgx driver to connect to PostgreSQL
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database interface {
	CreateUser(*User) error
	GetUserByID(int) (string, error)
	GetUserByUsername(string) (string, error)
	GetUsers() ([]*GetUsersResponse, error)
	UpdateUser(*User) error
	DeleteUser(*User) error
}

type PostgresDB struct {
	db *sql.DB
}

func ConnectDB() (*PostgresDB, error) {
	connStr := "postgres://postgres:IdQ2o.vm=,WNOQ_7MsY,o3VCZseAoI@credentialstoredb.cr22sw42g2wm.us-east-1.rds.amazonaws.com:5432"
	//Pass the driver name and the connection string
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}
	//defer db.Close()
	//Call db.Ping() to check the connection
	if pingErr := db.Ping(); pingErr != nil {
		return nil, pingErr
	}

	return &PostgresDB{
		db: db,
	}, nil
}

/*func (s *PostgresDB) Init() error {
	return s.createUserTable()
}

func (s *PostgresDB) createUserTable() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS users (
		user_id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(50),
		salt VARCHAR(50),
		is_admin BOOLEAN,
		created_at TIMESTAMP NOT NULL
	);`)
	return err
} */

func (s *PostgresDB) CreateUser(user *User) error {
	_, err := s.db.Exec(`
	INSERT INTO users(username, password, email, salt, is_admin, created_at)
	VALUES($1, $2, $3, $4, $5, $6);`, user.Username, user.Password, user.Email, user.Salt, user.IsAdmin, user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresDB) GetUserByID(id int) (string, error) {
	row := s.db.QueryRow(`
	SELECT username
	FROM users
	WHERE user_id=$1;`, id)
	var username string
	err := row.Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func (s *PostgresDB) GetUserByUsername(userName string) (string, error) {
	row := s.db.QueryRow(`
	SELECT username
	FROM users
	WHERE username$1;`, userName)
	var username string
	err := row.Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func (s *PostgresDB) GetUsers() ([]*GetUsersResponse, error) {
	rows, err := s.db.Query("SELECT * FROM users")
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

func (s *PostgresDB) UpdateUser(*User) error {
	return nil
}

func (s *PostgresDB) DeleteUser(*User) error {
	return nil
}
