package services_test

import (
	"errors"
	"testing"

	"github.com/NorioKe/mysql_demo_use_gorm/application/services"
	"github.com/NorioKe/mysql_demo_use_gorm/domain/repositories"

	"github.com/NorioKe/mysql_demo_use_gorm/domain/models"
	"github.com/NorioKe/mysql_demo_use_gorm/interfaces/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockTxManager := mocks.NewMockTransactionManager(ctrl)

	service := services.NewOrderService(mockUserRepo, mockOrderRepo, mockTxManager)

	t.Run("成功创建订单", func(t *testing.T) {
		// 初始化用户（消费总额200）
		user := &models.User{ID: 3, TotalConsumption: 200}
		orderAmount := 500.0

		// 预期用户状态更新
		expectedUser := &models.User{
			ID:               3,
			TotalConsumption: 200 + orderAmount, // 700
		}

		// 模拟仓储调用
		mockUserRepo.EXPECT().FindByID(uint64(3)).Return(user, nil)

		// 模拟事务管理器
		mockTxManager.EXPECT().Transaction(gomock.Any()).
			DoAndReturn(func(fn func() error) error {
				// 验证用户保存
				mockUserRepo.EXPECT().UpdateTotalConsumption(expectedUser).Return(int8(1), nil)

				// 验证订单生成
				mockOrderRepo.EXPECT().Save(gomock.Any()).
					Do(func(order *models.Order) {
						assert.Equal(t, user.ID, order.UserID)
						assert.Equal(t, orderAmount, order.Amount)
						assert.True(t, order.IsValid)
					}).Return(uint64(1001), nil)
				return fn()
			})

		// 执行测试
		err := service.CreateOrder(services.CreateOrderCommand{
			UserID: 3,
			Amount: orderAmount,
		})
		assert.NoError(t, err)
	})
}

func TestCreateOrder_UserNotFound(t *testing.T) {
	// 初始化mock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 确保所有EXPECT被验证

	// 创建mock对象
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockTxManager := mocks.NewMockTransactionManager(ctrl)

	// 初始化服务
	service := services.NewOrderService(mockUserRepo, mockOrderRepo, mockTxManager)

	t.Run("用户不存在时报错", func(t *testing.T) {
		mockUserRepo.EXPECT().
			FindByID(uint64(999)).
			Return(nil, repositories.ErrorNotFound).
			Times(1)

		err := service.CreateOrder(services.CreateOrderCommand{
			UserID: 999,
			Amount: 100,
		})
		assert.ErrorContains(t, err, "Not Found")
	})
}

func TestCreateOrder_InvalidAmount(t *testing.T) {
	// 初始化mock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 确保所有EXPECT被验证

	// 创建mock对象
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockTxManager := mocks.NewMockTransactionManager(ctrl)

	// 初始化服务
	service := services.NewOrderService(mockUserRepo, mockOrderRepo, mockTxManager)

	t.Run("订单创建失败", func(t *testing.T) {
		user := &models.User{ID: 1, TotalConsumption: 500}
		mockUserRepo.EXPECT().
			FindByID(uint64(1)).
			Return(user, nil)

		err := service.CreateOrder(services.CreateOrderCommand{
			UserID: 1,
			Amount: -50,
		})
		assert.ErrorContains(t, err, "订单创建失败")
	})
}

func TestCreateOrder_RollbackTx(t *testing.T) {
	// 初始化mock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 确保所有EXPECT被验证

	// 创建mock对象
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockTxManager := mocks.NewMockTransactionManager(ctrl)

	// 初始化服务
	service := services.NewOrderService(mockUserRepo, mockOrderRepo, mockTxManager)

	t.Run("用户保存失败触发回滚", func(t *testing.T) {
		user := &models.User{ID: 2, TotalConsumption: 1000}
		// order := &models.Order{UserID: 2, Amount: 300}

		mockUserRepo.EXPECT().FindByID(uint64(2)).Return(user, nil)

		mockTxManager.EXPECT().Transaction(gomock.Any()).
			DoAndReturn(func(fn func() error) error {
				mockUserRepo.EXPECT().UpdateTotalConsumption(user).Return(int8(1), errors.New("db error"))
				return fn()
			})

		err := service.CreateOrder(services.CreateOrderCommand{UserID: 2, Amount: 300})
		assert.ErrorContains(t, err, "db error")
	})
}

func TestCreateOrder_ErrorUpdateUsersTable(t *testing.T) {
	// 初始化mock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 确保所有EXPECT被验证

	// 创建mock对象
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockTxManager := mocks.NewMockTransactionManager(ctrl)

	// 初始化服务
	service := services.NewOrderService(mockUserRepo, mockOrderRepo, mockTxManager)

	t.Run("users表更新行数错误", func(t *testing.T) {
		user := &models.User{ID: 2, TotalConsumption: 1000}
		// order := &models.Order{UserID: 2, Amount: 300}

		mockUserRepo.EXPECT().FindByID(uint64(2)).Return(user, nil)

		mockTxManager.EXPECT().Transaction(gomock.Any()).
			DoAndReturn(func(fn func() error) error {
				mockUserRepo.EXPECT().UpdateTotalConsumption(user).Return(int8(2), nil)
				return fn()
			})

		err := service.CreateOrder(services.CreateOrderCommand{UserID: 2, Amount: 300})
		assert.ErrorContains(t, err, "users表更新行数错误")
	})
}

