package model

type CartItem struct {
	ID       string `json:"id" db:"id"`
	CartID   string `json:"cart_id" db:"cart_id"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}
