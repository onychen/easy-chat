package main

import (
	"flag"
	"fmt"
	"sync"

	"easy-chat/apps/task/mq/internal/config"
	"easy-chat/apps/task/mq/internal/handler"
	"easy-chat/apps/task/mq/internal/svc"
	"easy-chat/pkg/configserver"

	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/local/mq.yaml", "the config file")
var wg sync.WaitGroup

func main() {
	flag.Parse()

	var c config.Config
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.100.1:3379",
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2",
		Namespace:      "task",
		Configs:        "task-mq.yaml",
		ConfigFilePath: "./conf",
		LogLevel:       "DEBUG",
	})).MustLoad(&c, func(bytes []byte) error {
		var c config.Config
		configserver.LoadFromJsonBytes(bytes, &c)

		wg.Add(1)
		go func(c config.Config) {
			defer wg.Done()

			Run(c)
		}(c)
		return nil
	})
	if err != nil {
		panic(err)
	}

	wg.Add(1)
	go func(c config.Config) {
		defer wg.Done()

		Run(c)
	}(c)

	wg.Wait()
}

func Run(c config.Config) {
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	serviceGrop := service.NewServiceGroup()
	defer serviceGrop.Stop()

	svcCtx := svc.NewServiceContext(c)
	listen := handler.NewListen(svcCtx)
	for _, s := range listen.Services() {
		serviceGrop.Add(s)
	}
	fmt.Println("Starting mqueue server at...")

	defer serviceGrop.Stop()
	serviceGrop.Start()
}
