package memory

import (
	"context"
	"fmt"
)

type Service struct {
}

// 发送
func (s Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	fmt.Println(args)
	return nil
}
