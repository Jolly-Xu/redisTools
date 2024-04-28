package main

import "redisTools/redis"

func main() {
	r := redis.CreateRedis(redis.FastConfig("10.32.15.225", 6379, "xujialin"))
	print(r.GetToString("mykey"))

}
