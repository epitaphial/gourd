package gourd

import(
	"errors"
	"net/http"
	"path"
)

type routerGroup struct{
	prefix string
	middlewares []HandlerFunc
	parent *routerGroup
	engine *gourdEngine
}

func (rg *routerGroup) Group(groupPath string) (*routerGroup,error){
	for _,group := range rg.engine.groups{
		if rg.prefix + groupPath == group.prefix {
			panic(errors.New("duplicate group!"))
		}
	}
	newg := &routerGroup{
		prefix: rg.prefix + groupPath,
		parent: rg,
		engine: rg.engine,
	}
	rg.engine.groups = append(rg.engine.groups,newg)
	return newg,nil
}

func (rg *routerGroup) Route(path string, handlerInterface HandlerInterface){
	rg.engine.Route(rg.prefix + path, handlerInterface)
}

// RESTFUL路由
func (rg *routerGroup) Get(path string, hf HandlerFunc) {
	rg.engine.routeRest("GET", rg.prefix + path, hf)
}

func (rg *routerGroup) Post(path string, hf HandlerFunc) {
	rg.engine.routeRest("POST", rg.prefix + path, hf)
}

func (rg *routerGroup) Head(path string, hf HandlerFunc) {
	rg.engine.routeRest("HEAD", rg.prefix + path, hf)
}

func (rg *routerGroup) Put(path string, hf HandlerFunc) {
	rg.engine.routeRest("PUT", rg.prefix + path, hf)
}

func (rg *routerGroup) Delete(path string, hf HandlerFunc) {
	rg.engine.routeRest("DELETE", rg.prefix + path, hf)
}

func (rg *routerGroup) Connect(path string, hf HandlerFunc) {
	rg.engine.routeRest("CONNECT", rg.prefix + path, hf)
}

func (rg *routerGroup) Options(path string, hf HandlerFunc) {
	rg.engine.routeRest("OPTIONS", rg.prefix + path, hf)
}

func (rg *routerGroup) Trace(path string, hf HandlerFunc) {
	rg.engine.routeRest("TRACE", rg.prefix + path, hf)
}

func (rg *routerGroup) Patch(path string, hf HandlerFunc) {
	rg.engine.routeRest("PATCH", rg.prefix + path, hf)
}

// 中间件应用
func (rg *routerGroup) Use(middlewares ...HandlerFunc) {
	rg.middlewares = append(rg.middlewares,middlewares...)
}

func (rg *routerGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(rg.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(context *Context) {
		file := context.Param["staticpath"]
		if _, err := fs.Open(file); err != nil {
			context.WriteHeader(http.StatusNotFound)
		} else {
			fileServer.ServeHTTP(context.writer, context.req)
		}
	}
}

func (rg *routerGroup) StaticDir(relativePath string, root string) {
	handler := rg.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*staticpath")
	rg.Get(urlPattern, handler)
}