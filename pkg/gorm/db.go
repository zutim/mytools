package db

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type ResolverConf struct {
	Dsn    string
	Tables []string
}

type Conf struct {
	Dsn         string
	MaxIdle     int
	MaxOpen     int
	MaxLifeTime int
	ResolverConf
	Log logger.Interface
}

func NewDb(conf *Conf) *gorm.DB {

	instance := New()

	if err := instance.Connect(conf.Dsn, &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 conf.Log,
	}); err != nil {
		panic(err)
	}

	if err := instance.EnableConnectionPool(conf.MaxIdle, conf.MaxOpen, time.Duration(conf.MaxLifeTime)*time.Second); err != nil {
		panic(err)
	}

	if conf.ResolverConf.Dsn != "" {
		dbresolvertmp := dbresolver.Config{
			Sources:           []gorm.Dialector{mysql.Open(conf.ResolverConf.Dsn)},
			TraceResolverMode: true,
		}

		args := make([]interface{}, len(conf.ResolverConf.Tables))
		for i, v := range conf.ResolverConf.Tables {
			args[i] = v
		}

		if err := instance.RegisterResolverConfig(dbresolvertmp, args...); err != nil {
			panic(err)
		}
	}

	return instance.DB
}
