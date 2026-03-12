package config

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB 根据配置初始化数据库，并返回 *gorm.DB
func InitDB(cfg *Config) (*gorm.DB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("配置不能为空")
	}
	var db *gorm.DB
	var err error
	switch cfg.Database.Driver {
	case "mysql":
		m := cfg.Database.Mysql
		if m.Port == 0 {
			m.Port = 3306 // 默认端口
		}
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			m.User,
			m.Password,
			m.Host,
			m.Port,
			m.DBName,
		)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgres":
		p := cfg.Database.Postgres
		if p.Port == 0 {
			p.Port = 5432
		}
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
			p.Host,
			p.User,
			p.Password,
			p.DBName,
			p.Port,
		)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("不支持的数据库驱动程序: %s", cfg.Database.Driver)
	}
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}
	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层 sql.DB 失败: %w", err)
	}
	// 连接池配置，可根据需要调整
	sqlDB.SetMaxOpenConns(25)           // 最大打开连接数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)
	return db, nil
}
