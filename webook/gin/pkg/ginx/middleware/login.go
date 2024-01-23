package middleware

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}
func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}
func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		sess.Set("userId", id)
		updateTime := sess.Get("updateTime")
		//
		now := time.Now().UnixMilli()
		if updateTime == nil {
			sess.Set("updateTime", now)
			sess.Options(sessions.Options{
				MaxAge: 60 * 30, //
			})
			sess.Save()
		}
		updateTimeVal, _ := updateTime.(int64)
		if now-updateTimeVal > 60*1000 {
			sess.Set("updateTime", now)
			sess.Options(sessions.Options{
				MaxAge: 60 * 30,
			})
			sess.Save()
		}
		fmt.Printf("test-%v\n", id)
	}
}
