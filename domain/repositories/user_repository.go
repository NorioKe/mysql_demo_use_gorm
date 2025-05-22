package repositories

import (
	// "mysql_demo_use_gorm/domain/models"
	"github.com/NorioKe/mysql_demo_use_gorm/domain/models"
)

type UserRepository interface {
	FindByID(id uint64) (*models.User, error)
	FindByEmail(email string) (*models.User, error)         // 查询用户
	Save(user *models.User) (uint64, error)                 // 保存用户信息, 返回用户ID
	UpdateTotalConsumption(user *models.User) (int8, error) // 返回更新的条数
}
