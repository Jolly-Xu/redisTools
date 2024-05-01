package main

import (
	"encoding/hex"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"redisTools/redis"
)

type userinfo struct {
	Username string
	Password string
}

func test(redisContext *redis.InterceptorContext) {
	println("后置处理器开始")
	//args := redisContext.Cmd().Args()
	//i := args[0]
	res := redisContext.Cmd().Res()
	u := userinfo{}
	s := string(res)
	decodeString, err := hex.DecodeString(s)
	if err != nil {
		return
	}
	err = msgpack.Unmarshal(decodeString, &u)
	if err != nil {
		return
	}

}

func getWithStruct(vt interface{}, arr []byte) {
	err := msgpack.Unmarshal(arr, vt)
	if err != nil {
		fmt.Println("不能将[]byte转为", vt)
	}
	return
}

func main() {
	//cli := redis.CreateRedisCli(redis.FastConfig("10.32.2.37", 6379, "xujialin"))
	//println(cli.Get("test"))
	//ip := net.ParseIP("192.168.1.1")
	//cli.Set("ip2", ip)
	//cli.Set("userinfo", marshal)
	//cli.AddPostInterceptor(test)
	//println(cli.Get("userinfo"))
}
