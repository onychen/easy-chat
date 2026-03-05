package handler

import (
	"easy-chat/apps/im/api/internal/logic"
	"easy-chat/apps/im/api/internal/svc"
	"easy-chat/apps/im/api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getChatLogReadRecordsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetChatLogReadRecordsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetChatLogReadRecordsLogic(r.Context(), svcCtx)
		resp, err := l.GetChatLogReadRecords(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
