package module

import "wxbot/pkg/api"

type Module interface {
	Name() string
	Text(message *api.TextMessage) error
}
