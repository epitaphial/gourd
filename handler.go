package gourd

// HandlerInterface是包含了诸多方法的接口
type HandlerInterface interface {
	Prepare()
	Get()
	Post()
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
}

func (handler *Handler) Prepare() {
}

func (handler *Handler) Get() {
}

func (handler *Handler) Post() {
}
