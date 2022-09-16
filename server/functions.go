package server

// CodeFnType 功能函数类型
type CodeFnType func(c *Client, code string, data []byte)

// DefaultFn 默认功能码处理函数：什么也不做
func defaultFn(c *Client, code string, data []byte) {
}