func TestCreateOrder_RollbackTxOrder(t *testing.T) {
	// 初始化mock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 确保所有EXPECT被验证

	// 创建mock对象
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockTxManager := mocks.NewMockTransactionManager(ctrl)

	// 初始化服务
	service := services.NewOrderService(mockUserRepo, mockOrderRepo, mockTxManager)

	t.Run("orders表插入错误导致回滚", func(t *testing.T) {
		user := &models.User{ID: 2, TotalConsumption: 1000}
		// order := &models.Order{UserID: 2, Amount: 300}

		mockUserRepo.EXPECT().FindByID(uint64(2)).Return(user, nil)

		mockTxManager.EXPECT().Transaction(gomock.Any()).
			DoAndReturn(func(fn func() error) error {
				mockUserRepo.EXPECT().UpdateTotalConsumption(user).Return(int8(1), nil)
				mockOrderRepo.EXPECT().Save(gomock.Any()).Return(uint64(1001), errors.New("db error"))
				return fn()
			})

		err := service.CreateOrder(services.CreateOrderCommand{UserID: 2, Amount: 300})
		assert.ErrorContains(t, err, "db error")
	})
}

// 订单失效测试
func TestInvalidAmount_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 初始化mock对象
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTxManager := mocks.NewMockTransactionManager(ctrl)

	// 创建服务实例
	service := services.NewOrderService(mockUserRepo, mockOrderRepo, mockTxManager)

	t.Run("成功失效订单并扣减消费", func(t *testing.T) {
		// 设置订单预期
		mockOrderRepo.EXPECT().
			FindByID(uint64(1001)).
			Return(&models.Order{
				OrderID: 1001,
				UserID:  2001,
				Amount:  500.0,
				IsValid: true,
			}, nil)

		// 设置用户预期
		mockUserRepo.EXPECT().
			FindByID(uint64(2001)).
			Return(&models.User{
				ID:               2001,
				TotalConsumption: 1500.0,
			}, nil)

		// 事务管理器预期
		mockTxManager.EXPECT().
			Transaction(gomock.Any()).
			DoAndReturn(func(fn func() error) error {
				// 验证事务内操作
				mockOrderRepo.EXPECT().
					UpdateValidity(uint64(1001), false).
					Return(int8(1), nil)

				mockUserRepo.EXPECT().
					UpdateTotalConsumption(&models.User{
						ID:               2001,
						TotalConsumption: 1000.0, // 500元扣减
					}).
					Return(int8(1), nil)

				return fn()
			})

		err := service.InvalidateOrder(services.InvalidateOrderCommand{OrderID: 1001})
		assert.NoError(t, err)
	})
}

func TestInvalidAmount_OrderNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 初始化mock对象
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTxManager := mocks.NewMockTransactionManager(ctrl)

	// 创建服务实例
	service := services.NewOrderService(mockUserRepo, mockOrderRepo, mockTxManager)

	t.Run("订单不存在时报错", func(t *testing.T) {
		mockOrderRepo.EXPECT().
			FindByID(uint64(999)).
			Return(nil, repositories.ErrorInvalid)

		err := service.InvalidateOrder(services.InvalidateOrderCommand{OrderID: 999})
		assert.ErrorContains(t, err, "Invalid")
	})
}

func TestInvalidAmount_OrderInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 初始化mock对象
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTxManager := mocks.NewMockTransactionManager(ctrl)

	// 创建服务实例
	service := services.NewOrderService(mockUserRepo, mockOrderRepo, mockTxManager)

	t.Run("重复失效订单时报错", func(t *testing.T) {
		mockOrderRepo.EXPECT().
			FindByID(uint64(1002)).
			Return(&models.Order{
				OrderID: 1002,
				IsValid: false,
			}, nil)

		err := service.InvalidateOrder(services.InvalidateOrderCommand{OrderID: 1002})
		assert.ErrorContains(t, err, "订单已失效")
	})
}

func TestInvalidAmount_UserNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 初始化mock对象
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTxManager := mocks.NewMockTransactionManager(ctrl)

	// 创建服务实例
	service := services.NewOrderService(mockUserRepo, mockOrderRepo, mockTxManager)

	t.Run("订单关联用户不存在时报错", func(t *testing.T) {
		mockOrderRepo.EXPECT().
			FindByID(uint64(1003)).
			Return(&models.Order{
				OrderID: 1003,
				UserID:  3001,
				IsValid: true,
			}, nil)

		mockUserRepo.EXPECT().
			FindByID(uint64(3001)).
			Return(nil, repositories.ErrorNotFound)

		err := service.InvalidateOrder(services.InvalidateOrderCommand{OrderID: 1003})
		assert.ErrorContains(t, err, "Not Found")
	})
}

