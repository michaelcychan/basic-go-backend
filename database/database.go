package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/michaelcychan/basic-go-backend/config"
)

func Connect() (*sql.DB, error) {
	dbUser := config.Config("DBUSER")
	dbPassword := config.Config("DBPASSWORD")
	dbURL := config.Config("DBURL")
	dbPort := config.Config("DBPORT")
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/kingdom?sslmode=disable", dbUser, dbPassword, dbURL, dbPort)

	db, dbErr := sql.Open("postgres", connStr)
	if dbErr != nil {
		log.Fatalf("Database error: %s", dbErr)
		return nil, dbErr
	}
	return db, nil
}
