package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"sync"
)

// DB 定义一个全局变量用于service层使用,在main中初始化
var (
	mysqlOnce sync.Once
	db        *gorm.DB
	dbErr     error
)

func Init() {
	mysqlOnce.Do(func() {
		username := "root"        //账号
		password := "828119"      //密码
		host := "192.168.163.133" //数据库地址，可以是Ip或者域名
		port := 3306              //数据库端口
		Dbname := "shopping-go"   //数据库名
		timeout := "10s"          //连接超时，10秒

		// root:root@tcp(127.0.0.1:3306)/gorm?
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
		db, dbErr = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info), // 打印SQL日志（调试用）
		})
		if dbErr != nil {
			panic("连接数据库失败, error=" + dbErr.Error())
		}
		// 连接成功
		log.Println("连接数据库成功，数据库名：" + Dbname)
	})
}

func DB() *gorm.DB {
	return db
}
