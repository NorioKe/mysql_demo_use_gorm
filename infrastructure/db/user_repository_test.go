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

func setupTestUserDB(t *testing.T) *gorm.DB {
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
	if !dbConn.Migrator().HasTable(&models.User{}) {
		if err := dbConn.Migrator().CreateTable(&models.User{}); err != nil {
			t.Fatal(err)
		}
	} else {
		// 已有表时检查字段差异
		if err := dbConn.AutoMigrate(&models.User{}); err != nil {
			t.Fatal(err)
		}
	}

	// 清空表时使用Delete替代TRUNCATE
	if err := dbConn.Exec("DELETE FROM users").Error; err != nil {
		t.Fatal(err)
	}

	return dbConn
}

func TestUserRepository_FindByID(t *testing.T) {
	// 连接db
	dbConn := setupTestUserDB(t)
	repo := db.NewGormUserRepository(dbConn)

	t.Run("成功查找到用户ID", func(t *testing.T) {
		user := &models.User{
			ID:               uint64(1001),
			Name:             "test",
			Email:            "test@example.com",
			TotalConsumption: 2000,
		}

		// 执行保存
		if _, err := repo.Save(user); err != nil {
			t.Fatal(err)
		}

		// 查找
		foundUser, err := repo.FindByID(user.ID)
		// 验证数据
		assert.NoError(t, err)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, "test", foundUser.Name)
		assert.Equal(t, "test@example.com", foundUser.Email)
		assert.Equal(t, float64(2000), foundUser.TotalConsumption)
	})

	t.Run("无法找到用户ID", func(t *testing.T) {
		// 查找
		_, err := repo.FindByID(uint64(1002))
		assert.ErrorIs(t, err, repositories.ErrorNotFound)
	})

	// 清空环境
	if err := dbConn.Exec("DELETE FROM users").Error; err != nil {
		t.Fatal(err)
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	// 连接db
	dbConn := setupTestUserDB(t)
	repo := db.NewGormUserRepository(dbConn)

	t.Run("成功查找到用户邮箱", func(t *testing.T) {
		user := &models.User{
			ID:               uint64(1001),
			Name:             "test",
			Email:            "test@example.com",
			TotalConsumption: 2000,
		}

		// 执行保存
		if _, err := repo.Save(user); err != nil {
			t.Fatal(err)
		}

		// 查找
		foundUser, err := repo.FindByEmail(user.Email)
		// 验证数据
		assert.NoError(t, err)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, "test", foundUser.Name)
		assert.Equal(t, "test@example.com", foundUser.Email)
		assert.Equal(t, float64(2000), foundUser.TotalConsumption)
	})

	t.Run("无法找到用户邮箱", func(t *testing.T) {
		// 查找
		_, err := repo.FindByEmail("cannot@find.com")
		assert.ErrorIs(t, err, repositories.ErrorNotFound)
	})

	// 清空环境
	if err := dbConn.Exec("DELETE FROM users").Error; err != nil {
		t.Fatal(err)
	}
}

func TestUserRepository_Save(t *testing.T) {
	// 连接db
	dbConn := setupTestUserDB(t)
	repo := db.NewGormUserRepository(dbConn)

	t.Run("成功保存用户", func(t *testing.T) {
		user := &models.User{
			Name:             "test",
			Email:            "test@example.com",
			TotalConsumption: 0,
		}

		// 执行保存
		userID, err := repo.Save(user)
		assert.NoError(t, err)
		assert.NotZero(t, userID)

		// 验证数据
		foundUser, err := repo.FindByID(userID)
		assert.NoError(t, err)
		assert.Equal(t, "test", foundUser.Name)
		assert.Equal(t, "test@example.com", foundUser.Email)
		assert.Equal(t, float64(0), foundUser.TotalConsumption)
	})

	// 清空环境
	if err := dbConn.Exec("DELETE FROM users").Error; err != nil {
		t.Fatal(err)
	}
}

func TestUserRepository_UpdateTotalConsumption(t *testing.T) {
	// 连接db
	dbConn := setupTestUserDB(t)
	repo := db.NewGormUserRepository(dbConn)

	t.Run("成功更新用户的消费总额", func(t *testing.T) {
		user := &models.User{
			Name:             "test",
			Email:            "test@example.com",
			TotalConsumption: 0,
		}
		userID, err := repo.Save(user)
		assert.NoError(t, err)

		// 更新
		user.TotalConsumption += 1000
		affected_rows, err := repo.UpdateTotalConsumption(user)
		assert.NoError(t, err)
		assert.Equal(t, int8(1), int8(affected_rows))

		// 验证
		foundUser, err := repo.FindByID(userID)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.Email, foundUser.Email)
		assert.Equal(t, user.Name, foundUser.Name)
		assert.Equal(t, user.TotalConsumption, foundUser.TotalConsumption)
	})

	t.Run("更新用户消费总额时发现找不到用户", func(t *testing.T) {
		user := &models.User{
			ID:               888,
			Name:             "test",
			Email:            "test@example.com",
			TotalConsumption: 0,
		}
		_, err := repo.Save(user)
		assert.NoError(t, err)

		// 修改ID
		user.ID = 999
		affected_num, err := repo.UpdateTotalConsumption(user)
		assert.NoError(t, err)
		assert.Equal(t, int8(0), affected_num)
	})

	// 清空环境
	if err := dbConn.Exec("DELETE FROM users").Error; err != nil {
		t.Fatal(err)
	}
}
