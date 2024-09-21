package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"myhomesv/internal/domain/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var tables = []interface{}{
	&models.User{},
	&models.ResetToken{},
}

func main() {
	db := initDB()
	defer func() {
		fmt.Println("closing database connection")
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get database instance: %v", err)
		}
		sqlDB.Close()
		fmt.Println("database connection closed")
		fmt.Println("migration sucessful")
	}()
	// マイグレーション
	for _, table := range tables {
		if !db.Migrator().HasTable(table) {
			db.AutoMigrate(table)
		}
	}
}

var DB *gorm.DB

func initDB() *gorm.DB {
	// .envファイルを読み込む
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// 環境変数からデータベース接続情報を取得
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PW")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dbname := os.Getenv("MYSQL_BUDGET")

	// データベース名を指定せずに MySQL 接続用の DSN (Data Source Name) を作成
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, password, host, port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// データベースが存在しない場合は作成
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbname))
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	// データベース名を含むように DSN を更新します
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname)

	// データベースに接続
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return gormDB
}
