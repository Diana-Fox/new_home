package main

func main() {
	//r := gin.Default() //用engine来监听端口
	////通过middleware实现跨域
	////https://github.com/gin-contrib/cors
	//r.Use(cors.New(cors.Config{
	//	//允许的地址来源，*为所有，一般就公司仅有的域名就行了
	//	AllowOrigins: []string{"https://foo.com"},
	//	//允许请求的方法
	//	//AllowMethods: []string{"PUT", "PATCH"},
	//	//允许带上的请求头
	//	AllowHeaders: []string{"Content-Type", "Authorization", "x-jwt-token"},
	//	//允许带上的响应头
	//	ExposeHeaders: []string{"Content-Length"},
	//	//是否允许携带cookie
	//	AllowCredentials: true,
	//	//请求来源比较复杂的时候，用这个方法来判断是否允许
	//	//AllowOriginFunc: func(origin string) bool {
	//	//	return origin == "https://github.com"
	//	// //比如是否包含公司的域名
	//	//},
	//	MaxAge: 12 * time.Hour,
	//}))
	////store := cookie.NewStore([]byte("secret"))
	////接入
	//store, err := redis.NewStore(16, "tcp", config.Config.Redis.Addr, "", []byte("secret"))
	//if err != nil {
	//	panic(err)
	//}
	//redisclient := redisClient.NewClient(&redisClient.Options{
	//	Addr: config.Config.Redis.Addr,
	//})
	//r.Use(ratelimit.NewBuilder(redisclient, time.Second, 100).Build())
	//r.Use(sessions.Sessions("mysession", store))
	//r.Use(middleware.NewLoginMiddlewareJWTBuilder().
	//	IgnorePaths("/users/loginJWT").
	//	IgnorePaths("/users/login_sms").
	//	IgnorePaths("/users/login_sms/code/send").
	//	IgnorePaths("/users/signup").Build())
	////r.Use(middleware.NewLoginMiddlewareBuilder().
	////	IgnorePaths("/users/login").
	////	IgnorePaths("/users/signup").Build())
	//user := initUser(initDB(), redisclient)
	//user.RegisterRoutes(r)
	r := InitWebService()

	r.Run(":18080") // 监听并在 0.0.0.0:8080 上启动服务
}

//func initUser(db *gorm.DB, rdb redisClient.Cmdable) *web.UserHandler {
//	ud := dao.NewUserDAO(db)
//	uc := cache.NewUserCache(rdb, 10)
//	repo := repository.NewUserRepository(ud, uc)
//	svc := service.NewUserService(repo)
//	codeCache := cache.NewCodeCache(rdb)
//	codeRepo := repository.NewCodeRepository(codeCache)
//	me := memory.Service{}
//	codeSvc := service.NewCodeService(codeRepo, me)
//	u := web.NewUserHandler(svc, codeSvc)
//	return u
//}
//func initDB() *gorm.DB {
//	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13306)/webook"), &gorm.Config{})
//	if err != nil {
//		panic(err)
//	}
//	err = dao.InitTable(db)
//	if err != nil {
//		panic(err)
//	}
//	return db
//}
