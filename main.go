package main

import (
	"fmt"
	"go-learn/router"
	"log"
	"net/http"
)

func main() {

	router := router.InitializeRouter()

	// Start server
	port := ":8080"
	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
