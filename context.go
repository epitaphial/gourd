package gourd

import (
	"net/http"
	"encoding/json"
)

type ParamMap map[string]string

type Context struct {
	writer http.ResponseWriter
	req    *http.Request
	Method string
	Param ParamMap
}

// 初始化上下文的操作，包括请求响应、方法
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		writer: w,
		req:    r,
		Method: r.Method,
		Param: make(ParamMap),
	}
}

func (context *Context)setParam(pm ParamMap) {
	for k,v := range pm{
		context.Param[k] = v
	}
}

func (context *Context) WriteString(text string) {
	context.writer.Write([]byte(text))
}

func (context *Context) SetStatus(code int) {
	context.writer.WriteHeader(code)
}

func (context *Context) Query(key string) string {
	return context.req.Form.Get(key)
}

func (context *Context) Redirect(code int,path string){
	http.Redirect(context.writer,context.req, path, code)
}

func (context *Context) WriteJson(data interface{}){
	context.writer.Header().Set("Content-Type","application/json")
	if jsonData,err := json.Marshal(data);err == nil{
		context.writer.Write(jsonData)
	} else {
		// 错误处理
	}
}