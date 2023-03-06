package orders

import (
	"controller/orderOps"
	"encoding/json"
	"net/http"
)

type updateOrderReq struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func UpdateOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON payload from the request
	var update updateOrderReq
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the status of the corresponding order in the database
	err = orderOps.UpdateOrderStatus(update.ID, update.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success response
	response := map[string]string{"message": "Order status updated successfully"}
	json.NewEncoder(w).Encode(response)
}
