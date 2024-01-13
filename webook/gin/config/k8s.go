//go:build k8s

// 使用k8s标签
package config

// 用k8s的域名来访问
var Config = config{
	DB: DBConfig{
		DNS: "root:root@tcp(webook-live-mysql:11309)/webook", //本地连接
	},
	Redis: RedisConfig{
		Addr: "webook-live-redis:11479",
	},
}
