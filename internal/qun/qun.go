package qun

import (
	"wxbot/pkg/api"
)

func New() *Qun {
	return &Qun{}
}

type Qun struct {
	ch map[string]chan *api.TextMessage
}

func (q *Qun) Text(message *api.TextMessage) {

}
