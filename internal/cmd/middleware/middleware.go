package middleware

import "context"

type Context struct {
	context.Context
}

type Middleware func(ctx *Context)

func NewContext() *Context {
	return &Context{
		Context: context.Background(),
	}
}

func (ctx *Context) Set(key, value interface{}) {
	ctx.Context = context.WithValue(ctx.Context, key, value)
}

func Compute(ctx *Context, middlewares ...Middleware) func() {
	return func() {
		for _, m := range middlewares {
			m(ctx)
		}
	}
}
