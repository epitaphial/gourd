package gourd

import (
	"log"
	"net/http"
	"strings"
)

// gourdEngine是整个框架的主引擎，实现了ServeHTTP方法，替换http包原有的DefaultMux
// 成员上，包含路由列表
// 方法上，包含构造方法Gourd、注册路由方法Route，启动引擎方法Run
type gourdEngine struct {
	*routerGroup                // engine作为顶层的group
	rm           *routerManager // 所有的router
	groups       []*routerGroup // engine管理的所有的group
}

// Gourd方法是框架引擎的构造方法，返回引擎的指针
func Gourd() *gourdEngine {
	engine := &gourdEngine{
		rm: newRouterManager(),
	}
	engine.routerGroup = &routerGroup{engine: engine}
	engine.groups = []*routerGroup{engine.routerGroup}
	// 默认使用recovery中间件
	engine.Use(Recovery())
	return engine
}

// the function ListenAndServe receive a interface handler,
// all types have the ServeHTTP function can implement the interface
func (engine *gourdEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 设置上下文、动态路由参数
	context := NewContext(w, r)
	for _, group := range engine.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			context.handlerfuncs = append(context.handlerfuncs, group.middlewares...)
		}
	}
	engine.rm.handle(context)
}

// Run方法通过调用http包的ListenAndServe方法，启动服务器
func (engine *gourdEngine) Run(port string) {
	log.Printf("Running in port %s\n", port[1:])
	err := http.ListenAndServe(port, engine)
	if err != nil {
		panic(err)
	}
}

// Route方法向路径注册相应的方法,此时非rest模式
func (engine *gourdEngine) Route(path string, handlerInterface HandlerInterface) {
	handlerInterface.setRestful(false)
	err := engine.rm.addRouter(path, handlerInterface)
	if err != nil {
		panic(err)
	}
}

// RESTFUL路由
func (engine *gourdEngine) Get(path string, hf HandlerFunc) {
	engine.routeRest("GET", path, hf)
}

func (engine *gourdEngine) Post(path string, hf HandlerFunc) {
	engine.routeRest("POST", path, hf)
}

func (engine *gourdEngine) Head(path string, hf HandlerFunc) {
	engine.routeRest("HEAD", path, hf)
}

func (engine *gourdEngine) Put(path string, hf HandlerFunc) {
	engine.routeRest("PUT", path, hf)
}

func (engine *gourdEngine) Delete(path string, hf HandlerFunc) {
	engine.routeRest("DELETE", path, hf)
}

func (engine *gourdEngine) Connect(path string, hf HandlerFunc) {
	engine.routeRest("CONNECT", path, hf)
}

func (engine *gourdEngine) Options(path string, hf HandlerFunc) {
	engine.routeRest("OPTIONS", path, hf)
}

func (engine *gourdEngine) Trace(path string, hf HandlerFunc) {
	engine.routeRest("TRACE", path, hf)
}

func (engine *gourdEngine) Patch(path string, hf HandlerFunc) {
	engine.routeRest("PATCH", path, hf)
}

func (engine *gourdEngine) routeRest(method string, path string, hf HandlerFunc) {
	hi, ifFind, _ := engine.rm.findRouter(path)
	if ifFind {
		if hi.ifRestHandler() { // 是rest类handler，只需要判断方法层面上有没有重复注册
			if restHandlers := hi.getRestHandlers(); restHandlers != nil {
				if _, dep := restHandlers[method]; dep {
					// 如果restful路由方法重复
					panic("can not register restful route when there is already an unrestful route!")
				} else {
					hi.setRestHandler(method, hf)
				}
			}
		} else { // 不是rest类handler，不满足要求，从而也不允许注册
			panic("can not register restful route when there is already an unrestful route!")
		}
	} else {
		hd := Handler{
			handlerFuncs: make(map[string]HandlerFunc),
			restful:      true,
		}
		hd.handlerFuncs[method] = hf
		var hi HandlerInterface
		hi = &hd
		err := engine.rm.addRouter(path, hi)
		if err != nil {
			panic(err)
		}
	}
}
