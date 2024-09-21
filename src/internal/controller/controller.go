package controller

import (
	"myhomesv/internal/usecase"
)

type IBudgetController interface {
	IAuthController
	// ここにインターフェースを登録する
}

type budgetController struct {
	bu usecase.IBudgetUsecase
}

func NewBudgetController(bu usecase.IBudgetUsecase) IBudgetController {
	return &budgetController{bu}
}
