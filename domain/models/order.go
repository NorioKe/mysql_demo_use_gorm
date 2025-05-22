package models

import (
	"errors"
	"time"
)

type Order struct {
	// gorm.Model
	OrderID   uint64    `gorm:"primaryKey;autoIncrement;column:order_id;comment:订单ID"`
	UserID    uint64    `gorm:"column:user_id;index:idx_user_id;comment:关联用户ID"`
	Amount    float64   `gorm:"column:amount;type:decimal(12,2);not null;comment:订单金额"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;comment:创建时间"`
	IsValid   bool      `gorm:"column:is_valid;type:tinyint(1);default:1;comment:有效性标识(0:无效 1:有效)"`
}

// Invalidate: 订单失效（触发消费总额调整）
func (o *Order) Invalidate() (float64, error) {
	if o.IsValid == false {
		return 0, errors.New("订单已失效")
	}
	o.IsValid = false
	return o.Amount, nil
}
