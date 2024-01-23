package failover

import (
	"context"
	"errors"
	"new_home/webook/gin/internal/service/sms"
	"sync/atomic"
)

type FailoverSMSService struct {
	svcs []sms.Service
	idx  uint64
}

func NewFailoverSMSService(svcs []sms.Service) sms.Service {
	return &FailoverSMSService{
		svcs: svcs,
	}
}
func (f *FailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tpl, args, numbers...)
		if err == nil { //发送成功
			return err
		}
		//做好监控
	}
	return errors.New("全部服务商都失败了")
}

func (f FailoverSMSService) SendV1(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.AddUint64(&f.idx, 1) //取下一个节点
	length := uint64(len(f.svcs))
	for i := idx; i < idx+length; i++ {
		svc := f.svcs[i%length]
		err := svc.Send(ctx, tpl, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			return err
		}
	}
	return errors.New("全部服务商都失败了")
}
