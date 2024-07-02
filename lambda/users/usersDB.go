package main

import (
	"database/sql"

	//We are using the pgx driver to connect to PostgreSQL
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database interface {
	CreateUser(*User) error
	GetUserByID(int) (*User, error)
	GetUserByUsername(string) (*User, error)
	UpdateUser(*User) error
	DeleteUser(*User) error
}

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB() (*PostgresDB, error) {
	connStr := "postgres://postgres:geogGkshTw^Tw5BU5bT1rP_iTOWqpn@credentialstoredb.cr22sw42g2wm.us-east-1.rds.amazonaws.com:5432"
	//Pass the driver name and the connection string
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	//Call db.Ping() to check the connection
	if pingErr := db.Ping(); pingErr != nil {
		return nil, pingErr
	}

	return &PostgresDB{
		db: db,
	}, nil
}

func (s *PostgresDB) Init() error {
	return s.createUserTable()
}

func (s *PostgresDB) createUserTable() error {
	sqlStatement := `
		CREATE TABLE IF NOT EXISTS users(
		user_id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(50) NOT NULL,
		salt VARCHAR(50) NOT NULL,
		is_admin BOOLEAN NOT NULL,
		created_at TIMESTAMP NOT NULL)`
	
	_, err := s.db.Exec(sqlStatement)
	return err
}

func (s *PostgresDB) CreateUser(*User) error {
	return nil
}

func (s *PostgresDB) GetUserByID(id int) (*User, error) {
	return nil, nil
}

func (s *PostgresDB) GetUserByUsername(userName string) (*User, error) {
	return nil, nil
}
func (s *PostgresDB) UpdateUser(*User) error {
	return nil
}

func (s *PostgresDB) DeleteUser(*User) error {
	return nil
}
