package usecase

import (
	"myhomesv/internal/repository"
	"myhomesv/pkg/myenv"
)

type IBudgetUsecase interface {
	IAuthUsecase
	// ここにインターフェースを登録する
}

type BudgetUsecase struct {
	br        repository.IBudgetRepository
	env_sv    myenv.Server
	env_gmail myenv.Gmail
}

func NewBudgetUsecase(br repository.IBudgetRepository, env_sv myenv.Server, env_gmail myenv.Gmail) IBudgetUsecase {
	return &BudgetUsecase{br, env_sv, env_gmail}
}
