package friend

import (
	"context"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/apps/user/rpc/userclient"
	"easy-chat/pkg/ctxdata"
	"easy-chat/pkg/xerr"

	"easy-chat/apps/social/api/internal/svc"
	"easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInLogic) FriendPutIn(req *types.FriendPutInReq) (resp *types.FriendPutInResp, err error) {
	// todo: add your logic here and delete this line

	uid := ctxdata.GetUId(l.ctx)

	userResp, err := l.svcCtx.User.FindUser(l.ctx, &userclient.FindUserReq{
		Ids: []string{req.UserId},
	})
	if err != nil {
		return nil, err
	}
	if len(userResp.User) == 0 {
		return nil, xerr.NewMsg("用户不存在")
	}

	_, err = l.svcCtx.Social.FriendPutIn(l.ctx, &socialclient.FriendPutInReq{
		UserId:  uid,
		ReqUid:  req.UserId,
		ReqMsg:  req.ReqMsg,
		ReqTime: req.ReqTime,
	})

	return
}
