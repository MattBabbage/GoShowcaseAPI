package main

import (
	"log"
)

func main() {
	store, err := NewPostgressStore()

	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	if err := store.Init(); err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	server := NewAPIServer(":3000", store)
	server.Run()
}
