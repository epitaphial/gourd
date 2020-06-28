package gourd

import (
	"net/http"
)

// gourdEngine是整个框架的主引擎，实现了ServeHTTP方法，替换http包原有的DefaultMux
// 成员上，包含路由列表
// 方法上，包含构造方法Gourd、注册路由方法Route，启动引擎方法Run
type gourdEngine struct {
	rg *routerGroup
}

// Gourd方法是框架引擎的构造方法，返回引擎的指针
func Gourd() *gourdEngine {
	return &gourdEngine{
		rg: newRouterGroup(),
	}
}

// the function ListenAndServe receive a interface handler,
// all types have the ServeHTTP function can implement the interface
func (engine *gourdEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 通过路径查找到路由
	handlerInterface, ifFind ,params := engine.rg.findRouter(r.URL.Path)
	// 设置上下文、动态路由参数
	context := NewContext(w, r)
	context.setParam(params)
	if ifFind {
		// 监听方法
		handlerInterface.setContext(context)
		switch context.Method {
		case "GET":
			handlerInterface.Get()
		case "POST":
			handlerInterface.Post()
		}
	} else {
		context.SetStatus(http.StatusNotFound)
		context.WriteString("404notfound")
	}
}

// Run方法通过调用http包的ListenAndServe方法，启动服务器
func (engine *gourdEngine) Run() {
	http.ListenAndServe(":8080", engine)
}

// Route方法向路径注册相应的方法
func (engine *gourdEngine) Route(path string, handlerInterface HandlerInterface) {
	rg := engine.rg
	err := rg.addRouter(path, handlerInterface)
	if err != nil{
		// fmt.Println(err)
	}
}
