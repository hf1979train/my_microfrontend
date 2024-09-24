package db

import (
	"fmt"
	"log"
	"myhomesv/pkg/myenv"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDBAll(env myenv.Mysql) (*gorm.DB, error) {
	//func NewDBAll(path_env string) (*gorm.DB, *gorm.DB) {
	db_budget, err := NewDB(env, "MYSQL_BUDGET")
	//	return db_fx, db_budget
	return db_budget, err
}

func NewDB(env myenv.Mysql, label_db_name string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		env.User, env.Password, env.Host, env.Port, env.TableNameBudgeet)

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
		return nil, err
	}
	fmt.Println("Connceted")
	return db, nil
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
