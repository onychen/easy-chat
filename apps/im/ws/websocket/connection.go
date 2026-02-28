package websocket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Conn struct {
	*websocket.Conn

	s *Server

	Uid string

	idleMu            sync.Mutex // guard the following
	idle              time.Time
	maxConnectionIdle time.Duration

	done chan struct{}
}

func NewConn(s *Server, w http.ResponseWriter, r *http.Request) *Conn {

	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Error("upgrade http conn err", err)
		return nil
	}

	conn := &Conn{
		Conn:              c,
		s:                 s,
		idle:              time.Now(),
		done:              make(chan struct{}),
		maxConnectionIdle: s.opt.maxConnectionIdle,
	}

	go conn.keepalive()
	return conn
}

func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	// 开始忙碌
	messageType, p, err = c.Conn.ReadMessage()
	c.idle = time.Time{}
	return
}

func (c *Conn) WriteMessage(messageType int, data []byte) error {
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	err := c.Conn.WriteMessage(messageType, data)
	// 当写操作完成后当前连接就会进入空闲状态，并记录空闲的时间
	c.idle = time.Now()
	return err
}

func (c *Conn) Close() error {
	select {
	case <-c.done:
	default:
		return c.Conn.Close()
	}

	return c.Conn.Close()
}

// 长连接检测机制(思路是grpc的keepalive)
func (c *Conn) keepalive() {
	idleTimer := time.NewTimer(c.maxConnectionIdle)
	defer idleTimer.Stop()

	for {
		select {
		case <-idleTimer.C:
			c.idleMu.Lock()
			defer c.idleMu.Unlock()
			idle := c.idle

			fmt.Printf("idle %v, maxIdle %v \n", c.idle, c.maxConnectionIdle)
			if idle.IsZero() { // The connection is non-idle.
				c.idleMu.Unlock()
				idleTimer.Reset(c.maxConnectionIdle)
				continue
			}
			val := c.maxConnectionIdle - time.Since(idle)
			fmt.Printf("val %v \n", val)
			c.idleMu.Unlock()
			if val <= 0 {
				// The connection has been idle for a duration of keepalive.MaxConnectionIdle or more.
				// Gracefully close the connection.
				c.s.Close(c)
				return
			}
			idleTimer.Reset(val)
		case <-c.done:
			fmt.Println("客户端结束连接")
			return
		}
	}
}
