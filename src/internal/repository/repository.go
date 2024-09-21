package repository

import "gorm.io/gorm"

type IBudgetRepository interface {
	IUserRepository
	// ここにインターフェースを登録する
}

type BudgetRepository struct {
	db *gorm.DB
}

func NewBudgetRepository(db *gorm.DB) IBudgetRepository {
	return &BudgetRepository{db}
}
