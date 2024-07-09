package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	//We are using the pgx driver to connect to PostgreSQL
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	_ "github.com/lib/pq"
)

type Database interface {
	CreateUser(*User) error
	GetUserByUsername(string) (*GetUsersResponse, error)
	GetUsers() ([]*GetUsersResponse, error)
	UpdateUser(*User) error
	DeleteUser(*User) error
}

type PostgresDB struct {
	db *sql.DB
}

func ConnectDB() (*PostgresDB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	region := os.Getenv("AWS_REGION")

	var dbEndpoint string = fmt.Sprintf("%s:%s", host, port)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Println("Error loading default AWS SDK config")
		return nil, err
	}

	authenticationToken, err := auth.BuildAuthToken(
		context.TODO(), dbEndpoint, region, user, cfg.Credentials)
	if err != nil {
		log.Printf("Error building auth token to database. Endpoint: %s, Region: %s, User: %s", dbEndpoint, region, user)
		return nil, err
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, authenticationToken, dbName)
	//Pass the driver name and the connection string
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error connecting to database. Host: %s, Port: %s, User: %s, dbName: %s", host, port, user, dbName)
		return nil, err
	}
	//Call db.Ping() to check the connection
	if pingErr := db.Ping(); pingErr != nil {
		log.Printf("Error pinging database")
		return nil, pingErr
	}

	return &PostgresDB{
		db: db,
	}, nil
}

func (s *PostgresDB) CreateUser(user *User) error {
	_, err := s.db.Exec(`
	INSERT INTO users(username, password, email, salt, is_admin, created_at)
	VALUES($1, $2, $3, $4, $5, $6);`, user.Username, user.Password, user.Email, user.Salt, user.IsAdmin, user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresDB) GetUserByUsername(userName string) (*GetUsersResponse, error) {
	row := s.db.QueryRow("SELECT username, email FROM users WHERE username=$1;", userName)
	user := &GetUsersResponse{}
	err := row.Scan(&user.Username, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *PostgresDB) GetUsers() ([]*GetUsersResponse, error) {
	rows, err := s.db.Query("SELECT username, email FROM users;")
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
