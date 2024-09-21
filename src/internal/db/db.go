package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDBAll(path_env string) *gorm.DB {
	//func NewDBAll(path_env string) (*gorm.DB, *gorm.DB) {
	db_budget := NewDB(path_env, "MYSQL_BUDGET")
	//	return db_fx, db_budget
	return db_budget
}

func NewDB(path_env string, label_db_name string) *gorm.DB {
	err := godotenv.Load(path_env)
	if err != nil {
		log.Fatalln(err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PW"), os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"), os.Getenv(label_db_name))

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	//	db.Logger = db.Logger.LogMode(logger.Info)

	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connceted")
	return db
}

func CLoseDBAll(db_budget *gorm.DB) {
	//func CLoseDBAll(db_fx *gorm.DB, db_budget *gorm.DB) {
	// CloseDB(db_fx)
	CloseDB(db_budget)
}

func CloseDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Fatalln(err)
	}
}

// https://gorm.io/ja_JP/docs/connecting_to_the_database.html
// [logger]https://qiita.com/isaka1022/items/4b37481ec216e2fbf507
