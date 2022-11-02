package internal

import (
	"embed"
	"errors"
	"fmt"
	"syscall"
	"wxbot/pkg/logger"

	"github.com/shirou/gopsutil/process"
)

var DLLFS embed.FS

type WeChatService struct {
	port int

	pid    int32
	driver *syscall.LazyDLL
}

func NewWechatService(port int) (*WeChatService, error) {
	sv := &WeChatService{
		port: port,
	}
	if err := sv.Start(); err != nil {
		return nil, err
	}

	return sv, nil
}

func (w *WeChatService) Start() error {
	bit := 32 << (^uint(0) >> 63)
	logger.Infof("开启微信监听 系统位数 %d", bit)
	if bit == 32 {
		w.driver = syscall.NewLazyDLL("dll/wxDriver.dll")
	} else {
		w.driver = syscall.NewLazyDLL("dll/wxDriver64.dll")
	}

	pid, err := w.findWehatOrNew()
	if err != nil {
		return err
	}
	w.pid = pid
	startListen := w.driver.NewProc("start_listen")
	startListen.Call(uintptr(pid), uintptr(w.port))

	return nil
}

func (w *WeChatService) Stop() error {
	logger.Infof("关闭微信监听")
	stopListen := w.driver.NewProc("stop_listen")
	stopListen.Call(uintptr(w.pid))

	return nil
}

func (w *WeChatService) findWehatOrNew() (int32, error) {
	pid, err := w.findWechatPid()
	if pid > 0 {
		logger.Infof("检测到微信已经打开 pid %d", pid)
		return pid, nil
	}
	if errors.Is(err, ErrorWechatPidNotFound) {
		newWechat := w.driver.NewProc("new_wechat")
		v, _, _ := newWechat.Call()
		logger.Infof("微信未打开，尝试开启微信 pid %d", pid)
		return int32(v), nil
	}

	return 0, err
}

var ErrorWechatPidNotFound = errors.New("wechat pid not found")

func (w *WeChatService) findWechatPid() (int32, error) {
	pids, err := process.Pids()
	if err != nil {
		return 0, err
	}
	for _, pid := range pids {
		process, err := process.NewProcess(pid)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		name, _ := process.Name()
		if name == "WeChat.exe" {
			return pid, nil
		}
	}

	return 0, ErrorWechatPidNotFound
}
