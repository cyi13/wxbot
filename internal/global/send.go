package global

import (
	"sync"
	"wxbot/pkg/api"
	"wxbot/pkg/logger"
)

var Api *api.Api
var once sync.Once

var textChan = make(chan api.MsgSendTextData, 100)

func SendText(wxid, msg string) {
	once.Do(func() {
		go sendText()
	})
	textChan <- api.MsgSendTextData{
		Wxid: wxid,
		Msg:  msg,
	}
}

func sendText() {
	for v := range textChan {
		tmp := v
		_, err := Api.RequestMsgSendText(&tmp)
		if err != nil {
			logger.Warnf("global send message error %s, data %+v", err.Error(), tmp)
		}
	}
}
