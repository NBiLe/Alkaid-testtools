package menu

import (
	_db "github.com/fuyaocn/evaluatetools/db"
	_L "github.com/fuyaocn/evaluatetools/log"
)

// ClearAccounts 清空所有数据
type ClearAccounts struct {
	MenuSubItem
}

// InitMenu 初始化
func (ths *ClearAccounts) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.MenuSubItem.InitMenu(parent, key)
	ths.Exec = ths.execute

	ths.title = []string{
		LEnglish: "Clear all local stellar accounts",
	}
	return ths
}

func (ths *ClearAccounts) execute(isSync bool) {
	_L.LoggerInstance.Info(" ** Clear local stellar accounts ** \r\n")
	_L.LoggerInstance.InfoPrint(" > Clear all account from local database BEGIN ...\r\n")
	err := _db.DataBaseInstance.ClearAllAccount()
	if err != nil {
		_L.LoggerInstance.ErrorPrint("Clear local stellar account from database has error :\r\n%+v\r\n", err)
	}
	_L.LoggerInstance.InfoPrint(" > Clear all account from local database Complete\r\n")

	if !isSync {
		ths.ASyncChan <- 0
	}
}
