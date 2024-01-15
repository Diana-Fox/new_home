package sms

import "context"

// 抽象短信接口
type Service interface {
	Send(ctx context.Context, tpl string, args []string, numbers ...string) error
}
