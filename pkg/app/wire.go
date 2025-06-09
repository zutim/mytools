//go:build wireinject
// +build wireinject

//+gen:wire

package app

import (
	"github.com/zutim/mytools/pkg/config"
)

// InitializeApp 使用 wire 构建 App 实例
func InitializeApp(globalCfg *config.AppConfig) *App {
	app := GetDefaultApp()
	app.SetGlobalConfig(globalCfg)

	// 注册组件
	app.RegisterComponent(NewMysqlComponent())
	app.RegisterComponent(NewLogComponent())

	return app
}
