package gourd

import (
	"errors"
	"net/http"
	"strings"
)

// 如果路径为/admin/:name/info,该节点保存的信息包括斜杠开始的路径名
// 子节点、是否是通配符
// 优先级 静态>动态>通配
type routerNode struct {
	subPath          string                 // 当前节点的子路径
	childrenNodes    map[string]*routerNode //当前节点的子节点
	ifDynamic        bool                   // 是否为动态匹配节点（：）
	ifWildcard       bool                   // 是否为通配节点（*）只能用在最末尾
	ifEndPath        bool                   // 是否该节点是已注册的路由路径的终点
	handlerInterface HandlerInterface       // 节点相应的处理器接口
}

// 路由组，包含一个根节点
type routerManager struct {
	rootNode *routerNode
	engine   *gourdEngine
}

// 返回一个RouterGroup的组
func newRouterManager(engine *gourdEngine) *routerManager {
	return &routerManager{
		rootNode: &routerNode{
			subPath:          "",
			childrenNodes:    make(map[string]*routerNode),
			ifDynamic:        false,
			ifWildcard:       false,
			ifEndPath:        false,
			handlerInterface: nil,
		},
		engine: engine,
	}
}

// 添加router
func (r *routerManager) addRouter(path string, hi HandlerInterface) error {
	rn := r.rootNode
	subPaths := dividePath(path)
	return addNode(rn, subPaths, hi)
}

// 递归函数，向trie树添加节点
func addNode(rn *routerNode, subPaths []string, hi HandlerInterface) (err error) {
	err = nil
	sp := subPaths[0]
	ifDynamicNode := strings.Contains(sp, ":")
	ifWildcardNode := strings.Contains(sp, "*")
	if ifDynamicNode && sp[1:2] != ":" || ifWildcardNode && sp[1:2] != "*" {
		return errors.New("format error!")
	}
	if node, ok := rn.childrenNodes[sp]; ok {
		// 在map中找到了节点，此时有三种情况：
		// 1.节点ifEndPath为true且subPath是最后一个
		// 此时说明节点重复了
		// 2.subPath是最后一个但ifEndPath不为true，此时说明曾经注册过更长的路由
		// 如注册过/admin/delete，现在注册/admin，只需要把相关节点ifEndPath置true
		// 3.都不是，则还需继续向下匹配
		if len(subPaths) == 1 {
			if node.ifEndPath {
				err = errors.New("duplicate router!")
			} else {
				node.ifEndPath = true
			}
		} else {
			err = addNode(node, subPaths[1:], hi)
		}
	} else {
		// 注意，添加节点时不需要考虑查找时的问题，即：
		// 静态节点map匹配到，但递归下去却找不到节点，所以要回溯给动态节点
		// 在添加节点时，静态节点优先级高于动态节点，所以在这里不需要回溯
		//未找到节点，分以下几种情况：
		// 1.将注册的是动态节点，有动态节点，并且是EndPath且subPath是最后一个
		// 此时节点重复
		// 2.将注册的是静态节点或非1中情况的动态节点，直接添加节点，并继续递归
		// 3.将注册的是通配节点
		if ifDynamicNode {
			// 此时将注册的是动态节点
			hasDynamicNode := false
			for key, node := range rn.childrenNodes {
				// 查找到map已经注册了动态节点的情况
				if strings.Contains(key, ":") {
					hasDynamicNode = true
					if len(subPaths) == 1 {
						if node.ifEndPath {
							// 情况1，节点重复
							err = errors.New("duplicate router!")
						} else {
							// 未重复，但此时subPaths长度为1，到结尾，直接置ifEndPath为true
							node.ifEndPath = true
						}
					} else {
						// subPath未到结尾，且查找到动态节点，继续递归添加节点
						err = addNode(node, subPaths[1:], hi)
					}
				}
			}
			if !hasDynamicNode {
				// 如果未找到动态节点，注册动态节点
				rn.childrenNodes[sp] = &routerNode{
					subPath:          sp,
					childrenNodes:    make(map[string]*routerNode),
					ifDynamic:        true,
					ifEndPath:        len(subPaths) == 1,
					handlerInterface: hi,
				}
				// 如果subPaths长度不为一，则需要继续递归添加节点
				if len(subPaths) != 1 {
					err = addNode(rn.childrenNodes[sp], subPaths[1:], hi)
				}
			}
		} else if ifWildcardNode {
			//将注册的是通配节点，此时有以下几种情况
			// 1.对应的map有通配节点，此时重复
			// 2.此通配节点对应的subPath不为1，报错
			// 3.对应的map无通配节点，则注册相应节点
			if len(subPaths) != 1 {
				err = errors.New("wildcard should be the last node")
			} else {
				// 此时将注册的是通配节点
				hasWildcardNode := false
				for key, _ := range rn.childrenNodes {
					// 查找到map已经注册了通配节点的情况
					if strings.Contains(key, "*") {
						err = errors.New("duplicate router!")
						hasWildcardNode = true
					}
				}
				if !hasWildcardNode {
					// 如果未找到通配节点，注册通配节点
					rn.childrenNodes[sp] = &routerNode{
						subPath:          sp,
						childrenNodes:    make(map[string]*routerNode),
						ifWildcard:       true,
						ifEndPath:        true,
						handlerInterface: hi,
					}
				}
			}
		} else {
			// 未找到动态、通配节点，注册静态节点
			rn.childrenNodes[sp] = &routerNode{
				subPath:          sp,
				childrenNodes:    make(map[string]*routerNode),
				ifDynamic:        false,
				ifEndPath:        len(subPaths) == 1,
				handlerInterface: hi,
			}
			// 如果subPaths长度不为一，则需要继续递归添加节点
			if len(subPaths) != 1 {
				err = addNode(rn.childrenNodes[sp], subPaths[1:], hi)
			}
		}
	}
	return err
}

