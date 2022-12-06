package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"wxbot/internal"
	"wxbot/internal/config"
	"wxbot/internal/global"
	"wxbot/internal/qun"
	"wxbot/pkg/api"
	"wxbot/pkg/logger"
)

func Execute() {
	Init()
	//启动微信监听
	sv, err := internal.NewWechatService(global.Config.ListenPort)
	if err != nil {
		log.Fatal(err)
	}
	defer sv.Stop()

	selfInfo, err := global.Api.RequestGetSelfInfo()
	if err != nil {
		log.Fatal(err)
	}
	logger.Infof("用户信息 %+v\n", selfInfo)
	global.Self = selfInfo.SelfInfo

	//订阅消息
	read, err := internal.NewMessageRead(global.Config.WeChatApiAddress)
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

// todo 优化
func Init() {
	// 配置文件加载
	conf, err := config.Load("conf/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	global.Config = conf

	//日志
	logger.InitDefault(conf.Log, conf.LogLevel)

	//api
	a := api.New(conf.WeChatApiAddress)
	global.Api = a

	//qun处理模块
	qun.DefaultHandler = qun.New(global.Config.QunManager)
}
