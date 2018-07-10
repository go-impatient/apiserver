package src

import (
	"github.com/moocss/apiserver/src/config"
	"github.com/moocss/apiserver/src/service"
)

var (
	// 版本
	Version = "v0.0.3"

	// 配置参数
	Conf config.ConfYaml

	// StatStorage implements the storage interface
	Storage service.Database
)
