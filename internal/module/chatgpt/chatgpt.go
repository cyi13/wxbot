package chatgpt

import (
	"fmt"
	"wxbot/internal/global"
	"wxbot/pkg/api"
	"wxbot/pkg/logger"

	"github.com/m1guelpf/chatgpt-telegram/src/chatgpt"
	"github.com/m1guelpf/chatgpt-telegram/src/config"
)

type Service struct {
	chat         *chatgpt.ChatGPT
	messageLimit int //消息限制，超过的返回错误

	message chan *api.TextMessage
	room    map[string]Conversation
}

type Conversation struct {
	ConversationID string
	LastMessageID  string
}

type DefaultHandler *Service

func New(token string) *Service {
	chat := chatgpt.Init(config.Config{
		OpenAISession: token,
	})

	s := &Service{
		chat:         &chat,
		messageLimit: 30,
		//排队处理
		message: make(chan *api.TextMessage, 50),
		room:    make(map[string]Conversation),
	}
	go s.call()
	return s
}

func (s *Service) Name() string {
	return "聊天机器人"
}

func (s *Service) Text(message *api.TextMessage) error {
	if len(s.message) > s.messageLimit {
		logger.Infof("chatgpt message len %d reach limit", len(s.message))
		global.SendText(message.Sender, "抱歉，消息堆积过多，请稍后再发送")
		return nil
	}
	//处理下数据过多问题
	s.message <- message
	logger.Infof("chatgpt message len %d", len(s.message))
	if len(s.message) > 2 {
		global.SendText(message.Sender, fmt.Sprintf("当前消息排队位置 %d", len(s.message)))
	}
	return nil
}

func (s *Service) call() {
	for v := range s.message {
		tmp := v
		s.text(tmp)
	}
}

func (s *Service) text(message *api.TextMessage) error {
	logger.Infof("chatgpt %+v", message)
	rs, err := s.callChatgpt(message.Sender, message.Message)
	if err != nil {
		global.SendText(message.Sender, "抱歉，机器人服务请求失败，请联系管理员")

		return err
	}
	global.SendText(message.Sender, rs)
	return nil
}

func (s *Service) Add(wxid string) {
	if _, ok := s.room[wxid]; ok {
		return
	}
	s.room[wxid] = Conversation{}
}

// 请求chatgpt接口
func (s *Service) callChatgpt(wxid, message string) (string, error) {
	before := s.room[wxid]
	resp, err := s.chat.SendMessage(message, before.ConversationID, before.LastMessageID)
	if err != nil {
		logger.Warnf("call chatgpt error %s ", err.Error())
		return "", err
	}

	var (
		res            string
		conversationid string
		messageid      string
	)

	for v := range resp {
		logger.Debugf(v.Message)
		res = v.Message
		conversationid = v.ConversationId
		messageid = v.MessageId
	}

	s.room[wxid] = Conversation{
		ConversationID: conversationid,
		LastMessageID:  messageid,
	}
	return res, nil
}
