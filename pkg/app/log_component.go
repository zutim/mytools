// pkg/app/log_component.go
package app

import (
	"mytool3/pkg/log2"
	"os"
	"path/filepath"
)

type LogComponent struct{}

func NewLogComponent() Component {
	return &LogComponent{}
}

func (l *LogComponent) Name() string {
	return "log"
}

func (l *LogComponent) Init(path any) (any, error) {
	// 获取全局配置
	app := GetDefaultApp()
	logCfg := app.GetGlobalConfig().Log
	path2 := path.(string)
	if err := os.MkdirAll(filepath.Dir(path2), os.ModePerm); err != nil {
		return nil, err
	}

	logger := log2.InitLogger(
		func(options *log2.LoggerOptions) {

			options.Compress = logCfg.Compress

			if path2 != "" {
				options.Path = path2
			}

			if logCfg.MaxSize > 0 {
				options.MaxSize = logCfg.MaxSize
			}

			if logCfg.MaxAge > 0 {
				options.MaxAge = logCfg.MaxAge
			}

			if logCfg.MaxBackups > 0 {
				options.MaxBackups = logCfg.MaxBackups
			}
		},
	)

	return logger.Sugar(), nil
}

func (l *LogComponent) Close(tenantId any) error {
	return nil
}

func (l *LogComponent) HealthCheck() bool {
	return true
}
