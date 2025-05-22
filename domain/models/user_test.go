package models_test

import (
	"testing"

	"github.com/NorioKe/mysql_demo_use_gorm/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestUser_CreateOrder(t *testing.T) {
	// 正常情况测试
	t.Run("用户正确创建订单", func(t *testing.T) {
		user := &models.User{ID: 1001}
		amount := 100

		order, err := user.CreateOrder(user.ID, float64(amount))
		// 比对
		assert.NoError(t, err)
		assert.Equal(t, user.ID, order.UserID)
		assert.Equal(t, float64(amount), order.Amount)
		assert.True(t, order.IsValid)
	})
	t.Run("金额为负数应返回错误", func(t *testing.T) {
		user := models.User{ID: 1002}
		amount := -50.0

		order, err := user.CreateOrder(user.ID, amount)

		assert.ErrorContains(t, err, "消费金额不能为负数")
		assert.Nil(t, order)
	})
}

func TestUser_AddConsumption(t *testing.T) {
	// 正常情况测试
	t.Run("PositiveAmount", func(t *testing.T) {
		user := &models.User{TotalConsumption: 100}
		if err := user.AddConsumption(50); err != nil {
			t.Fatalf("添加金额失败: %v", err)
		}
		if user.TotalConsumption != 150 {
			t.Errorf("期望 150，实际 %.2f", user.TotalConsumption)
		}
	})

	// 异常情况测试
	t.Run("NegativeAmount", func(t *testing.T) {
		user := &models.User{TotalConsumption: 100}
		if err := user.AddConsumption(-120); err == nil {
			t.Error("余额不足错误未触发")
		}
	})
}