func TestInvalidAmount_AmountDesFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 初始化mock对象
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTxManager := mocks.NewMockTransactionManager(ctrl)

	// 创建服务实例
	service := services.NewOrderService(mockUserRepo, mockOrderRepo, mockTxManager)

	t.Run("用户余额不足时报错", func(t *testing.T) {
		mockOrderRepo.EXPECT().
			FindByID(uint64(1004)).
			Return(&models.Order{
				OrderID: 1004,
				UserID:  2002,
				Amount:  1000.0,
				IsValid: true,
			}, nil)

		mockUserRepo.EXPECT().
			FindByID(uint64(2002)).
			Return(&models.User{
				ID:               2002,
				TotalConsumption: 500.0,
			}, nil)

		err := service.InvalidateOrder(services.InvalidateOrderCommand{OrderID: 1004})
		assert.ErrorContains(t, err, "消费总额不足")
	})
}

func TestInvalidAmount_TxRollBack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 初始化mock对象
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTxManager := mocks.NewMockTransactionManager(ctrl)

	// 创建服务实例
	service := services.NewOrderService(mockUserRepo, mockOrderRepo, mockTxManager)

	t.Run("事务内操作失败时回滚", func(t *testing.T) {
		mockOrderRepo.EXPECT().
			FindByID(uint64(1005)).
			Return(&models.Order{
				OrderID: 1005,
				UserID:  2003,
				Amount:  200.0,
				IsValid: true,
			}, nil)

		mockUserRepo.EXPECT().
			FindByID(uint64(2003)).
			Return(&models.User{
				ID:               2003,
				TotalConsumption: 1000.0,
			}, nil)

		mockTxManager.EXPECT().
			Transaction(gomock.Any()).
			DoAndReturn(func(fn func() error) error {
				mockOrderRepo.EXPECT().
					UpdateValidity(uint64(1005), false).
					Return(int8(0), errors.New("数据库连接失败"))

				// 用户保存不会被调用
				return fn()
			})

		err := service.InvalidateOrder(services.InvalidateOrderCommand{OrderID: 1005})
		assert.ErrorContains(t, err, "数据库连接失败")
	})
}

func TestInvalidAmount_ErrorUpdateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 初始化mock对象
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTxManager := mocks.NewMockTransactionManager(ctrl)

	// 创建服务实例
	service := services.NewOrderService(mockUserRepo, mockOrderRepo, mockTxManager)

	t.Run("orders表更新行数错误", func(t *testing.T) {
		mockOrderRepo.EXPECT().
			FindByID(uint64(1005)).
			Return(&models.Order{
				OrderID: 1005,
				UserID:  2003,
				Amount:  200.0,
				IsValid: true,
			}, nil)

		mockUserRepo.EXPECT().
			FindByID(uint64(2003)).
			Return(&models.User{
				ID:               2003,
				TotalConsumption: 1000.0,
			}, nil)

		mockTxManager.EXPECT().
			Transaction(gomock.Any()).
			DoAndReturn(func(fn func() error) error {
				mockOrderRepo.EXPECT().
					UpdateValidity(uint64(1005), false).
					Return(int8(2), nil)

				// 用户保存不会被调用
				return fn()
			})

		err := service.InvalidateOrder(services.InvalidateOrderCommand{OrderID: 1005})
		assert.ErrorContains(t, err, "orders表更新行数错误")
	})
}

func TestInvalidAmount_ErrorUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 初始化mock对象
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTxManager := mocks.NewMockTransactionManager(ctrl)

	// 创建服务实例
	service := services.NewOrderService(mockUserRepo, mockOrderRepo, mockTxManager)

	t.Run("users表更新行数错误", func(t *testing.T) {
		mockOrderRepo.EXPECT().
			FindByID(uint64(1005)).
			Return(&models.Order{
				OrderID: 1005,
				UserID:  2003,
				Amount:  200.0,
				IsValid: true,
			}, nil)

		mockUserRepo.EXPECT().
			FindByID(uint64(2003)).
			Return(&models.User{
				ID:               2003,
				TotalConsumption: 1000.0,
			}, nil)

		mockTxManager.EXPECT().
			Transaction(gomock.Any()).
			DoAndReturn(func(fn func() error) error {
				mockOrderRepo.EXPECT().
					UpdateValidity(uint64(1005), false).
					Return(int8(1), nil)
				mockUserRepo.EXPECT().
					UpdateTotalConsumption(gomock.Any()).
					Return(int8(2), nil)

				return fn()
			})

		err := service.InvalidateOrder(services.InvalidateOrderCommand{OrderID: 1005})
		assert.ErrorContains(t, err, "users表更新行数错误")
	})
}
