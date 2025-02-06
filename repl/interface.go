package repl

import "fmt"

const DEFAULT_INDENTATION uint = 0

type Context interface {
	ExprNum() uint
	Scope() string
	BumpExprNum() Context
}

type ReplContext struct {
	exprNum uint
	scope   string
}

func NewReplContext() *ReplContext {
	ctx := ReplContext{
		exprNum: 1,
		scope:   "main",
	}
	return &ctx
}

func (replCtx *ReplContext) ExprNum() uint {
	return replCtx.exprNum
}

func (replCtx *ReplContext) Scope() string {
	return replCtx.scope
}

func (replCtx *ReplContext) BumpExprNum() *ReplContext {
	ctx := ReplContext{
		exprNum: replCtx.exprNum + 1,
		scope:   replCtx.scope,
	}
	return &ctx
}

func Prompt(ctx *ReplContext) string {
	return fmt.Sprintf("lx(%s):%03d:%d> ", ctx.Scope(), ctx.ExprNum(), DEFAULT_INDENTATION)
}
