package main

import (
	"log"

	"github.com/ajiteshreddy7/yc-go-scraper/internal/db"
)

func main() {
	// Connect to database
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Check if admin user already exists
	_, _, _, err = database.GetUserByUsername("admin")
	if err == nil {
		log.Println("Admin user already exists, skipping initialization")
		return
	}

	// Create admin user
	err = database.CreateUser("admin", "password123")
	if err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	log.Println("✅ Admin user created successfully!")
	log.Println("📋 Login credentials:")
	log.Println("   Username: admin")
	log.Println("   Password: password123")
	log.Println("🔒 Please change the password after first login!")
}
