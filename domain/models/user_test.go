package models

import (
	"testing"
)

func TestUser_AddConsumption(t *testing.T) {
	// 正常情况测试
	t.Run("PositiveAmount", func(t *testing.T) {
		user := &User{TotalConsumption: 100}
		if err := user.AddConsumption(50); err != nil {
			t.Fatalf("添加金额失败: %v", err)
		}
		if user.TotalConsumption != 150 {
			t.Errorf("期望 150，实际 %.2f", user.TotalConsumption)
		}
	})

	// 异常情况测试
	t.Run("NegativeAmount", func(t *testing.T) {
		user := &User{TotalConsumption: 100}
		if err := user.AddConsumption(-20); err == nil {
			t.Error("预期错误但未触发")
		}
	})
}
