package models

import (
	"errors"
	"strings"
	// "gorm.io/gorm"
)

type User struct {
	// gorm.Model  // 这个会引入CreatedAt、UpdatedAt等字段从而改变表结构
	ID               uint64  `gorm:"primaryKey;autoIncrement"`
	Name             string  `gorm:"type:varchar(100)"`
	Email            string  `gorm:"uniqueIndex;type:varchar(255)"` // 明确指定类型和长度
	TotalConsumption float64 `gorm:"type:decimal(12,2);default:0"`
}

// CreateUser: 创建用户
func CreateUser(name string, email string) (*User, error) {
	if isValidEmail(email) == false {
		return nil, errors.New("邮箱格式不正确")
	}
	return &User{
		ID:               uint64(0),
		Name:             name,
		Email:            email,
		TotalConsumption: 0,
	}, nil
}

// CreateOrder: 用户创建订单
func (u *User) CreateOrder(userid uint64, amount float64) (*Order, error) {
	if amount < 0 {
		return nil, errors.New("消费金额不能为负数")
	}

	return &Order{
		OrderID: 0,
		UserID:  userid,
		Amount:  amount,
		IsValid: true,
	}, nil
}

// AddConsumption: 修改消费总额
func (u *User) AddConsumption(amount float64) error {
	u.TotalConsumption += amount
	if u.TotalConsumption < 0 {
		return errors.New("消费总额不足")
	}
	return nil
}

func isValidEmail(email string) bool {
	// 检查后缀
	suffixes := []string{"@qq.com", "@163.com", "@example.com"}
	for _, suffix := range suffixes {
		if strings.HasSuffix(email, suffix) {
			// 确保@前至少有一个字符，且整个字符串只有一个@
			atIndex := len(email) - len(suffix)
			return atIndex > 0 && !strings.Contains(email[:atIndex], "@")
		}
	}
	return false
}
