package main

import (
	"log"
	"myhomesv/internal/controller"
	"myhomesv/internal/db"
	"myhomesv/internal/repository"
	"myhomesv/internal/router"
	"myhomesv/internal/usecase"
	"myhomesv/pkg/myenv"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// env -----------------
	var myenv = initialize()

	// web server ----------
	db_budget, err := db.NewDBAll(myenv.Mysql)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	//budget
	budgetRepository := repository.NewBudgetRepository(db_budget)
	budgetUsecase := usecase.NewBudgetUsecase(budgetRepository, myenv.Server, myenv.Gmail)
	budgetController := controller.NewBudgetController(budgetUsecase)

	//HTTPS非対応のサーバを起動する場合-------------
	// ルーターを設定
	r := router.SetupRouter(budgetController)
	// サーバーをポート8080で起動
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to Run Server: %v", err)
	}

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

func initialize() myenv.Env {
	// env -----------------
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}
	return myenv.Env{
		Server: myenv.Server{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
		Mysql: myenv.Mysql{
			User:             os.Getenv("MYSQL_USER"),
			Password:         os.Getenv("MYSQL_PW"),
			Host:             os.Getenv("MYSQL_HOST"),
			Port:             os.Getenv("MYSQL_PORT"),
			TableNameBudgeet: os.Getenv("MYSQL_BUDGET"),
		},
		Gmail: myenv.Gmail{
			PathCredentials: os.Getenv("GMAIL_CREDENTIALS"),
			PathToken:       os.Getenv("GMAIL_TOKEN"),
		},
	}
}
