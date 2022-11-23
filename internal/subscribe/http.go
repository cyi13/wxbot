package subscribe

import (
	"bytes"
	"encoding/json"
	"net/http"
	"wxbot/pkg/api"
	"wxbot/pkg/logger"
)

type Http struct {
	url []string
}

func NewHttp(url []string) *Http {
	return &Http{
		url: url,
	}
}

func (h *Http) Text(message *api.TextMessage) {
	if message == nil {
		return
	}
	data, _ := json.Marshal(message)

	for _, v := range h.url {
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
