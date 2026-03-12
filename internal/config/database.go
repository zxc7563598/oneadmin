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
		db, err = initMySQL(cfg)
	case "postgres":
		db, err = initPostgres(cfg)
	default:
		return nil, fmt.Errorf("不支持的数据库驱动程序: %s", cfg.Database.Driver)
	}
	if err != nil {
		return nil, err
	}
	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层 sql.DB 失败: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)
	return db, nil
}

// 初始化 MySQL
func initMySQL(cfg *Config) (*gorm.DB, error) {
	m := cfg.Database.Mysql
	if m.Port == 0 {
		m.Port = 3306
	}
	// 先连接 MySQL server（不指定数据库）
	serverDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
		m.User,
		m.Password,
		m.Host,
		m.Port,
	)
	serverDB, err := gorm.Open(mysql.Open(serverDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}
	// 创建数据库
	createSQL := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci",
		m.DBName,
	)
	if err := serverDB.Exec(createSQL).Error; err != nil {
		return nil, fmt.Errorf("创建数据库失败: %w", err)
	}
	// 连接目标数据库
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.User,
		m.Password,
		m.Host,
		m.Port,
		m.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}
	return db, nil
}

// 初始化 PostgreSQL
func initPostgres(cfg *Config) (*gorm.DB, error) {
	p := cfg.Database.Postgres
	if p.Port == 0 {
		p.Port = 5432
	}
	// 连接 postgres 默认数据库
	serverDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=postgres port=%d sslmode=disable",
		p.Host,
		p.User,
		p.Password,
		p.Port,
	)
	serverDB, err := gorm.Open(postgres.Open(serverDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}
	// 创建数据库（Postgres 没有 IF NOT EXISTS）
	createSQL := fmt.Sprintf(
		"CREATE DATABASE %s",
		p.DBName,
	)
	// 忽略数据库已存在的错误
	serverDB.Exec(createSQL)
	// 连接目标数据库
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		p.Host,
		p.User,
		p.Password,
		p.DBName,
		p.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}
	return db, nil
}
