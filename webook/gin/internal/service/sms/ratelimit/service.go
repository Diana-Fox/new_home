package ratelimit

import (
	"context"
	"fmt"
	"new_home/webook/gin/internal/service/sms"
	"new_home/webook/gin/pkg/ginx/ratelimit"
)

type Service struct {
	svc sms.Service
	r   ratelimit.Limiter
}

func NewService(svc Service, r ratelimit.Limiter) sms.Service {
	return &Service{
		svc: &svc,
		r:   r,
	}
}
func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	limit, err := s.r.Limit(ctx, "sms:tencent")
	if err != nil {
		return err
	}
	if limit {
		return fmt.Errorf("触发了限流")
	}
	err = s.svc.Send(ctx, tpl, args, numbers...)
	if err != nil {
		return err
	}
	return nil
}
