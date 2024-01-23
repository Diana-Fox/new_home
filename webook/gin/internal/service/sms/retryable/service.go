package retryable

import (
	"context"
	"errors"
	"new_home/webook/gin/internal/service/sms"
)

type Service struct {
	svc      sms.Service
	retryCnt int //计数会有并发问题，不能是单例
}

func NewService(svc sms.Service, retryCnt int) sms.Service {
	return &Service{svc: svc, retryCnt: retryCnt}
}
func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	for i := 0; i < s.retryCnt; i++ {
		err := s.Send(ctx, tpl, args, numbers...)
		if err == nil {
			return nil
		}
	}
	return errors.New("发送失败【考虑转为异步】")
}
