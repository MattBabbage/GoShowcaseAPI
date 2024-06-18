package main

import (
	"fmt"
	"log"
)

func main() {
	store, err := NewPostgressStore()

	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("%+v\n", store)

	server := NewAPIServer(":3000", store)
	server.Run()
}
