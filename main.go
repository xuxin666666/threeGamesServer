package main

import (
	"03.threeGames/controller"
	"03.threeGames/dao/mysql"
	"03.threeGames/logger"
	"03.threeGames/pkg/snowflake"
	"03.threeGames/routers"
	"03.threeGames/settings"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := settings.Init(); err != nil {
		fmt.Printf("settings init failed, err:%v\n", err)
		return
	}

	if err := logger.Init(); err != nil {
		fmt.Printf("logger init failed, err:%v\n", err)
		return
	}
	defer func() {
		if err :=zap.L().Sync(); err != nil {
			zap.L().Fatal("zap exiting failed", zap.Error(err))
		}
	}()
	zap.L().Debug("logger init success...")

	if err := mysql.Init(); err != nil {
		fmt.Printf("mysql init failed, err: %v\n", err)
		return
	}
	defer mysql.Close()

	if err := snowflake.Init(
		viper.GetString("startTime"), viper.GetInt64("machineID"),
	); err != nil {
		fmt.Printf("snowflake init failed, err: %v\n", err)
		return
	}

	// 初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
		return
	}

	// 注册路由
	r := routers.SetUp()
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", viper.GetInt("port")),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		fmt.Println("listening and serving HTTP on :", viper.GetInt("port"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: ", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<- quit // 阻塞在此，当接收到上述两种信号时才会往下执行

	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server showdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")

}




