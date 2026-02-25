package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"sync"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type Server struct {
	sync.RWMutex

	authentication Authentication
	routes         map[string]HandlerFunc
	addr           string

	connToUser map[*websocket.Conn]string
	userToConn map[string]*websocket.Conn
	patten     string

	upgrader websocket.Upgrader
	logx.Logger
}

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)

	return &Server{
		addr:   addr,
		patten: opt.patten,

		authentication: opt.Authentication,

		Logger:     logx.WithContext(context.Background()),
		routes:     make(map[string]HandlerFunc),
		connToUser: make(map[*websocket.Conn]string),
		userToConn: make(map[string]*websocket.Conn),

		upgrader: websocket.Upgrader{},
	}
}

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			s.Errorf("server handler ws recover err %v", r)
		}
	}()

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Error("upgrade http conn err", err)
		return
	}

	if !s.authentication.Auth(w, r) {
		conn.WriteMessage(websocket.TextMessage, []byte("auth err"))
		conn.Close()
		return
	}
	// 添加连接记录,会有并发问题
	s.addConn(conn, r)
	// 读取信息，完成请求，还需建立连接
	go s.handlerConn(conn)
}

func (s *Server) handlerConn(conn *websocket.Conn) {
	// 记录连接
	for {
		// 读取请求信息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket conn readMessage err %v, user Id %s", err, "")
			// 关闭并删除连接
			s.Close(conn)
			return
		}

		// 解析请求信息
		var message Message
		json.Unmarshal(msg, &message)

		// 处理请求
		if handler, ok := s.routes[message.Method]; ok {
			handler(s, conn, &message)
		} else {
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%v not found", message.Method)))
		}
	}
}

func (s *Server) addConn(conn *websocket.Conn, req *http.Request) {
	// 此处是map的写操作，在操作上会存在并发的可能问题
	uid := s.authentication.UserId(req)

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

func (s *Server) GetConn(uid string) *websocket.Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	return s.userToConn[uid]
}

func (s *Server) GetConns(uids ...string) []*websocket.Conn {
	if len(uids) == 0 {
		return nil
	}

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	res := make([]*websocket.Conn, 0, len(uids))
	for _, uid := range uids {
		res = append(res, s.userToConn[uid])
	}
	return res
}

func (s *Server) GetUsers(conns ...*websocket.Conn) []string {

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	var res []string
	if len(conns) == 0 {
		// 获取全部
		res = make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			res = append(res, uid)
		}
	} else {
		// 获取部分
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.connToUser[conn])
		}
	}

	return res
}

func (s *Server) Close(conn *websocket.Conn) {
	conn.Close()

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.connToUser[conn]
	delete(s.connToUser, conn)
	delete(s.userToConn, uid)
}

func (s *Server) SendByUserId(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}

	return s.Send(msg, s.GetConns(sendIds...)...)
}

func (s *Server) Send(msg interface{}, conns ...*websocket.Conn) error {
	if len(conns) == 0 {
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		if err = conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) AddRoutes(rs []Route) {
	for _, r := range rs {
		s.routes[r.Method] = r.Handler
	}
}

func (s *Server) Start() {
	http.HandleFunc(s.patten, s.ServerWs)
	http.ListenAndServe(s.addr, nil)
}

func (s *Server) Stop() {
	fmt.Println("stop server")
}
