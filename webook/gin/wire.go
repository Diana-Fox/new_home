//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"new_home/webook/gin/internal/repository"
	"new_home/webook/gin/internal/repository/cache"
	"new_home/webook/gin/internal/repository/dao"
	"new_home/webook/gin/internal/service"
	"new_home/webook/gin/internal/service/sms/memory"
	"new_home/webook/gin/internal/web"
	"new_home/webook/gin/ioc"
)

func InitWebService() *gin.Engine {
	wire.Build(
		ioc.InitDB,
		ioc.InitRedis,
		dao.NewUserDAO, cache.NewUserCache,
		repository.NewUserRepository,
		service.NewUserService,
		cache.NewCodeCache,
		repository.NewCodeRepository,
		memory.NewService,
		service.NewCodeService,
		web.NewUserHandler,
		//gin.Default,
		ioc.InitMiddlewares, ioc.InitGin,
	)
	return new(gin.Engine)
}
