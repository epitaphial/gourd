# Gourd Web Framework

![author](https://img.shields.io/badge/author-Curled-blueviolet.svg?style=plastic)
![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=plastic)

## DESCRIPTION

GOURD是一款轻量、简单、易用的golang web框架。

## INSTALLATION

假设用户已经安装好golang环境，安装方法如下：

```bash
$ go get -u github.com/epitaphial/gourd
```

在go源文件中使用：

```go
import "github.com/epitaphial/gourd"
```

## GLIMPSE

创建源文件main.go，内容如下：

```go
package main

import (
    "github.com/epitaphial/gourd"
)
type indexHandler struct {
	gourd.Handler
}

func (idx *indexHandler) Get() {
	idx.Ctx.WriteString(200, "hello,gourd!")
}
func main() {
	engine := gourd.Gourd()
	engine.Route("/", &indexHandler{})
    engine.Run(":8080")
}
```

打开控制台，运行：

```bash
$ go run main.go
```

访问http://127.0.0.1:8080/，即可看到文字：

"hello,gourd!"

## GOURD APIS

### 初始化

#### 1. 引擎初始化

```go
package main

import (
    "github.com/epitaphial/gourd"
)

func main() {
	engine := gourd.Gourd()
    engine.Run(":8080")
}
```

通过gourd包的Gourd方法，可以返回gourd引擎的一个实例，调用engine的方法Run，参数为监听的端口，这就完成了gourd的初始化与监听

### 路由相关

#### 1.一般路由

```go
package main

import (
    "github.com/epitaphial/gourd"
)
type indexHandler struct {
	gourd.Handler
}

func (idx *indexHandler) Get() {
	idx.Ctx.WriteString(200, "hello,gourd!")
}
func main() {
	engine := gourd.Gourd()
	engine.Route("/", &indexHandler{})
    engine.Run(":8080")
}
```

回到最开始的例子，路由的注册通过gourd的Route方法实现，该方法接受两个参数，string类型的路径，以及一个gourd.HandlerInterface类型的接口，该接口定义在handler.go文件内。我们可以创建一个结构，该结构继承了gourd.Handler结构，Handler结构实现了HandlerInterface的所有方法，因此我们创建的结构也能赋值给该类型接口。我们只需要重写相应的方法，支持的方法有：GET、POST、HEAD、PUT等。

#### 2.RESTFUL路由

```go
...
engine.Get("/get", func(context *gourd.Context) {
    context.WriteString(200, "get test")
})
...
```

gourd也支持RESTFUL路由，直接调用engine的Get、Post等方法，这些方法接受两个参数：路径和参数为gourd.Context指针的函数。

#### 3.动态路径路由

```go
....
engine.Route("/static/*relapath", &staticHandler{})
engine.Route("/:username/info", &userHandler{})
....
```

子路径：子路径是指两个'/'之间的路径，包含开头的'/'

gourd支持两种类型的动态路径路由。带\*的路由匹配0个或多个子路径，而带':'的路由只匹配一个路径，例如，"/static/\*relapath"可以匹配"/static/js/bootstrap.js","/static/","/static/css/bootstrap.css"。而"/:username/info"可以匹配"/curled/info"。

动态路径中的参数通过上下文\*gourd.Context的成员Param查询，例如：

```
context.Param["relapath"]
```

#### 4.路由分组

```go
group1, err := engine.Group("/user")
```

分组通过调用engine的Group方法产生一个新的分组，新的分组类型继承了engine的大多数方法。包括注册路由，使用中间件，以及产生新的分组，分组支持嵌套，示例如下：

```go
group2, err := group1.Group("/admin")
```

### 上下文相关

上下文，gourd.Context，是定义在context.go的一个结构体，通过对上下文的修改，以进一步改变返回的数据。

在非restful路由中，我们使用继承gourd.Handler的方式，Handler中包含了一个\*gourd.Context类型的成员Ctx，该成员即上下文。而在restful型路由中，该类路由注册时接受的参数就是带有\*gourd.Context类型的函数，直接操作即可。

#### 1.设置修改响应头

```go
// 向header中写入键值对
func (context *Context) SetHeader(key string, value string)
```

```go
// 设置响应码，完成header的修改，务必在SetHeader之后调用
func (context *Context) WriteHeader(code int)
```

#### 2.返回数据相关

```go
// 返回plaintext
func (context *Context) WriteString(code int, formart string, param ...interface{})
// 返回json
func (context *Context) WriteJson(code int, data interface{})
```

#### 3.路由重定向

```go
func (context *Context) Redirect(code int, path string)
```

#### 4.POST数据获取

```go
func (context *Context) Query(key string) string
```

#### 5.模板相关

```go
// 给渲染的模板添加数据
func (context *Context) AddData(key string, dataIt interface{}) 
// 渲染模板
func (context *Context) RenderHTML(code int, htmlPath string)
```

### 中间件

中间件可以在上下文执行前后进行一系列的操作，gourd提供了中间件的接口，中间件以带有Context参数的函数为签名，以路由组为粒度（也就是说，分组也继承了中间件的Use方法），以下是一个简单的中间件示例。

```go
engine.Use(func(context *gourd.Context) {
		// 开始计时
		t := time.Now()
		// 先处理后面的中间件以及上下文
		context.Next()
		// 最后再计算时间
		log.Printf("path %s in %v", context.Path, time.Since(t))
	})
```

通过Next方法，我们可以先处理后面的中间件，待后面的中间件执行完成后，再回到该中间件进行后续处理。中间件的顺序类似于队列，先进先出，先注册先调用。

### 静态文件处理

有时我们需要处理静态文件，如js、css等，gourd提供了相关的功能。

```go
engine.StaticDir("/sta", "./static")
```

通过调用StaticDir方法，把/sta路由映射到文件路径./static，以达到访问静态文件目的。（绝对路径和相对路径都被支持）

### Session控制

gourd提供了三个方法用于Session控制

```go
// 设置session
func (context *Context) SetSession(sessionName string, sessionValue interface{})
// 获取sessionName对应的session
func (context *Context) GetSession(sessionName string) (sessionValue interface{}, err error)
// 清除session
func (context *Context) DestroySession() (err error)
```

