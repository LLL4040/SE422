package main

import (
	"fmt"
	"net/http"
	"os"

	config "./config"
	connection "./connection"
	processor "./processor"
)

func main() {
	configure, err := config.NewConfigure()
	if err != nil {
		fmt.Printf("[ERROR] Parse Configure File Error: %v\n", err)
		return
	}

	// 启动redis连接
	fmt.Println("[INFO] Starting Redis Connection...")
	redis, err := connection.NewRedisAdaptor(configure)
	if err != nil {
		fmt.Printf("[ERROR] Redis connection fail")
		return
	}

	// 初始化redis计数器
	err = redis.InitCountService()
	if err != nil {
		fmt.Printf("[ERROR] Redis key count fail")
		return
	}
	countfunction := processor.CreateCounter(redis)

	// 启动LRU缓存
	fmt.Println("[INFO] Starting LRU Cache System...")
	lru, err := connection.NewLRU(redis)
	if err != nil {
		fmt.Printf("[ERROR] LRU init fail...")
		return
	}

	// 初始化短链接服务
	fmt.Println("[INFO] Starting Service...")
	baseprocessor := &processor.BaseProcessor{redis, configure, lru, countfunction}

	original := &processor.OriginalProcessor{baseprocessor}
	short := &processor.ShortProcessor{baseprocessor}

	// 启动http handler
	router := &processor.Router{configure, map[int]processor.Processor{
		0: short,
		1: original,
	}}

	port, _ := configure.GetPort()
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("[INFO]Service Starting addr :%v,port :%v\n", addr, port)
	err = http.ListenAndServe(addr, router)
	if err != nil {
		//logger.Error("Server start fail: %v", err)
		os.Exit(1)
	}
}
