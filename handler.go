package gourd

type HandlerFunc func(*Context)

// HandlerInterface是包含了诸多方法的接口
type HandlerInterface interface {
	Prepare()
	Get()
	Post()
	Head()
	Put()
	Delete()
	Connect()
	Options()
	Trace()
	Patch()
	setContext(ctx *Context)
	getRestHandler(string) HandlerFunc
	getRestHandlers() map[string]HandlerFunc
	setRestful(bool)
	ifRestHandler() bool
	setRestHandler(string, HandlerFunc)
}

// Handler实现了HandlerInterface接口的所有方法
// 包括成员Ctx，该成员含有上下文
type Handler struct {
	Ctx          *Context
	handlerFuncs map[string]HandlerFunc
	restful      bool
}

func (handler *Handler) getRestHandlers() (hfs map[string]HandlerFunc) {
	if handler.restful {
		hfs = handler.handlerFuncs
	}
	return
}

// 如果handler是restful路由，则获取其handler函数
func (handler *Handler) getRestHandler(method string) (hf HandlerFunc) {
	if handler.restful {
		hf = handler.handlerFuncs[method]
	}
	return
}

// 设置restful
func (handler *Handler) setRestful(set bool) {
	handler.restful = set
}

// 设置restful
func (handler *Handler) setRestHandler(method string, hf HandlerFunc) {
	handler.handlerFuncs[method] = hf
}

// 判断是否是restful
func (handler *Handler) ifRestHandler() bool {
	return handler.restful
}

// 该方法设置相关的上下文
func (handler *Handler) setContext(ctx *Context) {
	handler.Ctx = ctx
	handler.Ctx.req.ParseForm()
}

func (handler *Handler) Prepare() {
}

func (handler *Handler) Get() {
	handler.Ctx.WriteString(405, "405 Method Not Allowed")
}

func (handler *Handler) Post() {
	handler.Ctx.WriteString(405, "405 Method Not Allowed")
}

func (handler *Handler) Head() {
	handler.Ctx.WriteString(405, "405 Method Not Allowed")
}

func (handler *Handler) Put() {
	handler.Ctx.WriteString(405, "405 Method Not Allowed")
}

func (handler *Handler) Delete() {
	handler.Ctx.WriteString(405, "405 Method Not Allowed")
}

func (handler *Handler) Connect() {
	handler.Ctx.WriteString(405, "405 Method Not Allowed")
}

func (handler *Handler) Options() {
	handler.Ctx.WriteString(405, "405 Method Not Allowed")
}

func (handler *Handler) Trace() {
	handler.Ctx.WriteString(405, "405 Method Not Allowed")
}

func (handler *Handler) Patch() {
	handler.Ctx.WriteString(405, "405 Method Not Allowed")
}
