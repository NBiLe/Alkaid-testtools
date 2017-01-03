package menu

import (
	_L "github.com/fuyaocn/evaluatetools/log"
)

const (
	BACK_TO_MENU_FLAG = 99999
)

// ReturnParentMenu 返回上一级
type ReturnParentMenu struct {
	MenuSubItem
}

// InitMenu 初始化
func (ths *ReturnParentMenu) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.MenuSubItem.InitMenu(parent, key)
	ths.Exec = ths.execute

	ths.title = []string{
		LEnglish: "Go back",
	}
	return ths
}

func (ths *ReturnParentMenu) execute(isSync bool) {
	_L.LoggerInstance.Info(" ** go back ** \r\n")
	if !isSync {
		ths.ASyncChan <- BACK_TO_MENU_FLAG
	}
}
