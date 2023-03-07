package orderOps

import (
	"errors"
	"fmt"
	"strings"
	"xxx/src/constant"
	"xxx/src/services/sqlDb"

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
		return "", err
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

func FetchAllOrders(sortOrder, filterValue string, sortKey, filterKey int) (result []constant.Order, err error) {
	// Implement the code for fetching the orders from the database based on the query parameters
	var sortQuery, filterQuery string
	if sortKey != -1 {
		if sortOrder == "" {
			sortOrder = "ASC"
		}
		sortQuery = fmt.Sprintf(" ORDER BY %v %s", constant.OrderKeyValue[sortKey], sortOrder)
	}
	if filterKey != -1 && filterValue != "" {
		filterQuery = fmt.Sprintf("WHERE %v LIKE '%s%s%s' ", constant.OrderKeyValue[filterKey], "%", filterValue, "%")
	}

	query := make([]string, 1)
	query[0] = fmt.Sprintf("select id,status,currency_unit from orders %s%s", filterQuery, sortQuery)
	resp, err := sqlDb.RunQuery(query...)
	if err != nil {
		return
	} else if msg, ok := resp[0][0]["error"]; ok {
		return result, errors.New(msg.(string))
	}

	length := len(resp[0])
	if length == 0 {
		return
	}

	result, query = make([]constant.Order, length), make([]string, length)

	for i, v := range resp[0] {
		query[i] = fmt.Sprintf("SELECT id,description,price,qty from items WHERE order_ids LIKE '%s%v%s'", `%`, int(v["id"].(float64)), `%`)
	}

	resp, err = sqlDb.RunQuery(query...)
	if err != nil {
		return
	} else if msg, ok := resp[0][0]["error"]; ok {
		return result, errors.New(msg.(string))
	}

	for i, order := range resp {
		var item constant.Item
		for _, v := range order {
			item.Description = v["description"].(string)
			item.ID = fmt.Sprint(v["id"])
			item.Price = v["price"].(float64)
			item.Quantity = int(v["qty"].(float64))
			result[i].Total += item.Price
			result[i].Items = append(result[i].Items, item)
		}
		result[i].CurrencyUnit = resp[0][0]["currency_unit"].(string)
		result[i].Status = resp[0][0]["status"].(string)
		result[i].ID = fmt.Sprint(resp[0][0]["id"])
	}

	return

}
