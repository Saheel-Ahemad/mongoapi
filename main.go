package main

import (
	"fmt"
	"log"
	"mongoapi/router"
	"net/http"
)

func main() {
	fmt.Println("MongoDB API")
	r := router.Router()

	fmt.Println("Server is getting started...")

	// Attempt to start the server and handle any error
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
		return
	}

	// This will only run if the server starts successfully
	fmt.Println("Server started successfully!")
	fmt.Println("Start of server at http://localhost:8080")
	fmt.Println("Listening at port 8080 ...")
}
