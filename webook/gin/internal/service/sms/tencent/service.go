package tencent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	smsClient "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"new_home/webook/gin/internal/service/sms"
)

type service struct {
	appId    *string
	signName *string
	client   *smsClient.Client
}

func NewService(client *smsClient.Client, signName *string, appId *string) sms.Service {
	return &service{
		client:   client,
		signName: signName,
		appId:    appId,
	}
}
func (s *service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	req := smsClient.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = &tpl
	req.PhoneNumberSet = slice.Map[string, *string](numbers, func(idx int, src string) *string {
		return &src
	})
	req.TemplateParamSet = slice.Map[string, *string](args, func(idx int, src string) *string {
		return &src
	})
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	//failPhoneSet:=make([]string,0)
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			//	failPhoneSet = append(failPhoneSet, *status.PhoneNumber)
			//实际业务上，这里应该是个群发接口，返回失败的号码和内容给客户端，这里因为接口定义原因直接返回
			return fmt.Errorf("发送短信失败%s,%s", status.Code, status.Message)
		}
	}
	return nil
}
