package main

import (
	"log"
	"os"

	"github.com/MattBabbage/GoShowcaseAPI/app"
	"github.com/MattBabbage/GoShowcaseAPI/internal/storage"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load(".env")
	store, err := storage.NewPostgressStore(os.Getenv("STORAGE_CONNECTION_STRING"))

	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	if err := store.Init(); err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	server := app.NewAPIServer(":3000", store)
	server.Run()
}
