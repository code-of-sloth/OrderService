package main

import (
	"log"
	"net/http"

	"handler/orders"

	"github.com/gorilla/mux"
)

// Define the main function to start the HTTP server
func main() {

	// Create the Gorilla Mux router and register the API endpoints
	r := mux.NewRouter()
	r.HandleFunc("/orders/add", orders.AddOrderHandler).Methods("POST")
	r.HandleFunc("/orders/updateStatus", orders.UpdateOrderHandler).Methods("POST")
	r.HandleFunc("/orders/{id}", orders.FetchOrdersHandler).Methods("GET")
	r.HandleFunc("/orders", orders.FetchAllOrders).Methods("POST")

	// Start the HTTP server
	log.Println("Starting HTTP server on port 8080...")
	http.ListenAndServe(":8080", r)
}

// Implement the initializeDB function to initialize the database connection
