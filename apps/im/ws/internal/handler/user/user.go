package user

import (
	"easy-chat/apps/im/ws/internal/svc"
	websocketx "easy-chat/apps/im/ws/websocket"

	"github.com/gorilla/websocket"
)

func OnLine(svc *svc.ServiceContext) websocketx.HandlerFunc {
	return func(srv *websocketx.Server, conn *websocket.Conn, msg *websocketx.Message) {
		uids := srv.GetUsers()
		u := srv.GetUsers(conn)
		if len(u) == 0 {
			srv.Error("there is no user")
			return
		}
		err := srv.Send(websocketx.NewMessage(u[0], uids), conn)
		srv.Info("err ", err)
	}
}
