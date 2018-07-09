package config

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"strings"
)

var defaultConf = []byte(`
core:
  enabled: true
  mode: "debug"
  name: "apiserver"
  address: ""
  port: "9090"
  max_ping_count: 2
  jwt_secret: "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5"
  tls:
    port: "9098"
    cert_path: "src/config/server.crt"
    key_path: "src/config/server.key"
  auto_tls:
    enabled: false
    folder: ".cache"
    host: ""

log:
  writers: "file,stdout"
  logger_level: "DEBUG"
  logger_file: "log/apiserver.log"
  log_format_text: false
  rollingPolicy: "size"
  log_rotate_date: 1
  log_rotate_size: 1
  log_backup_count: 7

db:
  name: "db_apiserver"
  addr: "127.0.0.1:3306"
  username: "root"
  password: "root"
`)

type Config struct {
	Name string
}

type ConfYaml struct {
	Core SectionCore `yaml:"core"`
	Log  SectionLog  `yaml:"log"`
	Db   SectionDb   `yaml:"db"`
}

// SectionCore is sub section of config.
type SectionCore struct {
	Enabled			 bool						`yaml:"enabled"`
	Mode         string         `yaml:"mode"`
	Name         string         `yaml:"name"`
	Address      string         `yaml:"address"`
	Port         string         `yaml:"port"`
	MaxPingCount int            `yaml:"max_ping_count"`
	JwtSecret    string         `yaml:"jwt_secret"`
	TLS          SectionTLS     `yaml:"tls"`
	AutoTLS      SectionAutoTLS `yaml:"auto_tls"`
}

// SectionTLS support tls
type SectionTLS struct {
	Port       string `yaml:"port"`
	CertPath   string `yaml:"cert_path"`
	KeyPath    string `yaml:"key_path"`
}

// SectionAutoTLS support Let's Encrypt setting.
type SectionAutoTLS struct {
	Enabled bool   `yaml:"enabled"`
	Folder  string `yaml:"folder"`
	Host    string `yaml:"host"`
}

// SectionLog is sub section of config.
type SectionLog struct {
	Writers        string `yaml:"writers"`
	LoggerLevel    string `yaml:"logger_level"`
	LoggerFile     string `yaml:"logger_file"`
	LogFormatText  bool   `yaml:"log_format_text"`
	RollingPolicy  string `yaml:"rollingPolicy"`
	LogRotateDate  int    `yaml:"log_rotate_date"`
	LogRotateSize  int    `yaml:"log_rotate_size"`
	LogBackupCount int    `yaml:"log_backup_count"`
}

// SectionDb is sub section of config.
type SectionDb struct {
	Name     string `yaml:"name"`
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func Init(cfg string) (ConfYaml, error) {
	var confYaml ConfYaml
	c := Config{
		Name: cfg,
	}

	// 初始化配置文件
	confYaml, err := c.initConfig()
	if err != nil {
		return confYaml, nil
	}

	// 初始化日志包
	c.initLog(&confYaml)

	// 监控配置文件变化并热加载程序
	c.watchConfig()

	return confYaml, nil
}

// 初始化配置文件
func (c *Config) initConfig() (ConfYaml, error) {
	var confYaml ConfYaml

	// 设置配置文件格式为YAML
	viper.SetConfigType("yaml")

	// 读取匹配的环境变量
	viper.AutomaticEnv()

	// 读取环境变量的前缀为APISERVER
	viper.SetEnvPrefix("APISERVER")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if c.Name != "" {
		// 如果指定了配置文件，则解析指定的配置文件
		viper.SetConfigFile(c.Name)
	} else {
		// 如果没有指定配置文件，则解析默认的配置文件
		viper.AddConfigPath("src/config")
		viper.SetConfigName("config")
	}

	// viper解析配置文件
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		// load default config
		viper.ReadConfig(bytes.NewBuffer(defaultConf))
	}

	// Core
	confYaml.Core.Enabled = viper.GetBool("core.enabled")
	confYaml.Core.Mode = viper.GetString("core.mode")
	confYaml.Core.Name = viper.GetString("core.name")
	confYaml.Core.Address = viper.GetString("core.")
	confYaml.Core.Port = viper.GetString("core.port")
	confYaml.Core.MaxPingCount = viper.GetInt("core.max_ping_count")
	confYaml.Core.JwtSecret = viper.GetString("core.jwt_secret")
	confYaml.Core.TLS.Port = viper.GetString("core.tls.port")
	confYaml.Core.TLS.CertPath = viper.GetString("core.tls.cert_path")
	confYaml.Core.TLS.KeyPath = viper.GetString("core.tls.key_path")
	confYaml.Core.AutoTLS.Enabled = viper.GetBool("core.auto_tls.enabled")
	confYaml.Core.AutoTLS.Folder = viper.GetString("core.auto_tls.folder")
	confYaml.Core.AutoTLS.Host = viper.GetString("core.auto_tls.host")

	// Log
	confYaml.Log.Writers = viper.GetString("log.writers")
	confYaml.Log.LoggerLevel = viper.GetString("log.logger_level")
	confYaml.Log.LoggerFile = viper.GetString("log.logger_file")
	confYaml.Log.LogFormatText = viper.GetBool("log.log_format_text")
	confYaml.Log.RollingPolicy = viper.GetString("log.rollingPolicy")
	confYaml.Log.LogRotateDate = viper.GetInt("log.log_rotate_date")
	confYaml.Log.LogRotateSize = viper.GetInt("log.log_rotate_size")
	confYaml.Log.LogBackupCount = viper.GetInt("log.log_backup_count")

	// Db
	confYaml.Db.Name = viper.GetString("db.name")
	confYaml.Db.Addr = viper.GetString("db.addr")
	confYaml.Db.Username = viper.GetString("db.username")
	confYaml.Db.Password = viper.GetString("db.password")

	return confYaml, nil
}

// 初始化日志包
func (c *Config) initLog(conf *ConfYaml) {
	passLagerCfg := log.PassLagerCfg{
		Writers:        conf.Log.Writers,
		LoggerLevel:    conf.Log.LoggerLevel,
		LoggerFile:     conf.Log.LoggerFile,
		LogFormatText:  conf.Log.LogFormatText,
		RollingPolicy:  conf.Log.RollingPolicy,
		LogRotateDate:  conf.Log.LogRotateDate,
		LogRotateSize:  conf.Log.LogRotateSize,
		LogBackupCount: conf.Log.LogBackupCount,
	}
	log.InitWithConfig(&passLagerCfg)
}

// 监控配置文件变化并热加载程序
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})
}
