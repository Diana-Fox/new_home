package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"new_home/webook/gin/internal/web"
	"time"
)

func InitGin(mdls []gin.HandlerFunc, hdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	hdl.RegisterRoutes(server)
	return server
}
func InitMiddlewares(redisclient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			//允许的地址来源，*为所有，一般就公司仅有的域名就行了
			AllowOrigins: []string{"https://foo.com"},
			//允许请求的方法
			//AllowMethods: []string{"PUT", "PATCH"},
			//允许带上的请求头
			AllowHeaders: []string{"Content-Type", "Authorization", "x-jwt-token"},
			//允许带上的响应头
			ExposeHeaders: []string{"Content-Length"},
			//是否允许携带cookie
			AllowCredentials: true,
			//请求来源比较复杂的时候，用这个方法来判断是否允许
			//AllowOriginFunc: func(origin string) bool {
			//	return origin == "https://github.com"
			// //比如是否包含公司的域名
			//},
			MaxAge: 12 * time.Hour,
		}),
		//ratelimit.NewBuilder(redisclient, time.Second, 100).Build(),
		//middleware.NewLoginMiddlewareJWTBuilder().
		//	IgnorePaths("/users/loginJWT").
		//	IgnorePaths("/users/login_sms").
		//	IgnorePaths("/users/login_sms/code/send").
		//	IgnorePaths("/users/signup").Build(),
	}
}
