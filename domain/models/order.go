package models

import (
	"errors"
	"time"
)

type Order struct {
	OrderID   uint64
	UserID    uint64
	Amount    float64
	CreatedAt time.Time
	IsValid   bool
}

// Invalidate: 订单失效（触发消费总额调整）
func (o *Order) Invalidate() (float64, error) {
	if o.IsValid == false {
		return 0, errors.New("订单已失效")
	}
	o.IsValid = false
	return o.Amount, nil
}
