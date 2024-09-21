package main

import (
	"myhomesv/internal/controller"
	"myhomesv/internal/db"
	"myhomesv/internal/repository"
	"myhomesv/internal/router"
	"myhomesv/internal/usecase"
)

func main() {
	// web server ----------
	db_budget := db.NewDBAll(".env")
	//budget
	budgetRepository := repository.NewBudgetRepository(db_budget)
	budgetUsecase := usecase.NewBudgetUsecase(budgetRepository)
	budgetController := controller.NewBudgetController(budgetUsecase)

	//HTTPS非対応のサーバを起動する場合-------------
	// ルーターを設定
	r := router.SetupRouter(budgetController)
	// サーバーをポート8080で起動
	r.Run(":8080")

	//TODO: HTTPS対応のサーバを起動する場合-------------
	// TLS証明書とキーのパスを指定
	// certFile := "cfg/ssl/certfile.crt"
	// keyFile := "cfg/ssl/key.pem"
	// ルーターを設定
	// r := router.SetupRouter(budgetController, certFile, keyFile)
	// HTTPSサーバーを起動
	// if err := r.RunTLS(":443", certFile, keyFile); err != nil {
	// 	panic(err)
	// }
}
