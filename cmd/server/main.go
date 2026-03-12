package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zxc7563598/oneadmin/internal/bootstrap"
	"github.com/zxc7563598/oneadmin/internal/config"
	"github.com/zxc7563598/oneadmin/internal/migrate"
)

func main() {
	port := flag.Int("port", 9000, "服务端口")
	configPath := flag.String("config", "config.yaml", "配置文件路径")
	// 加载配置
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("无法加载配置: %v", err)
	}
	// 初始化数据库
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("无法初始化数据库: %v", err)
	}
	// 自动建表
	if err := migrate.Run(db); err != nil {
		panic(err)
	}
	// 初始化填充数据
	if err := migrate.Seed(db); err != nil {
		panic(err)
	}
	// 初始化 Gin 应用
	r := bootstrap.NewApp(db)
	// 创建 HTTP Server
	addr := fmt.Sprintf(":%d", *port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	// 启动服务
	go func() {
		log.Printf("服务在 %s 启动\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("监听错误: %v", err)
		}
	}()
	// 监听退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		syscall.SIGINT,  // ctrl+c
		syscall.SIGTERM, // docker stop
	)
	<-quit
	log.Println("关闭服务...")
	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("服务被迫关闭: %v", err)
	}
	log.Println("服务已退出")
}
