package initialisers

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb() {
	dsn := os.Getenv("DB")

	var err error
	//dsn opens a connection to the PostgreSQL database using the provided DSN (Data Source Name) and the GORM library. It assigns the resulting database connection to the global variable DB. If there is an error during the connection, it logs a fatal error message and terminates the program. If the connection is successful, it prints a confirmation message to the console.

	//It assigns the resulting database connection to the global variable DB
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	fmt.Println("✅ Connected to PostgreSQL")
}
