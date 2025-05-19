package models

import "errors"

type User struct {
	ID               uint64
	Name             string
	Email            string
	TotalConsumption float64
}

// AddConsumption: 增加消费总额
func (u *User) AddConsumption(amount float64) error {
	if amount < 0 {
		return errors.New("金额不能为负值")
	}
	u.TotalConsumption += amount
	return nil
}

// DeductConsumption: 减少消费总额
func (u *User) DeductConsumption(amount float64) error {
	if u.TotalConsumption < amount {
		return errors.New("消费总额不足")
	}
	u.TotalConsumption -= amount
	return nil
}
