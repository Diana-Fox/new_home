package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany       = errors.New("发送验证码太频繁")
	ErrCodeVerifyTooManyTime = errors.New("验证次数过多")
	ErrUnknownForCode        = errors.New("未知错误")
)

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}
type codeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) CodeCache {
	return &codeCache{
		client: client,
	}
}
func (c *codeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		//毫无问题
		return nil
	case -2:
		return ErrCodeSendTooMany
	default:
		//系统错误
		return errors.New("系统错误")
	}
}

func (c *codeCache) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		//毫无问题
		return true, nil
	case -1:
		return false, ErrCodeVerifyTooManyTime
	case -2:
		return false, nil
	default:
		//系统错误
		return false, ErrUnknownForCode
	}
}
func (c *codeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
