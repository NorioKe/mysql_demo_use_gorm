package repositories

import (
	"mysql_demo_use_gorm/domain/models"
	// "github.com/your-project/domain/models"
)

// OrderRepository 订单实体的数据访问契约
type OrderRepository interface {
	FindByID(orderID uint64) (*models.Order, error)
	Save(order *models.Order) error
	UpdateValidity(orderID uint64, isValid bool) error
}
