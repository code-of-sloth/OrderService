package orders

import (
	"encoding/json"
	"fmt"
	"net/http"
	"xxx/src/controller/orderOps"
)

type fetchAllOrdersReq struct {
	SortOrder   string `json:"sortOrder"`  //ASC Or DESC Default ASC
	SortKey     int    `json:"sortKey"`    //Default will be id
	FilterKeys  int    `json:"filterKeys"` //Default will be id
	FilterValue string `json:"filterValue"`
}

func FetchAllOrders(w http.ResponseWriter, r *http.Request) {
	var request fetchAllOrdersReq
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orders, err := orderOps.FetchAllOrders(request.SortOrder, request.FilterValue, request.SortKey, request.FilterKeys)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the list of orders as a JSON response
	json.NewEncoder(w).Encode(orders)

}
