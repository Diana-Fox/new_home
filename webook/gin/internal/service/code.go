package service

import (
	"context"
	"fmt"
	"math/rand"
	"new_home/webook/gin/internal/repository"
	"new_home/webook/gin/internal/service/sms"
)

// 写个假的模板id
const codeTplId = "2345678"

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}
type codeService struct {
	repo repository.CodeRepository
	sms  sms.Service
}

func NewCodeService(repo repository.CodeRepository, sms sms.Service) CodeService {
	return &codeService{repo, sms}
}
func (svc *codeService) Send(ctx context.Context, biz string, phone string) error {
	code := svc.generateCode()
	//先存储
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	//去发送
	err = svc.sms.Send(ctx, codeTplId, []string{code}, phone)
	return err
}
func (svc *codeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)

}
func (svc *codeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%6d", num)
}
