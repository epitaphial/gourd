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
}

// Handler实现了HandlerInterface接口的所有方法
// 包括成员Ctx，该成员含有上下文
type Handler struct {
	Ctx *Context
}


// 该方法设置相关的上下文
func (handler *Handler) setContext(ctx *Context) {
	handler.Ctx = ctx
	handler.Ctx.req.ParseForm()
}

func (handler *Handler) Prepare() {
}

func (handler *Handler) Get() {
	handler.Ctx.WriteString(405,"405 Method Not Allowed")
}

func (handler *Handler) Post() {
	handler.Ctx.WriteString(405,"405 Method Not Allowed")
}

func (handler *Handler) Head() {
	handler.Ctx.WriteString(405,"405 Method Not Allowed")
}

func (handler *Handler) Put() {
	handler.Ctx.WriteString(405,"405 Method Not Allowed")
}

func (handler *Handler) Delete() {
	handler.Ctx.WriteString(405,"405 Method Not Allowed")
}

func (handler *Handler) Connect() {
	handler.Ctx.WriteString(405,"405 Method Not Allowed")
}

func (handler *Handler) Options() {
	handler.Ctx.WriteString(405,"405 Method Not Allowed")
}

func (handler *Handler) Trace() {
	handler.Ctx.WriteString(405,"405 Method Not Allowed")
}

func (handler *Handler) Patch() {
	handler.Ctx.WriteString(405,"405 Method Not Allowed")
}