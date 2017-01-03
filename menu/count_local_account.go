package menu

import (
	"fmt"

	_db "github.com/fuyaocn/evaluatetools/db"
	_L "github.com/fuyaocn/evaluatetools/log"
)

// CountAccount 获取数量
type CountAccount struct {
	MenuSubItem
}

// InitMenu 初始化
func (ths *CountAccount) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.MenuSubItem.InitMenu(parent, key)
	ths.Exec = ths.execute

	ths.title = []string{
		LEnglish: "Get all local stellar account count",
	}
	return ths
}

func (ths *CountAccount) execute(isSync bool) {
	_L.LoggerInstance.Info(" ** Get local stellar account count ** \r\n")
	var factive int64 = 0
	err := _db.DataBaseInstance.GetAccountCount("F", false, &factive)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("Get local stellar account count from database has error :\r\n%+v\r\n", err)
		return
	}
	var allactive int64 = 0
	err = _db.DataBaseInstance.GetAccountCount("", true, &allactive)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("Get local stellar account count from database has error :\r\n%+v\r\n", err)
		return
	}

	_L.LoggerInstance.InfoPrint(" > Total Account    : %d\r\n", allactive)
	_L.LoggerInstance.InfoPrint(" > Active Account   : %d\r\n", allactive-factive)
	_L.LoggerInstance.InfoPrint(" > Unactive Account : %d\r\n", factive)

	fmt.Scanf("%s\n")

	if !isSync {
		ths.ASyncChan <- 0
	}
}
