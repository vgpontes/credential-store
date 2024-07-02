package users

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
