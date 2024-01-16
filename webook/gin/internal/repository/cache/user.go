package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"new_home/webook/gin/internal/domain"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, user domain.User) error
}

func NewUserCache(cmd redis.Cmdable) UserCache {
	return &userCache{
		cmd:        cmd,
		expiration: 10,
	}
}

type userCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (cache *userCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	val, err := cache.cmd.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err
}
func (cache *userCache) Set(ctx context.Context, user domain.User) error {
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return cache.cmd.Set(ctx, cache.key(user.Id), val, cache.expiration).Err()
}
func (cache *userCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
