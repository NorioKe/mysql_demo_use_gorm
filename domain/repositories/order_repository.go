package repositories

import (
	// "mysql_demo_use_gorm/domain/models"
	"github.com/NorioKe/mysql_demo_use_gorm/domain/models"
)

// OrderRepository 订单实体的数据访问契约
type OrderRepository interface {
	FindByID(orderID uint64) (*models.Order, error)
	Save(order *models.Order) (uint64, error)                  // 返回订单ID
	UpdateValidity(orderID uint64, isValid bool) (int8, error) // 返回影响的行数
}
