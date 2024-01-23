package memory

import (
	"context"
	"fmt"
	"new_home/webook/gin/internal/service/sms"
)

type Service struct {
}

func NewService() sms.Service {
	return &Service{}
}

// 发送
func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	fmt.Println(args)
	return nil
}
