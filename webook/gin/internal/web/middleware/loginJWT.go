package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"new_home/webook/gin/internal/web"
	"time"
)

type LoginMiddlewareBuilderJWT struct {
	paths []string
}

func NewLoginMiddlewareJWTBuilder() *LoginMiddlewareBuilderJWT {
	return &LoginMiddlewareBuilderJWT{}
}
func (l *LoginMiddlewareBuilderJWT) IgnorePaths(path string) *LoginMiddlewareBuilderJWT {
	l.paths = append(l.paths, path)
	return l
}
func (l *LoginMiddlewareBuilderJWT) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		token := ctx.GetHeader("x-jwt-token")
		if token == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//parse, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		//	return []byte("要有密钥"), nil
		//})
		clamis := &web.UserClaims{}
		parse, err := jwt.ParseWithClaims(token, clamis, func(token *jwt.Token) (interface{}, error) {
			return []byte("要有密钥"), nil
		})
		if err != nil {
			return
		}
		if err != nil || parse.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		now := time.Now()
		if clamis.ExpiresAt.Sub(now) < time.Second*50 {
			//10S一续约
			//续约，生成新token
			clamis.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			signedString, err := parse.SignedString([]byte("要有密钥"))
			if err != nil {
				ctx.String(http.StatusOK, "加密异常")
				return
			}
			ctx.Header("x-jwt-token", signedString)
		}
		ctx.Set("claims", clamis)
	}
}
