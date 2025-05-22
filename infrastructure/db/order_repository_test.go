package db_test

import (
	"testing"

	"github.com/NorioKe/mysql_demo_use_gorm/domain/models"
	"github.com/NorioKe/mysql_demo_use_gorm/domain/repositories"
	"github.com/NorioKe/mysql_demo_use_gorm/infrastructure/config"
	"github.com/NorioKe/mysql_demo_use_gorm/infrastructure/db"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestOrderDB(t *testing.T) *gorm.DB {
	// 使用测试的数据库
	cfg := &config.DatabaseConfig{
		Host:      "localhost",
		Port:      3306,
		User:      "gouser",
		Password:  "StrongPass123!",
		DBName:    "go_dev_test",
		Charset:   "utf8mb4",
		ParseTime: true,
	}

	dbConn, err := db.NewDB(cfg)
	assert.NoError(t, err, "数据库连接失败")

	// 迁移表结构
	// 替换AutoMigrate为同步表结构(不要用AutoMigrate避免表结构被改变)
	if !dbConn.Migrator().HasTable(&models.Order{}) {
		if err := dbConn.Migrator().CreateTable(&models.Order{}); err != nil {
			t.Fatal(err)
		}
	} else {
		// 已有表时检查字段差异
		if err := dbConn.AutoMigrate(&models.Order{}); err != nil {
			t.Fatal(err)
		}
	}

	// 清空表时使用Delete替代TRUNCATE
	// 清空环境
	if err := dbConn.Exec("DELETE FROM orders").Error; err != nil {
		t.Fatal(err)
	}
	if err := dbConn.Exec("DELETE FROM users").Error; err != nil {
		t.Fatal(err)
	}

	return dbConn
}

func TestOrderRepository_FindByID(t *testing.T) {
	// 连接db
	dbConn := setupTestOrderDB(t)
	repo := db.NewGormOrderRepository(dbConn)
	user_repo := db.NewGormUserRepository(dbConn)

	t.Run("成功查找到订单ID", func(t *testing.T) {
		// 由于外键约束，先把用户存进去
		user := &models.User{
			ID:               uint64(10001),
			Name:             "test",
			Email:            "test@example.com",
			TotalConsumption: 2000,
		}
		_, err := user_repo.Save(user)
		assert.NoError(t, err)

		order := &models.Order{
			OrderID: uint64(2025052110001),
			UserID:  uint64(10001),
			Amount:  float64(1000),
			IsValid: true,
		}

		// 执行保存
		_, err = repo.Save(order)
		assert.NoError(t, err)

		// 查找
		found_order, err := repo.FindByID(order.OrderID)
		// 验证数据
		assert.NoError(t, err)
		assert.Equal(t, uint64(2025052110001), found_order.OrderID)
		assert.Equal(t, uint64(10001), found_order.UserID)
		assert.Equal(t, float64(1000), found_order.Amount)
		assert.Equal(t, true, found_order.IsValid)
	})

	t.Run("无法找到用户ID", func(t *testing.T) {
		// 查找
		_, err := repo.FindByID(uint64(1002))
		assert.ErrorIs(t, err, repositories.ErrorNotFound)
	})

	// 清空环境
	if err := dbConn.Exec("DELETE FROM orders").Error; err != nil {
		t.Fatal(err)
	}
	if err := dbConn.Exec("DELETE FROM users").Error; err != nil {
		t.Fatal(err)
	}
}

func TestOrderRepository_Save(t *testing.T) {
	// 连接db
	dbConn := setupTestOrderDB(t)
	repo := db.NewGormOrderRepository(dbConn)
	user_repo := db.NewGormUserRepository(dbConn)

	t.Run("成功保存用户", func(t *testing.T) {
		// 由于外键约束，先把用户存进去
		user := &models.User{
			ID:               uint64(10001),
			Name:             "test",
			Email:            "test@example.com",
			TotalConsumption: 2000,
		}
		_, err := user_repo.Save(user)
		assert.NoError(t, err)

		order := &models.Order{
			OrderID: uint64(2025052110001),
			UserID:  uint64(10001),
			Amount:  float64(1000),
			IsValid: true,
		}

		// 执行保存
		_, err = repo.Save(order)
		assert.NoError(t, err)

		// 查找
		found_order, err := repo.FindByID(order.OrderID)
		// 验证数据
		assert.NoError(t, err)
		assert.Equal(t, uint64(2025052110001), found_order.OrderID)
		assert.Equal(t, uint64(10001), found_order.UserID)
		assert.Equal(t, float64(1000), found_order.Amount)
		assert.Equal(t, true, found_order.IsValid)
	})

	// 清空环境
	if err := dbConn.Exec("DELETE FROM orders").Error; err != nil {
		t.Fatal(err)
	}
	if err := dbConn.Exec("DELETE FROM users").Error; err != nil {
		t.Fatal(err)
	}
}

func TestOrderRepository_UpdateValidity(t *testing.T) {
	// 连接db
	dbConn := setupTestOrderDB(t)
	repo := db.NewGormOrderRepository(dbConn)
	user_repo := db.NewGormUserRepository(dbConn)

	t.Run("更新订单失效成功", func(t *testing.T) {
		// 由于外键约束，先把用户存进去
		user := &models.User{
			ID:               uint64(10001),
			Name:             "test",
			Email:            "test@example.com",
			TotalConsumption: 2000,
		}
		_, err := user_repo.Save(user)
		assert.NoError(t, err)

		order := &models.Order{
			OrderID: uint64(2025052110001),
			UserID:  uint64(10001),
			Amount:  float64(1000),
			IsValid: true,
		}
		_, err = repo.Save(order)
		assert.NoError(t, err)

		// 执行更新
		rows, err := repo.UpdateValidity(order.OrderID, false)
		assert.NoError(t, err)
		assert.Equal(t, int8(1), rows) // 验证影响行数

		// 验证数据
		found_order, err := repo.FindByID(order.OrderID)
		assert.NoError(t, err)
		assert.Equal(t, uint64(2025052110001), found_order.OrderID)
		assert.Equal(t, uint64(10001), found_order.UserID)
		assert.Equal(t, float64(1000), found_order.Amount)
		assert.Equal(t, false, found_order.IsValid)
	})

	t.Run("更新订单已失效或是不存在的订单", func(t *testing.T) {
		order := &models.Order{
			OrderID: uint64(2025052110001), // 上一个测试样例存储的
		}

		// 执行更新
		rows, err := repo.UpdateValidity(order.OrderID, false)
		assert.NoError(t, err)
		assert.Equal(t, int8(0), rows) // 验证影响行数

		// 验证数据
		found_order, err := repo.FindByID(order.OrderID)
		assert.NoError(t, err)
		assert.Equal(t, uint64(2025052110001), found_order.OrderID)
		assert.Equal(t, uint64(10001), found_order.UserID)
		assert.Equal(t, float64(1000), found_order.Amount)
		assert.Equal(t, false, found_order.IsValid)
	})

	// 清空环境
	if err := dbConn.Exec("DELETE FROM orders").Error; err != nil {
		t.Fatal(err)
	}
	if err := dbConn.Exec("DELETE FROM users").Error; err != nil {
		t.Fatal(err)
	}
}
