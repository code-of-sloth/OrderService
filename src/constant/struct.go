package constant

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
