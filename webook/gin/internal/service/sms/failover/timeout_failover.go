package failover

import (
	"context"
	"new_home/webook/gin/internal/service/sms"
	"sync/atomic"
)

type TimeoutFailoverSMSService struct {
	svcs []sms.Service
	idx  int32
	cnt  int32 //连续超时的个数
	//超过阈值就切换
	threshold int32 //阈值
}

func NewTimeoutFailoverSMSService(svcs []sms.Service, cnt int32) sms.Service {
	return &TimeoutFailoverSMSService{
		svcs: svcs,
		cnt:  cnt,
	}
}

func (t *TimeoutFailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt > t.threshold {
		//这里要切换
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			//切换成功
			atomic.StoreInt32(&t.cnt, 0) //重置阈值
		}
		//拿到下一个
		idx = atomic.LoadInt32(&t.cnt)
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tpl, args, numbers...)
	if err != nil {
		return err
	}
	switch err {
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
	case nil:
		//你的连续状态被打断了
		atomic.StoreInt32(&t.cnt, 0)
	default:
		return err
	}
	return err
}
