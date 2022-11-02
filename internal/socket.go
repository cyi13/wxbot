package internal

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
	"wxbot/pkg/api"
	"wxbot/pkg/logger"
)

type MessageRead struct {
	notifyURL []string

	textPort int
	server   *api.Api
	done     chan struct{}
}

const (
	HOOK_TEXT = iota
)

func NewMessageRead(wechatPort int, notifyURL []string) (*MessageRead, error) {
	sv := api.New(fmt.Sprintf("127.0.0.1:%d", wechatPort))
	return &MessageRead{
		notifyURL: notifyURL,

		textPort: 10808,
		server:   sv,
		done:     make(chan struct{}),
	}, nil
}

func (m *MessageRead) Start() {
	logger.Infof("消息订阅启动")
	m.SubscribeText()
}

func (m *MessageRead) Stop() {
	logger.Infof("消息订阅关闭")
	close(m.done)
}

// 文本消息订阅
func (m *MessageRead) SubscribeText() error {
	logger.Infof("开始文件消息订阅 socket 端口 %d", m.textPort)
	if err := m.starHook(HOOK_TEXT); err != nil {
		logger.Errorf("start text hook error %s", err.Error())
		return err
	}
	return m.subscribe(m.textPort, HOOK_TEXT, func(data []byte) {
		logger.Debugf("text msg %s", string(data))
		m.notify(data)
	})
}

func (m *MessageRead) starHook(htype int) error {
	var (
		rs  *api.CommonResult
		err error
	)
	switch htype {
	case HOOK_TEXT:
		rs, err = m.server.RequestMsgStartHook(&api.MsgStartHookData{
			Port: m.textPort,
		})
	default:
		return errors.New("unknow hook")
	}

	if err != nil {
		return err
	}

	if rs.Result == "OK" {
		return nil
	}

	return errors.New(rs.Result)
}

func (m *MessageRead) stopHook(htype int) error {
	var (
		rs  *api.CommonResult
		err error
	)
	switch htype {
	case HOOK_TEXT:
		rs, err = m.server.RequestMsgStopHook()
	default:
		return errors.New("unknow hook")
	}

	if err != nil {
		return err
	}

	if rs.Result == "OK" {
		return nil
	}

	return errors.New(rs.Result)
}

// 消息通知
func (m *MessageRead) notify(data []byte) {
	if len(m.notifyURL) == 0 || len(data) == 0 {
		return
	}
	for _, v := range m.notifyURL {
		resp, err := http.Post(v, "application/json", bytes.NewReader(data))
		if err != nil {
			logger.Errorf("%s 消息通知失败 %s", v, err.Error())
			continue
		}
		if resp.StatusCode != http.StatusOK {
			logger.Errorf("%s 消息通知失败 status code not 200, get %d", v, resp.StatusCode)
			continue
		}
	}
}

func (m *MessageRead) subscribe(port, hookType int, call func(data []byte)) error {
	address := fmt.Sprintf("127.0.0.1:%d", port)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	go func() {
		defer func() {
			listen.Close()
			m.starHook(hookType)
		}()
		for {
			select {
			case <-m.done:
				return
			default:
				conn, err := listen.Accept()
				if err != nil {
					log.Fatal(err)
				}
				go m.handler(conn, call)
			}
		}
	}()
	return nil
}

func (m *MessageRead) handler(conn net.Conn, call func(data []byte)) {
	defer conn.Close()
	var (
		rs  []byte
		buf = make([]byte, 2048)
		r   = bufio.NewReader(conn)
	)

	//微信程序不会主动关闭，手动设置200毫秒读取时间
	conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	for {
		n, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}
		rs = append(rs, buf[:n]...)
	}
	call(rs)
}
