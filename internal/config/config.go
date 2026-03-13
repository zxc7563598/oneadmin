package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type DatabaseMysqlConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type DatabasePostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type DatabaseConfig struct {
	Driver   string                 `yaml:"driver"`
	Mysql    DatabaseMysqlConfig    `yaml:"mysql"`
	Postgres DatabasePostgresConfig `yaml:"postgres"`
}

type RedisConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
}
type JWTConfig struct {
	Secret     string `yaml:"secret"`
	AccessTTL  int    `yaml:"access_ttl"`
	RefreshTTL int    `yaml:"refresh_ttl"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
}

// LoadConfig 解析 YAML
func LoadConfig(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	if err := yaml.Unmarshal(raw, c); err != nil {
		return nil, err
	}
	// 验证配置
	if err := ValidateConfig(c); err != nil {
		return nil, err
	}
	return c, nil
}

// ValidateConfig 验证配置的有效性
func ValidateConfig(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("配置不能为空")
	}
	switch cfg.Database.Driver {
	case "mysql":
		if cfg.Database.Mysql.Host == "" {
			return fmt.Errorf("mysql 配置错误：host 不能为空")
		}
		if cfg.Database.Mysql.Port == 0 {
			return fmt.Errorf("mysql 配置错误：port 不能为 0")
		}
		if cfg.Database.Mysql.User == "" {
			return fmt.Errorf("mysql 配置错误：user 不能为空")
		}
		if cfg.Database.Mysql.Password == "" {
			return fmt.Errorf("mysql 配置错误：password 不能为空")
		}
		if cfg.Database.Mysql.DBName == "" {
			return fmt.Errorf("mysql 配置错误：dbname 不能为空")
		}
	case "postgres":
		if cfg.Database.Postgres.Host == "" {
			return fmt.Errorf("postgres 配置错误：host 不能为空")
		}
		if cfg.Database.Postgres.Port == 0 {
			return fmt.Errorf("postgres 配置错误：port 不能为 0")
		}
		if cfg.Database.Postgres.User == "" {
			return fmt.Errorf("postgres 配置错误：user 不能为空")
		}
		if cfg.Database.Postgres.Password == "" {
			return fmt.Errorf("postgres 配置错误：password 不能为空")
		}
		if cfg.Database.Postgres.DBName == "" {
			return fmt.Errorf("postgres 配置错误：dbname 不能为空")
		}
	default:
		return fmt.Errorf("不支持的数据库驱动程序: %s", cfg.Database.Driver)
	}
	if len(cfg.JWT.Secret) < 32 {
		return fmt.Errorf("JWT 密钥长度不能低于 32 位")
	}
	if cfg.JWT.AccessTTL <= 0 {
		return fmt.Errorf("access ttl 必须大于 0")
	}
	if cfg.JWT.RefreshTTL <= 0 {
		return fmt.Errorf("refresh ttl 必须大于 0")
	}
	return nil
}
