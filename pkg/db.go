package pkg

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"xorm.io/xorm"
)

const (
	DB_HOST     = "localhost"
	DB_PORT     = 31016
	DB_USER     = "postgres"
	DB_PASSWORD = "abc123"
	DB_NAME     = "exampledb"
)

var Engine *xorm.Engine

func DbConnection() {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	var err error
	Engine, err = xorm.NewEngine("postgres", dbinfo)
	if err != nil {
		log.Fatalf("Failed to create engine: %v", err)
	}
	if err = Engine.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully!")
}
