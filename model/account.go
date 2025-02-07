package model

import "time"

type Account struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

type Transaction struct {
	ID        int       `json:"id"`
	FromID    int       `json:"from_id"`
	ToID      int       `json:"to_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
