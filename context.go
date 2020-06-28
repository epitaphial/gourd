package gourd

import (
	"net/http"
)

type Context struct {
	writer http.ResponseWriter
	req    *http.Request
	Method string
}

// 初始化上下文的操作，包括请求响应、方法
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		writer: w,
		req:    r,
		Method: r.Method,
	}
}

func (context *Context) WriteText(text string) {
	context.writer.Write([]byte(text))
}
