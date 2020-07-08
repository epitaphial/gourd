package gourd

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"
)

type ParamMap map[string]string

type Context struct {
	writer       http.ResponseWriter
	req          *http.Request
	Method       string
	Param        ParamMap
	Path         string
	data         map[string]interface{}
	mutex        sync.RWMutex
	index        int
	handlerfuncs []HandlerFunc // 中间件函数
	hi           HandlerInterface
	engine       *gourdEngine
	sessID       string
}

// 初始化上下文的操作，包括请求响应、方法
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		writer: w,
		req:    r,
		Method: r.Method,
		Param:  make(ParamMap),
		Path:   r.URL.Path,
		data:   make(map[string]interface{}),
		index:  -1,
		sessID: "",
	}
}

// 设置context的参数，用于动态和通配路由匹配时，获取相应参数的值
func (context *Context) setParam(pm ParamMap) {
	for k, v := range pm {
		context.Param[k] = v
	}
}

// 设置响应头的相关信息
func (context *Context) SetHeader(key string, value string) {
	context.mutex.Lock()
	defer context.mutex.Unlock()
	headMap := context.writer.Header()
	if _, ok := headMap[key]; !ok {
		context.writer.Header().Set(key, value)
	} else {
		context.writer.Header().Add(key, value)
	}
}

// 设置状态码，并完成响应头修改
func (context *Context) WriteHeader(code int) {
	context.mutex.Lock()
	defer context.mutex.Unlock()
	if code != http.StatusOK {
		context.writer.WriteHeader(code)
	}
}

// 向response报文写入内容
func (context *Context) write(str []byte) {
	context.writer.Write(str)
}

// 用于产生plaintext
func (context *Context) WriteString(code int, formart string, param ...interface{}) {
	context.SetHeader("Content-Type", "text/plain")
	context.WriteHeader(code)
	context.mutex.Lock()
	defer context.mutex.Unlock()
	fmt.Fprintf(context.writer, formart, param...)
}

// 用于post参数的查询
func (context *Context) Query(key string) string {
	return context.req.Form.Get(key)
}

// 用于路由重定向
func (context *Context) Redirect(code int, path string) {
	context.mutex.Lock()
	defer context.mutex.Unlock()
	http.Redirect(context.writer, context.req, path, code)
}

// 用于返回json信息
func (context *Context) WriteJson(code int, data interface{}) {
	if jsonData, err := json.Marshal(data); err == nil {
		context.SetHeader("Content-Type", "application/json")
		context.WriteHeader(code)
		context.write(jsonData)
	} else {
		// 错误处理
	}
}

// 用于给模板返回数据
func (context *Context) AddData(key string, dataIt interface{}) {
	context.mutex.Lock()
	defer context.mutex.Unlock()
	context.data[key] = dataIt
}

// 渲染模板
func (context *Context) RenderHTML(code int, htmlPath string) {
	context.WriteHeader(code)
	t := template.Must(template.ParseFiles(htmlPath))
	t.Execute(context.writer, context.data)
}

// 在中间件中使用
func (context *Context) Next() {
	context.index++
	for ; context.index < len(context.handlerfuncs); context.index++ {
		if context.handlerfuncs[context.index] == nil {
			handlerInterface := context.hi
			// 监听方法
			switch context.Method {
			case "GET":
				handlerInterface.Get()
			case "POST":
				handlerInterface.Post()
			case "HEAD":
				handlerInterface.Head()
			case "PUT":
				handlerInterface.Put()
			case "DELETE":
				handlerInterface.Delete()
			case "CONNECT":
				handlerInterface.Connect()
			case "OPTIONS":
				handlerInterface.Options()
			case "TRACE":
				handlerInterface.Trace()
			case "PATCH":
				handlerInterface.Patch()
			}
		} else {
			context.handlerfuncs[context.index](context)
		}
	}
}

func (context *Context) SetSession(sessionName string, sessionValue interface{}) {
	// 判断cookie是否已经存在
	sessCookie, err := context.req.Cookie(context.engine.smgr.cookieName)
	if err != nil {
		if context.sessID != "" {
			context.engine.smgr.setSession(context.sessID, sessionName, sessionValue)
		} else {
			sessId := getRandomString()
			cookie := http.Cookie{
				Name:     context.engine.smgr.cookieName,
				Value:    sessId,
				Path:     "/",
				MaxAge:   context.engine.smgr.maxLifeTime,
				HttpOnly: true,
			}
			context.SetHeader("Set-Cookie", cookie.String())
			// 客户端cookie注册
			context.engine.smgr.setSession(cookie.Value, sessionName, sessionValue)
			context.sessID = sessId
		}
	} else {
		// 客户端cookie注册
		context.engine.smgr.setSession(sessCookie.Value, sessionName, sessionValue)
	}
}

func (context *Context) GetSession(sessionName string) (sessionValue interface{}, err error) {
	if sessCookie, err1 := context.req.Cookie(context.engine.smgr.cookieName); err1 == nil {
		sessionValue, err = context.engine.smgr.getSessionValueBy(sessCookie.Value, sessionName)
	} else {
		err = errors.New("Can not find session in context.")
	}
	return
}

func (context *Context) DestroySession() (err error) {
	if sessCookie, err1 := context.req.Cookie(context.engine.smgr.cookieName); err1 != nil {
		err = errors.New("Can not find session in context.")
	} else {
		if err2 := context.engine.smgr.removeSessionValueBy(sessCookie.Value); err2 != nil {
			err = err2
		} else {
			cookie := http.Cookie{
				Name:    context.engine.smgr.cookieName,
				MaxAge:  -1,
				Expires: time.Now().Add(-100 * time.Hour), // Set expires for older versions of IE
				Path:    "/",
			}
			context.SetHeader("Set-Cookie", cookie.String())
			err = nil
		}
	}
	return
}
