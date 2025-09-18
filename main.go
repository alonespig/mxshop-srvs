package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
)

func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}

func main() {
	// dsn := "root:Root123456!@tcp(127.0.0.1:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"

	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags),
	// 	logger.Config{
	// 		SlowThreshold: time.Second,
	// 		LogLevel:      logger.Info,
	// 		Colorful:      true,
	// 	},
	// )

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	// 	Logger: newLogger,
	// 	NamingStrategy: schema.NamingStrategy{
	// 		SingularTable: true,
	// 	},
	// })

	// if err != nil {
	// 	panic(err)
	// }

	// _ = db.AutoMigrate(&model.User{})

	fmt.Println(genMd5("123456"))
}