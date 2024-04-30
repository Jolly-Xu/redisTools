package main

import "redisTools/redis"

type userinfo struct {
	Username string
	Password string
}

func test(redisContext *redis.InterceptorContext) {
	println("后置处理器开始")

}

func main() {
	cli := redis.CreateRedisCli(redis.FastConfig("10.32.2.37", 6379, "xujialin"))
	//println(cli.Get("test"))
	cli.Set("xujialin", "jolly")
}
