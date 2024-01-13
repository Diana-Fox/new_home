//go:build dev

// 没使用k8s标签
package config

var Config = config{
	DB: DBConfig{
		DNS: "localhost:13306", //本地连接
	},
	Redis: RedisConfig{
		Addr: "localhost:16379",
	},
}
