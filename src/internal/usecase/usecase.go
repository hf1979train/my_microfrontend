package usecase

import (
	"myhomesv/internal/repository"
)

type IBudgetUsecase interface {
	IAuthUsecase
	// ここにインターフェースを登録する
}

type BudgetUsecase struct {
	br repository.IBudgetRepository
}

func NewBudgetUsecase(br repository.IBudgetRepository) IBudgetUsecase {
	return &BudgetUsecase{br}
}
