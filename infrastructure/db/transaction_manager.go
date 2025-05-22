// infrastructure/db/transaction_manager.go
package db

import (
	"github.com/NorioKe/mysql_demo_use_gorm/domain/repositories"

	"gorm.io/gorm"
)

type GormTransactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) repositories.TransactionManager {
	return &GormTransactionManager{db: db}
}

func (m *GormTransactionManager) Transaction(fn func() error) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		return fn()
	})
}
