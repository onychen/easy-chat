package rpcserver

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	zerr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LogInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any,
	err error) {

	fmt.Printf("========== [RPC REQ] method: %s ==========\n", info.FullMethod)

	resp, err = handler(ctx, req)
	if err == nil {
		fmt.Printf("========== [RPC RESP] method: %s, success ==========\n", info.FullMethod)
		return resp, nil
	}

	fmt.Printf("========== [RPC ERR] method: %s, err: %v ==========\n", info.FullMethod, err)
	logx.WithContext(ctx).Errorf("[RPC ERR] method: %s, err: %v", info.FullMethod, err)

	causeErr := errors.Cause(err)
	if e, ok := causeErr.(*zerr.CodeMsg); ok {
		err = status.Error(codes.Code(e.Code), e.Msg)
	}

	return resp, err
}
