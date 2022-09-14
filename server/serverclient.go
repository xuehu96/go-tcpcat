package server

import (
	"bufio"
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net"
	"sync"
)

// NewClient 返回client实例
func NewClient(s *Server, conn net.Conn, id int64) *Client {
	c := &Client{
		id:   id,
		s:    s,
		conn: conn,
		mu:   sync.RWMutex{},
	}
	c.key = conn.RemoteAddr().String()

	s.mu.Lock()
	s.clients[c.key] = c
	s.mu.Unlock()

	log.Info("new Client", zap.Int64("id", id), zap.String("key", c.key))
	return c
}

// Start TCP服务器的client开始干活
func (c *Client) Start() {
	// OnAccept hook
	ok := c.s.hooks.OnAccept(c)
	if !ok {
		c.closed = true
		c.conn.Close()
		return
	}

	c.ctx, c.cancel = context.WithCancel(context.Background())

	go c.ClientReader()

	select {
	case <-c.ctx.Done():
		//退出

		if c.closed {
			return
		}
		// hook onclose
		go c.s.hooks.OnClose(c.key)

		c.conn.Close()
		c.closed = true

		// 删除s中的c
		c.s.mu.Lock()
		delete(c.s.clients, c.key)
		c.s.mu.Unlock()
		return
	}
}

// Close 关闭客户端
func (c *Client) Close() {
	c.cancel()
}

// ClientReader 循环读取客户端发的数据
func (c *Client) ClientReader() {
	defer c.Close()
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			// 1. 循环接收数据
			reader := bufio.NewReader(c.conn)
			buf := make([]byte, c.s.bufLen)
			n, err := reader.Read(buf[:])
			if err != nil {
				log.Debug("Client read err (close)", zap.Error(err))
				return
			}
			ok := c.s.hooks.OnReadData(c, buf, n)
			if !ok {
				continue
			}
			// 2. 获取功能码
			fnc := c.s.hooks.OnFnCode(buf)
			if fnc == "" {
				go defaultFn(c, fnc, buf, n)
				continue
			}
			// 3. 根据功能码 回调处理函数
			fn, ok := c.s.fns[fnc]
			if ok {
				//调用对应功能码的处理函数
				go fn(c, fnc, buf, n)
			} else {
				//调用默认处理函数
				go defaultFn(c, fnc, buf, n)
			}

		} // end select

	}
}

// ReplyData 向客户端回写数据
func (c *Client) ReplyData(data []byte) (err error) {
	if c.closed {
		go c.s.hooks.OnSendData(nil, data, errors.New("Client is closed"))
		return errors.New("Client is closed")
	}
	_, err = c.conn.Write(data)
	go c.s.hooks.OnSendData(c, data, err)
	return
}

// SendDataTo 向指定客户端发送数据
func (s *Server) SendDataTo(ClientKey string, data []byte) error {
	c, ok := s.clients[ClientKey]
	if !ok {
		go s.hooks.OnSendData(nil, data, errors.New("Client key not found"))
		return errors.New("Client key not found")
	}
	return c.ReplyData(data)
}
