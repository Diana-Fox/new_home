package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"new_home/webook/gin/internal/repository"
	"new_home/webook/gin/internal/repository/dao"
	"new_home/webook/gin/internal/service"
	"new_home/webook/gin/internal/web"
	"time"
)

func main() {
	r := gin.Default() //用engine来监听端口
	//通过middleware实现跨域
	//https://github.com/gin-contrib/cors
	r.Use(cors.New(cors.Config{
		//允许的地址来源，*为所有，一般就公司仅有的域名就行了
		AllowOrigins: []string{"https://foo.com"},
		//允许请求的方法
		//AllowMethods: []string{"PUT", "PATCH"},
		//允许带上的请求头
		AllowHeaders: []string{"Content-Type", "Authorization"},
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
	}))
	user := initUser(initDB())
	user.RegisterRoutes(r)
	r.Run(":18080") // 监听并在 0.0.0.0:8080 上启动服务
}
func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}
func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13306)/webook"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
