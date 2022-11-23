package internal

import (
	"strings"
	"wxbot/pkg/api"
	"wxbot/pkg/logger"
)

type Contact struct {
	api *api.Api
}

func NewContact(api *api.Api) *Contact {
	return &Contact{
		api: api,
	}
}

type AbnormalStatusFriend struct {
	api.ContractInfo

	Status int
	Remark string
}

var FriendStatusDesc = map[int]string{
	0x0:  "Unknown",
	0xB0: "被删除",
	0xB1: "是好友",
	0xB2: "已拉黑",
	0xB5: "被拉黑",
}

// AbnormalStatusFriend 异常状态好友
func (c *Contact) AbnormalStatusFriend() ([]*AbnormalStatusFriend, error) {
	list, err := c.Friends()
	if err != nil {
		return nil, err
	}

	var rs []*AbnormalStatusFriend
	for _, v := range list {
		statusRs, err := c.api.RequestContactCheckStatus(&api.ContactCheckStatusData{
			Wxid: v.Wxid,
		})
		if err != nil {
			logger.Warnf("get friend status error %s, wxid %s %s", err.Error(), v.Wxid, v.WxNumber)
			continue
		}
		if statusRs.Status == 0xB1 {
			continue
		}

		rs = append(rs, &AbnormalStatusFriend{
			ContractInfo: *v,
			Status:       statusRs.Status,
			Remark:       FriendStatusDesc[statusRs.Status],
		})
	}
	return rs, err
}

// Friends 好友列表
func (c *Contact) Friends() ([]*api.ContractInfo, error) {
	list, err := c.api.RequestContactGetList()
	if err != nil {
		return nil, err
	}
	var rs []*api.ContractInfo
	for _, v := range list.Data {
		if v.WxType != 3 {
			continue
		}

		//公众号等
		if strings.HasPrefix(v.Wxid, "gh_") {
			continue
		}
		if v.WxVerifyFlag > 0 {
			continue
		}

		//一些特殊的微信号
		specialWxid := []string{
			"floatbottle",
			"fmessage",
			"qqmail",
			"filehelpe",
		}
		flag := false
		for _, wxid := range specialWxid {
			if v.Wxid == wxid {
				flag = true
			}
		}
		if flag {
			continue
		}

		tmp := v
		rs = append(rs, &tmp)
	}

	return rs, nil
}
