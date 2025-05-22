package services

import (
	"errors"
	"fmt"

	// "github.com/NorioKe/mysql_demo_use_gorm/domain/models"
	"github.com/NorioKe/mysql_demo_use_gorm/domain/repositories"
)

// OrderAppService 订单应用服务（事务编排中心）
type OrderAppService struct {
	userRepo  repositories.UserRepository
	orderRepo repositories.OrderRepository
	txManager repositories.TransactionManager // 事务管理器
}

func NewOrderService(ur repositories.UserRepository, or repositories.OrderRepository,
	tm repositories.TransactionManager) *OrderAppService {
	return &OrderAppService{userRepo: ur, orderRepo: or, txManager: tm}
}

// CreateOrderCommand 创建订单命令
type CreateOrderCommand struct {
	UserID uint64
	Amount float64 // 订单金额
}

// CreateOrder 业务流程
func (s *OrderAppService) CreateOrder(cmd CreateOrderCommand) error {
	// 1. 获取用户（不再自动创建）
	user, err := s.userRepo.FindByID(cmd.UserID)
	if errors.Is(err, repositories.ErrorNotFound) {
		return errors.New("用户不存在") // 明确返回业务错误
	} else if err != nil {
		return err
	}

	// 2. 生成订单
	order, err := user.CreateOrder(user.ID, cmd.Amount)
	if err != nil {
		return fmt.Errorf("订单创建失败: %w", err)
	}

	// 3. 金额校验
	if err := user.AddConsumption(cmd.Amount); err != nil {
		return fmt.Errorf("金额校验失败: %w", err)
	}

	// 4. 开启事务
	return s.txManager.Transaction(func() error {
		affect_num, err := s.userRepo.UpdateTotalConsumption(user)
		if affect_num != 1 {
			return errors.New("users表更新行数错误")
		}
		if err != nil {
			return err
		}
		if _, err = s.orderRepo.Save(order); err != nil {
			return err
		}
		return nil
	})
}

// InvalidateOrderCommand 订单失效命令
type InvalidateOrderCommand struct {
	OrderID uint64
}

// InvalidateOrder 订单失效流程
func (s *OrderAppService) InvalidateOrder(cmd InvalidateOrderCommand) error {
	// 获取订单
	order, err := s.orderRepo.FindByID(cmd.OrderID)
	if errors.Is(err, repositories.ErrorInvalid) {
		return errors.New("订单不存在")
	} else if err != nil {
		return err
	}

	// 检查订单是否有效
	if order.IsValid == false {
		return errors.New("订单已失效")
	}

	// 检查是否有对应用户
	user, err := s.userRepo.FindByID(order.UserID)
	if errors.Is(err, repositories.ErrorNotFound) {
		return errors.New(fmt.Sprintf("用户%d不存在", order.UserID))
	} else if err != nil {
		return err
	}

	// 扣除
	err = user.AddConsumption(-order.Amount)
	if err != nil {
		return err
	}

	// 开启事务
	return s.txManager.Transaction(func() error {
		affect_num, err := s.orderRepo.UpdateValidity(order.OrderID, false)
		if err != nil {
			return err
		}
		if affect_num != 1 {
			return errors.New("orders表更新行数错误")
		}

		affect_num, err = s.userRepo.UpdateTotalConsumption(user)
		if err != nil {
			return err
		}
		if affect_num != 1 {
			return errors.New("users表更新行数错误")
		}
		return nil
	})
}
