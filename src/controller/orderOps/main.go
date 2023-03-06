package orderOps

import (
	"constant"
	"errors"
	"fmt"
	"services/sqlDb"
	"strings"

	"github.com/dchest/uniuri"
)

type Order struct {
	ID           string  `json:"id"`
	Status       string  `json:"status"`
	Items        []Item  `json:"items"`
	Total        float64 `json:"total"`
	CurrencyUnit string  `json:"currencyUnit"`
}

type Item struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

func StoreOrder(itemId []string, status string) (string, error) {
	// Implement the code for storing the order details in the database
	if len(itemId) == 0 {
		return "", errors.New("No Item is Selected")
	}

	orderId := uniuri.NewLenChars(11, []byte("1234567890"))
	query := make([]string, 2)
	if status != "" {
		query[0] = fmt.Sprintf("INSERT into orders (id,status) values (%s,'%s')", orderId, status)
	} else {
		query[0] = fmt.Sprintf("INSERT into orders (id) values (%s)", orderId)
	}
	query[1] = fmt.Sprintf("UPDATE items SET order_ids=CONCAT(order_ids,'-%s-') WHERE id in (%s)", strings.Join(itemId, ","))

	resp, err := sqlDb.RunQuery(query...)
	if err != nil {
		return err
	}

	for _, v := range resp {
		if msg, ok := v[0]["error"]; ok {
			return "", errors.New(msg.(string))
		}
	}

	return orderId, nil
}

func UpdateOrderStatus(id string, status string) error {
	// Implement the code for updating the status of the order with the given ID in the database
	query := fmt.Sprintf("UPDATE orders SET status='%s' WHERE id=%s", status, id)
	resp, err := sqlDb.RunQuery(query)
	if err != nil {
		return err
	} else if msg, ok := resp[0][0]["error"]; ok {
		return errors.New(msg.(string))
	}

	return nil
}

func FetchOrders(id string) (result constant.Order, err error) {
	// Implement the code for fetching the orders from the database based on the query parameters
	if id == "" {
		return result, errors.New("Order Id is empty")
	}

	query := make([]string, 2)
	query[0] = fmt.Sprintf("SELECT id,status,currency_unit from orders WHERE id=%s", id)
	query[1] = fmt.Sprintf("SELECT id,description,price,qty from items WHERE order_ids LIKE '%s%s%s'", `%`, id, `%`)
	resp, err := sqlDb.RunQuery(query...)
	if err != nil {
		return
	} else if len(resp[0]) == 0 || len(resp[1]) == 0 {
		return result, errors.New("Order and/or items empty")
	}
	for _, v := range resp {
		if msg, ok := v[0]["error"]; ok {
			return result, errors.New(msg.(string))
		}
	}

	var item constant.Item
	for _, v := range resp[1] {
		item.Description = v["description"].(string)
		item.ID = fmt.Sprint(v["id"])
		item.Price = v["price"].(float64)
		item.Quantity = int(v["qty"].(float64))
		result.Total += item.Price
		result.Items = append(result.Items, item)
	}

	result.CurrencyUnit = resp[0][0]["currency_unit"].(string)
	result.Status = resp[0][0]["status"].(string)
	result.ID = fmt.Sprint(resp[0][0]["status"])

	return
}
