package main

import (
	"fmt"
	"log"

	"github.com/NorioKe/mysql_demo_use_gorm/application/services"

	"github.com/NorioKe/mysql_demo_use_gorm/infrastructure/db"

	"github.com/NorioKe/mysql_demo_use_gorm/infrastructure/config"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	gorm_DB, err := db.NewDB(&cfg.Database)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 初始化仓储（repository）
	user_repo := db.NewGormUserRepository(gorm_DB)
	order_repo := db.NewGormOrderRepository(gorm_DB)
	tx_repo := db.NewTransactionManager(gorm_DB)

	// 初始化应用服务
	user_service := services.NewUserAppService(user_repo)
	order_service := services.NewOrderService(user_repo, order_repo, tx_repo)

	// 示例1: 创建用户
	user_id, err := user_service.CreateNewUser(services.CreateNewUserCommand{
		Name:  "Xiao Hong",
		Email: "xiaohong@163.com",
	})
	if err != nil {
		log.Fatalf("创建用户失败: %v", err)
	}
	fmt.Printf("User ID: %d\n", user_id)

	// 创建订单
	err = order_service.CreateOrder(services.CreateOrderCommand{
		UserID: uint64(1),
		Amount: float64(1000),
	})
	if err != nil {
		log.Fatalf("创建订单失败: %v", err)
	}
}
