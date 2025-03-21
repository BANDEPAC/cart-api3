package models

type Cart struct {
	ID    string     `json:"id" db:"id"`
	Items []CartItem `json:"items" `
}
