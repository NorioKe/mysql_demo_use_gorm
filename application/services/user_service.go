package services

import (
	"errors"
	"fmt"

	"github.com/NorioKe/mysql_demo_use_gorm/domain/models"
	"github.com/NorioKe/mysql_demo_use_gorm/domain/repositories"
)

// UserAppService 用户应用服务（事务编排中心）
type UserAppService struct {
	userRepo repositories.UserRepository
}

func NewUserAppService(ur repositories.UserRepository) *UserAppService {
	return &UserAppService{
		userRepo: ur,
	}
}

// 新增用户的命令
type CreateNewUserCommand struct {
	Name  string
	Email string // 邮箱为unique
}

// CreateNewUser: 新增用户
func (u *UserAppService) CreateNewUser(cmd CreateNewUserCommand) (uint64, error) {
	// 1. 检查邮箱是否存在
	exist_user, err := u.userRepo.FindByEmail(cmd.Email)
	if err != nil && !errors.Is(err, repositories.ErrorNotFound) {
		return 0, fmt.Errorf("DB error: %w", err)
	}
	if exist_user != nil {
		return exist_user.ID, errors.New(fmt.Sprintf("邮箱%s已存在", cmd.Email))
	}

	// 2. 创建用户
	new_user, err := models.CreateUser(cmd.Name, cmd.Email)
	if err != nil {
		return 0, fmt.Errorf("创建用户失败: %w", err)
	}

	// 3. 存储用户
	userid, err := u.userRepo.Save(new_user)
	if err != nil {
		return 0, fmt.Errorf("DB error: %w", err)
	}
	return userid, nil
}
