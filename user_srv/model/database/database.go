package main

import (
	"crypto/rand"
	"fmt"
	"go-gin-rme/user_srv/model"
	"go-gin-rme/user_srv/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"math/big"
	"net/url"
	"os"
	"time"
)

var (
	DB  *gorm.DB
	err error
)

func main() {
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=%s",
		"realy",
		"Chen1224",
		"rm-bp16v29co3893zkyplo.mysql.rds.aliyuncs.com",
		"3306",
		"remi_user_srv",
		"utf8",
		url.QueryEscape("Asia/Shanghai"))
	NewLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold: time.Second, //慢查询阈值
		LogLevel:      logger.Info, //log lever
		Colorful:      true,        //禁用彩色打印
	})

	// 连接数据库
	DB, err = gorm.Open(mysql.Open(args), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: NewLogger,
	})
	if err != nil {
		panic("failed to open database: " + err.Error())
	}
	for i := 0; i < 10; i++ {
		var p, _ = rand.Int(rand.Reader, big.NewInt(8999))
		user := model.User{
			NickName: fmt.Sprintf("op-%d", i),
			Mobile:   fmt.Sprintf("1531046%d", p.Int64()+1000),
			Password: util.BcryptHash("qwe123"),
		}
		DB.Save(&user)
	}

	//数据库迁移表，第一次启动后，可以注释掉
	err = DB.AutoMigrate(&model.User{})
	if err != nil {

		return
	}
	sqlDB, _ := DB.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

}
