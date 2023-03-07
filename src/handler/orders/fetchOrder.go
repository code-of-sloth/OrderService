package orders

import (
	"encoding/json"
	"net/http"
	"xxx/src/controller/orderOps"
)

func FetchOrdersHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the query parameters from the request
	queryParams := r.URL.Query()
	id := queryParams.Get("id")

	// Fetch the orders from the database based on the query parameters
	orders, err := orderOps.FetchAllOrders("", id, -1, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the list of orders as a JSON response
	json.NewEncoder(w).Encode(orders)
}
