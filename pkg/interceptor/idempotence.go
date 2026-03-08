package interceptor

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"easy-chat/pkg/xerr"
)

type Idempotent interface {
	Identify(ctx context.Context, method string) string
	IsIdempotentMethod(fullMethod string) bool
	TryAcquire(ctx context.Context, id string) (resp interface{}, isAcquire bool)
	SaveResp(ctx context.Context, id string, resp interface{}, respErr error) error
}

var (
	TKey = "easy-chat-idempotence-task-id"
	DKey = "easy-chat-idempotence-dispatch-key"
)

func ContextWithVal(ctx context.Context) context.Context {
	return context.WithValue(ctx, TKey, utils.NewUuid())
}

func NewIdempotenceClient(idempotent Idempotent) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		identify := idempotent.Identify(ctx, method)

		ctx = metadata.NewOutgoingContext(ctx, map[string][]string{
			DKey: {identify},
		})

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func NewIdempotenceServer(idempotent Idempotent) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		identify := metadata.ValueFromIncomingContext(ctx, DKey)
		if len(identify) == 0 || !idempotent.IsIdempotentMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		fmt.Println("----", "请求进入 幂等性处理 ", identify)

		r, isAcquire := idempotent.TryAcquire(ctx, identify[0])
		if isAcquire {
			resp, err = handler(ctx, req)
			fmt.Println("---- 执行任务", identify)

			if err := idempotent.SaveResp(ctx, identify[0], resp, err); err != nil {
				return resp, err
			}

			return resp, err
		}

		fmt.Println("----- 任务在执行", identify)

		if r != nil {
			fmt.Println("--- 任务已经执行完了 ", identify)
			return r, nil
		}

		return nil, errors.WithStack(xerr.New(int(codes.DeadlineExceeded), fmt.Sprintf("存在其他任务在执行 id %v", identify[0])))
	}
}

var (
	DefaultIdempotent       = new(defaultIdempotent)
	DefaultIdempotentClient = NewIdempotenceClient(DefaultIdempotent)
)

type defaultIdempotent struct {
	*redis.Redis
	*collection.Cache
	method map[string]bool
}

func NewDefaultIdempotent(c redis.RedisConf) Idempotent {
	cache, err := collection.NewCache(60 * 60)
	if err != nil {
		panic(err)
	}

	return &defaultIdempotent{
		Redis: redis.MustNewRedis(c),
		Cache: cache,
		method: map[string]bool{
			"/social.social/GroupCreate": true,
		},
	}
}

func (d *defaultIdempotent) Identify(ctx context.Context, method string) string {
	id := ctx.Value(TKey)
	rpcId := fmt.Sprintf("%v.%s", id, method)
	return rpcId
}

func (d *defaultIdempotent) IsIdempotentMethod(fullMethod string) bool {
	return d.method[fullMethod]
}

func (d *defaultIdempotent) TryAcquire(ctx context.Context, id string) (resp interface{}, isAcquire bool) {
	retry, err := d.SetnxEx(id, "1", 60*60)
	if err != nil {
		return nil, false
	}

	if retry {
		return nil, true
	}

	resp, _ = d.Cache.Get(id)
	return resp, false
}

func (d *defaultIdempotent) SaveResp(ctx context.Context, id string, resp interface{}, respErr error) error {
	d.Cache.Set(id, resp)
	return nil
}
