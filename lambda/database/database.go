package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	_ "github.com/lib/pq"
)

type CredentialStoreDB struct {
	db *sql.DB
}

func NewCredentialStoreDB() (*CredentialStoreDB, error) {
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

	return &CredentialStoreDB{
		db: db,
	}, nil
}

func (c *CredentialStoreDB) GetDatabase() *sql.DB {
	return c.db
}
