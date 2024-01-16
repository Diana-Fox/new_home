package web

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"new_home/webook/gin/internal/domain"
	"new_home/webook/gin/internal/service"
	svcmocks "new_home/webook/gin/internal/service/mock"
	"testing"
)

func TestUserHandler_Login(t *testing.T) {
	testCases := []struct {
		name string
		//email string
		//password string
		input    string
		mock     func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		wantCode int
		wantBody string
	}{
		{
			name: "成功的测试",
			input: `
					{
				   "email":"123@qq.com",
				   "password":"123456AAa!"
				}
				`,
			//email: "123@qq.com",
			//password: "123456AAa!",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				usersvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return usersvc, codeSvc
			},
			wantCode: 200,
			wantBody: "注册成功",
		},
	}
	//

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//使用Mock模拟
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			service := gin.Default()           //启动个服务器
			h := NewUserHandler(tc.mock(ctrl)) //获得模拟组件
			h.RegisterRoutes(service)          //注册路由
			resp := httptest.NewRecorder()     //存储网络请求的返回值
			t.Log(resp)
			req, err := http.NewRequest(http.MethodPost, "/users/login",
				bytes.NewBuffer([]byte(tc.input))) //最后一个参数是不是可以回头做成直接json序列化的串
			t.Log(err)
			req.Header.Set("Content-Type", "application/json")
			//http请求gin框架的入口
			service.ServeHTTP(resp, req)
			assert.Equal(t, 200, resp.Code)
			assert.Equal(t, "登录成功", resp.Body.String())
		})
	}
}
func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t) //创建控制器
	defer ctrl.Finish()             //调用结束的时候判断是否符合预期
	usersvc := svcmocks.NewMockUserService(ctrl)
	//调用Sinup方法，预期的参数，参数和返回值
	usersvc.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.User{}, errors.New("模拟的异常返回值"))
	//调用和返回值
	//login, err := usersvc.Login(context.Background(), "email", "pwd")
	//if err != nil {
	//	return
	//}
}
