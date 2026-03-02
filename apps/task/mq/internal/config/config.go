package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
	service.ServiceConf

	ListenOn string

	Mysql struct {
		DataSource string
	}

	Cache cache.CacheConf

	MsgChatTransfer kq.KqConf

	Redisx redis.RedisConf
	Mongo  struct {
		Url string
		Db  string
	}

	Ws struct {
		Host string
	}
}
