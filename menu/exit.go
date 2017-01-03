package menu

import (
	"os"

	_L "github.com/fuyaocn/evaluatetools/log"
)

// ExitApp 退出钱包程序
type ExitApp struct {
	MenuSubItem
}

// InitMenu 初始化
func (ths *ExitApp) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.MenuSubItem.InitMenu(parent, key)
	ths.Exec = ths.execute

	ths.title = []string{
		LEnglish: "Exit",
	}
	return ths
}

func (ths *ExitApp) execute(isSync bool) {
	_L.LoggerInstance.Info(" ** Exit app ** \r\n")
	os.Exit(0)
}
