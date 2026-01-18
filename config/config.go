package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire string `mapstructure:"expire"`
}

// func Load() *Config {
// 	// 简化配置加载，实际应该使用 Viper
// 	return &Config{
// 		Server: ServerConfig{
// 			Port: "8080",
// 			Host: "0.0.0.0",
// 			Mode: "debug",
// 		},
// 		Database: DatabaseConfig{
// 			Host:     "localhost",
// 			Port:     3306,
// 			Username: "root",
// 			Password: "password",
// 			DBName:   "mydb",
// 		},
// 		JWT: JWTConfig{
// 			Secret: "your-secret-key-change-in-production",
// 			Expire: "24h",
// 		},
// 	}
// }

func init() {
	// 设置配置文件名称（不含扩展名）
	viper.SetConfigName("config")
	// 设置配置文件类型
	viper.SetConfigType("yaml")
	// 添加配置文件搜索路径
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.app")

	// 读取环境变量
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	// 设置默认值
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.mode", "debug")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Error reading config file: %v", err)
		log.Println("Using default values and environment variables")
		log.Fatalln("Config load failed")
	} else {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}
}

func Load() *Config {
	var config Config

	// 使用 viper.Unmarshal 将配置数据解析到 config 结构体中
	// 如果解析失败，返回错误
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln("Config load failed")
	}

	// 返回解析成功的配置对象
	return &config
}
