package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"wxbot/internal"
	"wxbot/internal/config"
	"wxbot/internal/global"
	"wxbot/pkg/api"
	"wxbot/pkg/logger"
)

func Execute() {
	//配置文件加载
	conf, err := config.Load("conf/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	global.Config = conf

	//日志初始化
	logger.InitDefault(conf.Log, conf.LogLevel)

	//启动微信监听
	sv, err := internal.NewWechatService(8000)
	if err != nil {
		log.Fatal(err)
	}
	defer sv.Stop()

	//测试api
	a := api.New(fmt.Sprintf("127.0.0.1:%d", conf.ListenPort))
	selfInfo, err := a.RequestGetSelfInfo()
	if err != nil {
		log.Fatal(err)
	}
	logger.Infof("用户信息 %+v\n", selfInfo)

	//订阅消息
	read, err := internal.NewMessageRead(conf.ListenPort)
	if err != nil {
		log.Fatal(err)
	}
	read.Start()
	defer read.Stop()

	// 退出命令监听
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}
