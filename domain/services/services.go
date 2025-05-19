package services

import "github.com/NorioKe/mysql_demo_use_gorm/domain/repositories"

type ConsumptionService struct {
	userRepo  repositories.UserRepository
	orderRepo repositories.OrderRepository
}

func NewConsumptionService(ur repositories.UserRepository, or repositories.OrderRepository) *ConsumptionService {
	return &ConsumptionService{userRepo: ur, orderRepo: or}
}

// AdjustConsumption 领域服务方法示例（后续在应用层调用）
func (s *ConsumptionService) AdjustConsumption(orderID uint64) error {
	// 具体实现将在应用层开发时补充
	return nil
}
