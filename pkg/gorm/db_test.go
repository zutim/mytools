package db

import (
	"fmt"
	"testing"
)

func TestDb(t *testing.T) {

	db := NewDb(&Conf{
		Dsn:         "demo:demo123@tcp(192.168.0.178:3306)/onebet?charset=utf8mb4&parseTime=True&loc=Local",
		MaxOpen:     4,
		MaxIdle:     10,
		MaxLifeTime: 100,
		ResolverConf: ResolverConf{
			Dsn:    "demo:demo123@tcp(192.168.0.178:3306)/onebet_sport?charset=utf8mb4&parseTime=True&loc=Local",
			Tables: []string{"tbl_tickets"},
		},
		//Log: logger.New(
		//	log.Writer{
		//		Log: log.NewLogMap().WithOptionPath(LoggerOptions{}),
		//	},
		//	logger.Config{
		//		SlowThreshold:             200 * time.Millisecond,
		//		IgnoreRecordNotFoundError: true,
		//		Colorful:                  true,
		//		LogLevel:                  logger.Warn,
		//	}),
	})

	type Tickets struct {
		Id string
	}

	var ti Tickets

	if err := db.Table("tbl_tickets").First(&ti).Error; err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(ti.Id)

	type Users struct {
		Id       int    `json:"id"`
		Mobile   string `json:"mobile"`
		Nickname string `json:"nickname"`
	}

	var u Users
	if err := db.Table("tbl_users").First(&u).Error; err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(u)
}
