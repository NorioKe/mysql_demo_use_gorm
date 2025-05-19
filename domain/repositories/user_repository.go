package repositories

import (
	// "mysql_demo_use_gorm/domain/models"
	"github.com/NorioKe/mysql_demo_use_gorm/domain/models"
)

type UserRepository interface {
	FindByID(id uint64) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Save(user *models.User) error
	Transaction(func(repo UserRepository) error) error // 事务管理接口
}
