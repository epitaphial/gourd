package gourd

import (
	"net/http"
)

// gourdEngine是整个框架的主引擎，实现了ServeHTTP方法，替换http包原有的DefaultMux
// 成员上，包含路由列表
// 方法上，包含构造方法Gourd、注册路由方法Route，启动引擎方法Run
type gourdEngine struct {
	*routerGroup // engine作为顶层的group
	rm *routerManager // 所有的router
	groups []*routerGroup // engine管理的所有的group
}

// Gourd方法是框架引擎的构造方法，返回引擎的指针
func Gourd() *gourdEngine {
	engine := &gourdEngine{
		rm: newRouterManager(),
	}
	engine.routerGroup = &routerGroup{engine:engine}
	engine.groups = []*routerGroup{engine.routerGroup}
	return engine
}

// the function ListenAndServe receive a interface handler,
// all types have the ServeHTTP function can implement the interface
func (engine *gourdEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 通过路径查找到路由
	handlerInterface, ifFind ,params := engine.rm.findRouter(r.URL.Path)
	// 设置上下文、动态路由参数
	context := NewContext(w, r)
	context.setParam(params)
	if ifFind {
		// 监听方法
		handlerInterface.setContext(context)
		handlerInterface.Prepare()
		//fmt.Printf("METHOD-%s-PATH-%s\n",context.Method,r.URL.Path)
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
		context.WriteString(http.StatusNotFound,"404 Not Found")
	}
}

// Run方法通过调用http包的ListenAndServe方法，启动服务器
func (engine *gourdEngine) Run() {
	http.ListenAndServe(":8080", engine)
}

// Route方法向路径注册相应的方法
func (engine *gourdEngine) Route(path string, handlerInterface HandlerInterface) {
	err := engine.rm.addRouter(path, handlerInterface)
	if err != nil{
		panic(err)
	}
}
