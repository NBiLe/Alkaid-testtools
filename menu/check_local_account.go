package menu

import (
	_dis "github.com/fuyaocn/evaluatetools/distribute"
	_L "github.com/fuyaocn/evaluatetools/log"
)

// CheckAccounts 检查所有数据
type CheckAccounts struct {
	MenuSubItem
}

// InitMenu 初始化
func (ths *CheckAccounts) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.MenuSubItem.InitMenu(parent, key)
	ths.Exec = ths.execute

	ths.title = []string{
		LEnglish: "Check all local stellar accounts",
	}
	return ths
}

func (ths *CheckAccounts) execute(isSync bool) {
	_L.LoggerInstance.Info(" ** Check local stellar accounts ** \r\n")
	_L.LoggerInstance.InfoPrint(" > Check all account from local database BEGIN ...\r\n")
	ca := &_dis.CAMainController{}
	ca.Init(0, 0, 0, nil)
	ca.Execute()
	_L.LoggerInstance.InfoPrint(" > Check all account from local database Complete\r\n")

	if !isSync {
		ths.ASyncChan <- 0
	}
}
