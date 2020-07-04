package gourd

import (
	"net/http"
	"encoding/json"
	"html/template"
	"sync"
	"fmt"
)

type ParamMap map[string]string

type Context struct {
	writer http.ResponseWriter
	req    *http.Request
	Method string
	Param ParamMap
	data map[string]interface{}
	mutex sync.RWMutex
}

// 初始化上下文的操作，包括请求响应、方法
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		writer: w,
		req:    r,
		Method: r.Method,
		Param: make(ParamMap),
		data: make(map[string]interface{}),
	}
}

func (context *Context)setParam(pm ParamMap) {
	for k,v := range pm{
		context.Param[k] = v
	}
}

func (context *Context) SetHeader(key string, value string) {
	context.mutex.Lock()
	defer context.mutex.Unlock()
	headMap := context.writer.Header()
	if _,ok := headMap[key];!ok{
		context.writer.Header().Set(key, value)
	} else{
	}
}


func (context *Context) WriteHeader(code int) {
	context.mutex.Lock()
	defer context.mutex.Unlock()
	if code != http.StatusOK{
		context.writer.WriteHeader(code)
	}
}

func (context *Context) write(str []byte){
	context.writer.Write(str)
}

// 用于产生plaintext
func (context *Context) WriteString(code int,formart string,param ...interface{}) {
	context.SetHeader("Content-Type","text/plain")
	context.WriteHeader(code)
	fmt.Fprintf(context.writer,formart,param...)
}

func (context *Context) Query(key string) string {
	return context.req.Form.Get(key)
}

func (context *Context) Redirect(code int,path string){
	context.mutex.Lock()
	defer context.mutex.Unlock()
	headMap := context.writer.Header()
	if len(headMap) == 0 {
		http.Redirect(context.writer,context.req, path, code)
	}
}

func (context *Context) WriteJson(code int,data interface{}){
	if jsonData,err := json.Marshal(data);err == nil{
		context.SetHeader("Content-Type","application/json")
		context.WriteHeader(code)
		context.write(jsonData)
	} else {
		// 错误处理
	}
}

func (context *Context) AddData(key string,dataIt interface{}) {
	context.mutex.Lock()
	defer context.mutex.Unlock()
	context.data[key] = dataIt
}

func (context *Context) RenderHTML(htmlPath string) {
	t := template.Must(template.ParseFiles(htmlPath))
	t.Execute(context.writer,context.data)
}