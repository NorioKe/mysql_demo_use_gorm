package models

import (
	"testing"
)

func TestOrder_Invalidate(t *testing.T) {
	order := &Order{Amount: 200, IsValid: true}

	// 有效订单失效测试
	t.Run("ValidOrder", func(t *testing.T) {
		amount, err := order.Invalidate()
		if err != nil {
			t.Fatal(err)
		}
		if amount != 200 || order.IsValid {
			t.Error("失效逻辑异常")
		}
	})

	// 重复失效测试
	t.Run("AlreadyInvalid", func(t *testing.T) {
		if _, err := order.Invalidate(); err == nil {
			t.Error("预期错误但未触发")
		}
	})
}
