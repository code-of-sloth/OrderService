package orders

import (
	"controller/orderOps"
	"encoding/json"
	"fmt"
	"net/http"
)

type addOrderReq struct {
	ItemIds []string `json:"itemIds"`
	Status  string   `Json:"status"`
}

func AddOrderHandler(w http.ResponseWriter, r *http.Request) {
	var request addOrderReq
	// Parse the JSON payload from the request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Store the order details in the database
	id, err := orderOps.StoreOrder(request.ItemIds, request.Status)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the ID of the newly created order
	response := map[string]string{"id": id}
	json.NewEncoder(w).Encode(response)
}
