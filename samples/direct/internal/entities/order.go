package entities

import "time"

type Order struct {
	ID        int64      `json:"id"`
	Amount    float64    `json:"amount"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func NewOrder(id int64, amount float64) *Order {
	return &Order{
		ID:        id,
		Amount:    amount,
		CreatedAt: time.Now(),
	}
}

func (o *Order) Update(amount float64) {
	o.Amount = amount
	o.UpdatedAt = &time.Time{}
}
