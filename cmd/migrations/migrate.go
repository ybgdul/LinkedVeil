package main

import (
	"linkedveil/initializers"
	"linkedveil/models"
	"log"
)

func main() { 
	initializers.LoadEnvs()
	initializers.ConnectDB()

	if err := initializers.DB.AutoMigrate(&models.User{}); err != nil { 
		log.Fatal("Migration failed: ", err)
	}
	
	log.Println("Migration completed succesfully")

}

// go mod migrations/migrate.go for migrations to work out