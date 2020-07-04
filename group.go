package gourd

import(
	"errors"
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

func (rg *routerGroup) Use(middlewares ...HandlerFunc) {
	rg.middlewares = append(rg.middlewares,middlewares...)
}