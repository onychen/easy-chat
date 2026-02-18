package logic

import (
	"context"
	"database/sql"
	"time"

	"easy-chat/apps/social/rpc/internal/svc"
	"easy-chat/apps/social/rpc/social"
	"easy-chat/apps/social/socialmodels"
	"easy-chat/pkg/xerr"

	"easy-chat/pkg/constants"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type FriendPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 好友业务：请求好友、通过或拒绝申请、好友列表
func (l *FriendPutInLogic) FriendPutIn(in *social.FriendPutInReq) (*social.FriendPutInResp, error) {
	// todo: add your logic here and delete this line

	// 申请人是否与目标是好友关系
	friends, err := l.svcCtx.FriendsModel.FindByUidAndFid(l.ctx, in.UserId, in.ReqUid)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friends by uid and fid err %v req %v ", err, in)
	}
	if friends != nil {
		return &social.FriendPutInResp{}, err
	}

	// 是否已经有过申请，申请是不成功，没有完成
	friendReqs, err := l.svcCtx.FriendRequestsModel.FindByReqUidAndUserId(l.ctx, in.ReqUid, in.UserId)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friendsRequest by rid and uid err %v req %v ", err, in)
	}
	if friendReqs != nil {
		if constants.HandlerResult(friendReqs.HandleResult.Int64) != constants.RefuseHandlerResult {
			return &social.FriendPutInResp{}, err
		}

		friendReqs.ReqMsg = sql.NullString{
			Valid:  true,
			String: in.ReqMsg,
		}
		friendReqs.ReqTime = time.Unix(in.ReqTime, 0)
		friendReqs.HandleResult = sql.NullInt64{
			Int64: int64(constants.NoHandlerResult),
			Valid: true,
		}
		friendReqs.HandleMsg = sql.NullString{}
		friendReqs.HandledAt = sql.NullTime{}

		err = l.svcCtx.FriendRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			if err := l.svcCtx.FriendRequestsModel.Update(ctx, session, friendReqs); err != nil {
				return errors.Wrapf(xerr.NewDBErr(), "update friendRequest err %v req %v ", err, friendReqs)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		return &social.FriendPutInResp{}, nil
	}

	// 创建申请记录
	_, err = l.svcCtx.FriendRequestsModel.Insert(l.ctx, &socialmodels.FriendRequests{
		UserId: in.UserId,
		ReqUid: in.ReqUid,
		ReqMsg: sql.NullString{
			Valid:  true,
			String: in.ReqMsg,
		},
		ReqTime: time.Unix(in.ReqTime, 0),
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandlerResult),
			Valid: true,
		},
	})

	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert friendRequest err %v req %v ", err, in)
	}

	return &social.FriendPutInResp{}, nil
}
