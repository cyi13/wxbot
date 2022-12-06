package demo

import (
	"wxbot/internal/global"
	"wxbot/pkg/api"
)

type Demo struct {
	api api.Api

	wxid string
}

func New(wxid string) *Demo {
	d := &Demo{
		wxid: wxid,
	}

	global.Api.RequestMsgSendText(&api.MsgSendTextData{
		Wxid: wxid,
		Msg:  "启动: " + d.Name(),
	})

	return d
}

func (d *Demo) Name() string {
	return "测试模块"
}

func (d *Demo) Text(message *api.TextMessage) error {
	return nil
}
