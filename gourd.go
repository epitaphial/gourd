package gourd

import (
	"net/http"
	"strings"
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
	// 设置上下文、动态路由参数
	context := NewContext(w, r)
	for _,group := range engine.groups{
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			context.handlerfuncs = append(context.handlerfuncs,group.middlewares...)
		}
	}
	engine.rm.handle(context)
}

// Run方法通过调用http包的ListenAndServe方法，启动服务器
func (engine *gourdEngine) Run(port string) {
	err := http.ListenAndServe(port, engine)
	if err != nil{
		panic(err)
	}
}

// Route方法向路径注册相应的方法
func (engine *gourdEngine) Route(path string, handlerInterface HandlerInterface) {
	err := engine.rm.addRouter(path, handlerInterface)
	if err != nil{
		panic(err)
	}
}
