package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
	"new_home/webook/gin/internal/domain"
	"new_home/webook/gin/internal/service"
)

// 所有跟用户相关的路由
type UserHandler struct {
	svc           service.UserService
	emailExp      *regexp.Regexp
	passwordExp   *regexp.Regexp
	brotherDayExp *regexp.Regexp
}

func NewUserHandler(svc service.UserService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
		brotherDayPattern    = `(?!0000)[0-9]{4}-((0[1-9]|1[0-2])-(0[1-9]|1[0-9]|2[0-8])|(0[13-9]|1[0-2])-(29|30)|(0[13578]|1[02])-31)`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	brotherDayExp := regexp.MustCompile(brotherDayPattern, regexp.None)
	return &UserHandler{
		svc:           svc,
		emailExp:      emailExp,
		passwordExp:   passwordExp,
		brotherDayExp: brotherDayExp,
	}
}
func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ur := server.Group("/users")
	ur.GET("/ping", func(context *gin.Context) {
		context.String(http.StatusOK, "pong")
		return
	})
	ur.POST("/signup", u.SignUp)
	ur.POST("/login", u.Login)
	ur.POST("/edit", u.Edit)
	ur.GET("/profile", u.Profile)
}
func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}

	var req SignUpReq
	// Bind 方法会根据 Content-Type 来解析你的数据到 req 里面
	// 解析错了，就会直接写回一个 400 的错误
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "你的邮箱格式不对")
		return
	}
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		// 记录日志
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含数字、特殊字符")
		return
	}
	fmt.Printf("%v", req)
	// 这边就是数据库操作
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		PassWord: req.Password,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
}
func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码错误")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	//取session进行设置
	fmt.Printf(user.Email)
	ctx.String(http.StatusOK, "登陆成功")
}
func (u *UserHandler) Edit(ctx *gin.Context) {
	type UserEditReq struct {
		Id         int64  `json:"id"`
		NickName   string `json:"nickName"`
		BrotherDay string `json:"brotherDay"`
		Biography  string `json:"biography"`
	}
	var req UserEditReq
	if err := ctx.BindJSON(&req); err != nil {
		return
	}
	ok, err := u.brotherDayExp.MatchString(req.BrotherDay)
	if err != nil {
		// 记录日志
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "生日格式错误")
		return
	}
	if len(req.NickName) > 16 {
		ctx.String(http.StatusOK, "昵称过长，请修改")
		return
	}
	if len(req.Biography) > 127 {
		ctx.String(http.StatusOK, "简介超长，请修改")
		return
	}
	err = u.svc.Edit(ctx, req.Id, req.NickName, req.BrotherDay, req.Biography)
	if err == service.ErrUserNotFound {
		ctx.String(http.StatusOK, "用户名不存在")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.String(http.StatusOK, "修改成功")
}
func (u *UserHandler) Profile(ctx *gin.Context) {
	type Id struct {
		Id int64 `json:"id"`
	}
	var req Id
	if err := ctx.BindJSON(&req); err != nil {
		return
	}
	user, err := u.svc.Profile(ctx, req.Id)
	if err == service.ErrUserNotFound {
		ctx.String(http.StatusOK, "用户不存在")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "获取失败")
	}
	ctx.JSON(http.StatusOK, user)
}
