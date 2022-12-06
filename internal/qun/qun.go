// Pacage qun
// 群管理功能处理逻辑
package qun

import (
	"regexp"
	"strings"
	"wxbot/internal/config"
	"wxbot/internal/global"
	"wxbot/internal/module/chatgpt"
	"wxbot/pkg/api"
	"wxbot/pkg/common"

	"wxbot/internal/module"
)

var DefaultHandler *Qun

func New(conf config.QunManager) *Qun {
	return &Qun{
		wxid: make(map[string][]module.Module),
	}
}

type Qun struct {
	self string
	conf config.QunManager

	msg  chan *api.TextMessage
	wxid map[string][]module.Module
}

// 对指定的群聊启用群管理功能
func (q *Qun) Add(wxid string) {
	q.wxid[wxid] = []module.Module{}
}

// 默认启用的模块
func (q *Qun) defaultModules(wxid string) []module.Module {
	return []module.Module{
		chatgpt.New(global.Config.Module.ChatGPT.AiSession),
	}
}

// 移除指定群的群管理功能
func (q *Qun) Remove(wxid string) {
	delete(q.wxid, wxid)
}

// 判断是否是群聊的微信id格式
func (q *Qun) IsQunWxid(wxid string) bool {
	return strings.HasSuffix(wxid, "@chatroom")
}

func (q *Qun) Text(message *api.TextMessage) {
	if message == nil {
		return
	}
	//默认的消息处理模块，主要用来启用/退出群聊管理等其他功能

	sender := message.Sender
	//先过滤非群聊的消息
	if !q.IsQunWxid(message.Sender) {
		return
	}

	// 暂时只处理@自己的消息
	atuser := strings.Split(message.Extra.Atuserlist, ",")
	if len(atuser) == 0 {
		return
	}
	if !common.StrInSlice(atuser, global.Self.WxID) {
		return
	}

	// if _, ok := q.wxid[sender]; !ok {
	// 	return
	// }

	//去除消息的@头部
	message.Message = q.filterMessage(message.Message)

	//消息逐个模块处理，返回nil表示由此模块处理
	for _, module := range q.defaultModules(sender) {
		if err := module.Text(message); err == nil {
			break
		}
	}
}

// todo 报错
func (q *Qun) handler(message *api.TextMessage) {

	text := q.filterMessage(message.Message)

	switch {
	//启用群管理功能
	case common.StrInSlice(q.conf.EnableWord, text):
		q.Add(message.Sender)
	//退出群管理功能
	case common.StrInSlice(q.conf.DisbaleWord, text):
		q.Remove(message.Message)
	}
}

// 过滤消息 把@某人 等数据过滤掉
func (q *Qun) filterMessage(s string) string {
	//TODO \u2005 1\4空格优化匹配
	reg := regexp.MustCompile(`(@\S+? )`)
	str := reg.ReplaceAllString(s, "")
	str = strings.TrimSpace(str)
	return str
}
