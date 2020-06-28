package gourd

import (
	"fmt"
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
	handlerInterface, ok := engine.rg.findRouter(r.URL.Path)
	if ok {
		context := NewContext(w, r)
		handlerInterface.setContext(context)
		switch context.Method {
		case "GET":
			handlerInterface.Get()
		case "POST":
			handlerInterface.Post()
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404notfound")
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
		fmt.Println(err)
	}
}
