package main

import (
	"redisTools/redis"
)

type userinfo struct {
	Username string
	Password string
}

func test(args ...interface{}) {
	println("后置处理器")
}

func main() {
	r := redis.CreateRedis(redis.FastConfig("10.32.15.225", 6379, "xujialin"))
	//println(r.GetToString("mykey"))
	r.AddPreInterceptor(func(args ...interface{}) {
		println("前置处理器")
	})
	r.AddPostInterceptor(test)
	println(r.GetToString("111"))
	//println(r.GetToString("test"))
	//r.Set("jolly", userinfo{
	//	Username: "hehhh",
	//	Password: "1111",
	//})
}
