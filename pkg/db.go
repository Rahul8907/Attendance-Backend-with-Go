package pkg

import (
	"log"

	_ "github.com/lib/pq"
	"xorm.io/xorm"
)

var engine *xorm.Engine

func DbConnection() {
	var err error
	engine, err = xorm.NewEngine("postgres", "user=postgres password=abc123 dbname=exampledb sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to create engine: %v", err)
	}

	// Optional: Test the connection
	if err = engine.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully!")
}
