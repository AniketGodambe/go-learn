package main

import (
	"fmt"
	"go-learn/controller"
	"go-learn/router"
	"log"
	"net/http"
)

func main() {

	if err := controller.Initialize(); err != nil {
		log.Fatalf("Failed to initialize todo data: %v", err)
	}

	router := router.InitializeRouter()

	// Start server
	port := ":8080"
	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
