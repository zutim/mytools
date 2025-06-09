package log2

import (
	"fmt"
	"go.uber.org/zap"
	"strconv"
)

type LogPathOption struct {
	LogPre     string
	TenantId   any
	ModuleName string
}

func WithModuleName(name string) func(*LogPathOption) {
	return func(o *LogPathOption) {
		o.ModuleName = name
	}
}

func WithTenantId(tenantId any) func(*LogPathOption) {
	return func(o *LogPathOption) {
		o.TenantId = tenantId
	}
}

func WithLogPre(pre string) func(*LogPathOption) {
	return func(o *LogPathOption) {
		o.LogPre = pre
	}
}

func GetLogPath(options ...func(*LogPathOption)) string {
	defaultOptions := LogPathOption{LogPre: "./", TenantId: "1", ModuleName: "default"}

	for _, setter := range options {
		setter(&defaultOptions)
	}
	return fmt.Sprintf("%s%s_%s.log", defaultOptions.LogPre, GetStringTenantId(defaultOptions.TenantId), defaultOptions.ModuleName)
}

func GetStringTenantId(tenantId any) string {
	path := "1"
	switch v := tenantId.(type) {
	case string:
		// 如果 tenantId 是 string 类型，直接使用
		path = v
	case int:
		// 如果 tenantId 是 int 类型，先将其转换为 string
		path = strconv.Itoa(v)
	default:
		// 如果 tenantId 是其他类型，处理错误情况
		fmt.Println("Error: tenantId is neither string nor int")
		return path
	}
	return path
}

func InitLog(path any) *zap.Logger {
	path2 := path.(string)
	return InitLogger((func(options *LoggerOptions) {
		options.Path = path2
	}))
}
