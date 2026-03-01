package main

import (
	"easy-chat/apps/task/mq/internal/config"
	"easy-chat/apps/task/mq/internal/handler"
	"easy-chat/apps/task/mq/internal/svc"
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/local/mq.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config

	conf.MustLoad(*configFile, &c)
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
