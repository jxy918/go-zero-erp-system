// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	DB struct {
		DataSource string
	}
	JWT struct {
		AccessSecret string
		AccessExpire int64
	}
	// CORS 配置
	CORS struct {
		AllowedOrigins string // 允许的域名，逗号分隔，如 "http://localhost:3000,https://example.com"
	}
	// 自定义业务指标配置
	Metrics struct {
		Enabled bool `json:",default=true"` // 是否启用自定义业务指标
	}
}
