package dao

import (
	"backend/app/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var db *gorm.DB

func init() {
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPwd := os.Getenv("MYSQL_PWD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	//mysqlUser := "root"
	//mysqlPwd := "{xISlh^j$MD9*j@O"
	//mysqlHost := "35.224.29.75:3306"
	//mysqlHost := "unix(/cloudsql/term2-harutaka-kohama:us-central1:uttc)"
	//mysqlDatabase := "hackathon"

	//mysqlUser := "test_user"
	//mysqlPwd := "password"
	//mysqlHost := "localhost:3306"
	//mysqlDatabase := "test_database"

	//dsn := fmt.Sprint("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia/Tokyo", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)
	dsn := fmt.Sprintf("%s:%s@%s/%s?charset=utf8mb4&parseTime=true&loc=Asia%2FTokyo", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)
	//dsn := "root:6xQ}qHkc\"tf}uvLH@tcp(35.224.29.75:3306)/hackathon?charset=utf8mb4&parseTime=True&loc=Asia%2FTokyo"

	//dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("fail: sql.Open, %v\n", err)
	}

	if err := db.AutoMigrate(&model.Account{}); err != nil {
		log.Fatalf("fail: db.AutoMigrate, %v\n", err)
	}
	if err := db.AutoMigrate(&model.Contribution{}); err != nil {
		log.Fatalf("fail: db.AutoMigrate, %v\n", err)
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("fail: db.AutoMigrate, %v\n", err)
	}
	if err := db.AutoMigrate(&model.Workspace{}); err != nil {
		log.Fatalf("fail: db.AutoMigrate, %v\n", err)
	}
}
