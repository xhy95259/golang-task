package config

import "time"

// Config 应用配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	Charset  string
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
}

// GetConfig 返回应用配置
func GetConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8090",
		},
		Database: DatabaseConfig{
			Driver:   "mysql",
			Host:     "localhost",
			Port:     "3306",
			Username: "root",
			Password: "root",
			DBName:   "blog_db2",
			Charset:  "utf8mb4",
		},
		JWT: JWTConfig{
			Secret:    "wsykxhy999",
			ExpiresIn: 24 * time.Hour,
		},
	}
}
