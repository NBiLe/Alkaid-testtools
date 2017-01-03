package distribute

import (
	_str "github.com/fuyaocn/evaluatetools/stellar"
)

// IController 基础控制器的接口定义
type IController interface {
	Init(l, bl, dl int, a *_str.AccountInfo) IController
	Execute()
}

// Executor 执行基础定义
type Executor struct {
	Level       int // 当前等级
	BaseLevel   int // 基础等级 第一级 之后将按照这个等级扩展
	DepthLevel  int // 深度 做多扩展几级
	BaseAccount *_str.AccountInfo
	SubAccounts []*_str.AccountInfo
	executor    []IController
}

// Init 初始化
func (ths *Executor) Init(l, bl, dl int, a *_str.AccountInfo) IController {
	ths.Level = l
	ths.BaseLevel = bl
	ths.DepthLevel = dl
	ths.BaseAccount = a
	ths.executor = make([]IController, ths.BaseLevel)
	return ths
}

// Execute 执行
func (ths *Executor) Execute() {
	panic("Must override 'Execute' function !!!")
}
