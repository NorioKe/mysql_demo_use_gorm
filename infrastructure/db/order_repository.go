package db

import (
	"errors"

	"github.com/NorioKe/mysql_demo_use_gorm/domain/models"

	"github.com/NorioKe/mysql_demo_use_gorm/domain/repositories"
	"gorm.io/gorm"
)

type GormOrderRepository struct {
	db *gorm.DB
}

func NewGormOrderRepository(db *gorm.DB) repositories.OrderRepository {
	return &GormOrderRepository{db: db}
}

func (r *GormOrderRepository) FindByID(orderID uint64) (*models.Order, error) {
	var order models.Order
	if err := r.db.First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repositories.ErrorNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *GormOrderRepository) Save(order *models.Order) (uint64, error) {
	if err := r.db.Save(order).Error; err != nil {
		return uint64(0), err
	}
	return order.OrderID, nil
}

func (r *GormOrderRepository) UpdateValidity(orderID uint64, isValid bool) (int8, error) {
	var valid_num int8
	if isValid {
		valid_num = 1
	} else {
		valid_num = 0
	}
	result := r.db.Model(&models.Order{}).
		Where("order_id = ?", orderID).
		Update("is_valid", valid_num)
	if result.Error != nil {
		return int8(0), result.Error
	}
	return int8(result.RowsAffected), nil
}
