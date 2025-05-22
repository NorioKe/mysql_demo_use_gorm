package db

import (
	"errors"

	"github.com/NorioKe/mysql_demo_use_gorm/domain/models"
	"github.com/NorioKe/mysql_demo_use_gorm/domain/repositories"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) repositories.UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) FindByID(id uint64) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repositories.ErrorNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repositories.ErrorNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) Save(user *models.User) (uint64, error) {
	if err := r.db.Save(user).Error; err != nil {
		return uint64(0), err
	}
	return user.ID, nil
}

func (r *GormUserRepository) UpdateTotalConsumption(user *models.User) (int8, error) {
	result := r.db.Model(user).
		Where("id = ?", user.ID).
		Update("total_consumption", user.TotalConsumption)
	affected_num, err := result.RowsAffected, result.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return int8(0), repositories.ErrorNotFound
		}
		return int8(0), err
	}
	return int8(affected_num), nil
}
