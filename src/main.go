package main

import (
	"log"
	"os"

	"github.com/MattBabbage/GoShowcaseAPI/internal/api"
	"github.com/MattBabbage/GoShowcaseAPI/internal/storage"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load("../.env")
	store, err := storage.NewPostgressStore(os.Getenv("STORAGE_CONNECTION_STRING"))

	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	if err := store.Init(); err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	server := api.NewAPIServer(":8080", store)
	server.Run()
}
