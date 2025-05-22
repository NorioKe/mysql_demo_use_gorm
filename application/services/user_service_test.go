package services_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/NorioKe/mysql_demo_use_gorm/application/services"

	"github.com/NorioKe/mysql_demo_use_gorm/domain/models"
	"github.com/NorioKe/mysql_demo_use_gorm/domain/repositories"
	"github.com/NorioKe/mysql_demo_use_gorm/interfaces/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service := services.NewUserAppService(mockUserRepo)

	t.Run("成功创建用户", func(t *testing.T) {
		cmd := services.CreateNewUserCommand{
			Name:  "Test User",
			Email: "test@example.com", // 有效邮箱
		}

		// 模拟 FindByEmail 返回无用户
		mockUserRepo.EXPECT().
			FindByEmail(cmd.Email).
			Return(nil, repositories.ErrorNotFound)

		// 这里不验证models.CreateUser是因为:
		// 1. 应用服务测试的核心是验证​​业务流程​​，而非领域模型的创建逻辑。
		// 2. 避免测试耦合​直接依赖 models.CreateUser 的返回结果会使得应用服务测试与领域模型实现​​强耦合​​

		// 模拟 Save 方法
		mockUserRepo.EXPECT().
			Save(gomock.Any()).
			Do(func(user *models.User) {
				assert.Equal(t, cmd.Name, user.Name)
				assert.Equal(t, cmd.Email, user.Email)
				assert.Equal(t, float64(0), user.TotalConsumption)
			}).
			Return(uint64(1001), nil) // 返回模拟的用户ID

		// 执行测试
		userID, err := service.CreateNewUser(cmd)
		assert.NoError(t, err)
		assert.Equal(t, uint64(1001), userID)
	})
}

func TestCreateNewUser_ExistedEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service := services.NewUserAppService(mockUserRepo)

	t.Run("邮箱已存在", func(t *testing.T) {
		cmd := services.CreateNewUserCommand{
			Name:  "Existed User",
			Email: "existed@example.com", // 有效邮箱
		}

		mockUserRepo.EXPECT().FindByEmail(cmd.Email).
			Return(&models.User{ID: 1001}, nil)
		// 执行
		_, err := service.CreateNewUser(cmd)
		assert.ErrorContains(t, err, fmt.Sprintf("邮箱%s已存在", cmd.Email))
	})
}

func TestCreateNewUser_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service := services.NewUserAppService(mockUserRepo)

	t.Run("数据查询失败", func(t *testing.T) {
		cmd := services.CreateNewUserCommand{
			Name:  "Error User",
			Email: "error@example.com", // 有效邮箱
		}

		mockUserRepo.EXPECT().FindByEmail(cmd.Email).
			Return(nil, errors.New("connection error"))
		// 执行
		_, err := service.CreateNewUser(cmd)
		assert.ErrorContains(t, err, "connection error")
	})
}

func TestCreateNewUser_ParamError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service := services.NewUserAppService(mockUserRepo)

	t.Run("邮箱参数错误", func(t *testing.T) {
		cmd := services.CreateNewUserCommand{
			Name:  "Error User",
			Email: "invalid email", // 有效邮箱
		}

		mockUserRepo.EXPECT().FindByEmail(cmd.Email).
			Return(nil, errors.New("邮件格式不正确"))
		// 执行
		_, err := service.CreateNewUser(cmd)
		assert.ErrorContains(t, err, "邮件格式不正确")
	})
}

func TestCreateNewUser_SaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service := services.NewUserAppService(mockUserRepo)

	t.Run("持久化保存错误", func(t *testing.T) {
		cmd := services.CreateNewUserCommand{
			Name:  "User",
			Email: "test@example.com", // 有效邮箱
		}

		mockUserRepo.EXPECT().FindByEmail(cmd.Email).
			Return(nil, nil)

		mockUserRepo.EXPECT().Save(gomock.Any()).
			Return(uint64(0), errors.New("disk full"))

		// 执行
		_, err := service.CreateNewUser(cmd)
		assert.ErrorContains(t, err, "disk full")
	})
}