// 在router中查找handlerInterface
func (r *routerManager) findRouter(path string) (hi HandlerInterface, ifFind bool, params ParamMap) {
	params = make(map[string]string)
	rn := r.rootNode
	subPaths := dividePath(path)
	hi, ifFind = findNode(rn, subPaths, &params)
	return
}

// 递归函数，在trie树中查找节点
func findNode(rn *routerNode, subPaths []string, params *ParamMap) (hi HandlerInterface, ifFind bool) {
	sp := subPaths[0]
	hi = nil
	ifFind = false
	// 查找分以下几种种情况
	// 1.查找到静态节点
	// 1.1.该节点为终点且subPaths长度为1到达结尾，此时返回节点。
	// 1.2.若该节点为终点，但subPath长度不为1，说明未找到节点
	// 1.3.若该节点不为终点，subPath长度为1，说明未找到节点
	// 1.4.若该节点不为终点，subPath长度不为1，继续向下遍历
	if node, ok := rn.childrenNodes[sp]; ok {
		if len(subPaths) == 1 {
			if node.ifEndPath {
				// 情况1.1
				hi, ifFind = node.handlerInterface, true
				return
			}
		} else {
			if !node.ifEndPath {
				// 情况1.4
				hi, ifFind = findNode(node, subPaths[1:], params)
			}
		}
	}
	// 注意情况：静态节点中查找到路由，但查找下去并无结果，
	// 此时就要回溯到动态节点中进行匹配
	// 2.查找是否有动态节点，对每个动态节点，有：
	// 2.1.该动态节点为末节点，subPaths长度为1，此时返回节点
	// 2.2.该动态节点为末节点，subPaths长度不为1，此时未找到节点
	// 2.3.该动态节点不是末节点，但subPath长度为1，此时未找到节点
	// 2.4.该动态节点不是末节点，且subPath也不为1，继续递归查找
	// 或者回溯到通配节点中进行匹配
	// 此时只有一种情况，查找到通配节点，直接返回节点
	if !ifFind {
		for key, node := range rn.childrenNodes {
			// 查找到map已经注册了动态节点的情况
			if strings.Contains(key, ":") {
				if node.ifEndPath {
					if len(subPaths) == 1 {
						// 情况2.1
						hi, ifFind = node.handlerInterface, true
						(*params)[key[2:]] = sp[1:]
					}
				} else {
					if len(subPaths) != 1 {
						// 情况2.4
						hi, ifFind = findNode(node, subPaths[1:], params)
						(*params)[key[2:]] = sp[1:]
					}
				}
			} else if strings.Contains(key, "*") {
				hi, ifFind = node.handlerInterface, true
				pathParam := ""
				for _, subPath := range subPaths {
					pathParam += subPath
				}
				(*params)[key[2:]] = pathParam[1:]
			}
		}
	}
	return
}

// 分割path为子path
func dividePath(path string) []string {
	paths := strings.Split(path, "/")
	paths = paths[1:]
	var subPath []string
	for _, p := range paths {
		subPath = append(subPath, "/"+p)
	}
	return subPath
}

func (rm *routerManager) handle(context *Context) {
	// 通过路径查找到路由
	handlerInterface, ifFind, params := rm.findRouter(context.req.URL.Path)
	if ifFind {
		context.setParam(params)
		handlerInterface.setContext(context)
		if restHandler := handlerInterface.getRestHandler(context.Method); restHandler != nil {
			context.handlerfuncs = append(context.handlerfuncs, restHandler)
		} else {
			context.hi = handlerInterface
			handlerInterface.Prepare()
			context.handlerfuncs = append(context.handlerfuncs, nil)
		}
	} else {
		context.handlerfuncs = append(context.handlerfuncs, func(c *Context) {
			context.WriteString(http.StatusNotFound, "404 Not Found")
		})
	}
	context.Next()
}
