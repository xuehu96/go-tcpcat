package server

import "go.uber.org/zap"

type Hook struct {
	OnListen   // Listen成功后回调
	OnAccept   // 新客户端连接调用
	OnReadData // 接收到客户端数据调用
	OnFnCode   // 从接收到的数据解析出功能代码
	OnSendData // 向客户端发送数据成功后调用
	OnClose    // 客户端主动断开连接或被断开连接后调用
	OnStop     // Listen关闭后回调
}

// OnListen Listen成功后回调
type OnListen func(s *Server)

// OnAccept 有新客户端连接，在这里可以做限流，或者关掉conn
// 返回true 允许连接， 返回false直接close 可以做来源或连接数判断
type OnAccept func(c *Client) bool

// OnReadData 接收到客户端数据调用
// 返回true 继续处理数据， 返回false 丢弃数据 可以做数据包格式校验，格式不对的丢掉
type OnReadData func(c *Client, buf []byte) bool

// OnFnCode 从接收到的数据解析出功能代码
type OnFnCode func(buf []byte) string

// OnSendData 向客户端发送数据成功后调用,可用于发送成功通知调用方
type OnSendData func(c *Client, buf []byte, err error)

// OnClose  客户端主动断开连接或被断开连接后调用
type OnClose func(clientKey string)

// OnStop Listen关闭后回调, 可用于清理打开的资源，不允许阻塞
type OnStop func()

// ---------- 默认回调函数 ------------

// defaultOnListen 默认的OnListen函数
func defaultOnListen(s *Server) {
	log.Info("xgtcp Server start listening...",
		zap.String("tcp", s.ln.Addr().String()),
		zap.Int("buf_len", s.bufLen))
	return
}

// defaultOnAccept 默认的OnAccept函数
func defaultOnAccept(c *Client) bool {
	log.Info("accept new connect", zap.String("addr", c.conn.RemoteAddr().String()))
	return true
}

// defaultOnReadData 默认的OnReadData函数
// 返回true 继续处理数据， 返回false 丢弃数据 可以做数据包格式校验，格式不对的丢掉
func defaultOnReadData(c *Client, buf []byte) bool {
	log.Debug("read success", zap.ByteString("data", buf), zap.Int("len", len(buf)))
	return true
}

// DefaultGetFnCode 默认的功能码解析函数
func defaultGetFnCode(buf []byte) string {
	return ""
}

// defaultOnSendData 默认的defaultOnSendData函数
// 参数1 发送的客户端，参数2 数据  参数3 发送的结果，是否出错，可以在这里写重发
func defaultOnSendData(c *Client, buf []byte, err error) {
	log.Info("data send", zap.String("to", c.key), zap.Error(err))
	return
}

// defaultOnClose 默认的defaultOnClose函数
func defaultOnClose(clientKey string) {
	log.Info("a Client close", zap.String("key", clientKey))
	return
}

// defaultOnStop Listen停止回调函数, 可用于清理打开的资源，不允许阻塞
func defaultOnStop() {
	log.Info("service stopped, exiting program...")
	return
}
