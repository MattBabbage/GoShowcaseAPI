package main

import "fmt"

func main() {
	fmt.Println("HelloWorld")
	server := NewAPIServer(":3000")
	server.Run()
}
