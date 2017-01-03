package main

import (
	"fmt"
	"runtime"
	"time"

	_ac "github.com/fuyaocn/evaluatetools/appconf"
	_db "github.com/fuyaocn/evaluatetools/db"
	_L "github.com/fuyaocn/evaluatetools/log"
	_m "github.com/fuyaocn/evaluatetools/menu"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	_m.MainMenuInstace.SetLanguageType(_m.LEnglish)

	_L.LoggerInstance = _L.NewLoggerInstance(fmt.Sprintf("NBiTester.%s", time.Now().Format("2006-01-02_15.04.05.000")))
	_L.LoggerInstance.OpenDebug = true
	_L.LoggerInstance.SetLogFunCallDepth(4)
	_L.LoggerInstance.InfoPrint("Read database configurations ...\r\n")

	_ac.ConfigInstance = _ac.NewConfigController()
	dbConf := _ac.ConfigInstance.GetDateBaseConf()
	if dbConf == nil {
		panic(0)
	}
	_L.LoggerInstance.InfoPrint("Init database ...\r\n")
	_db.DataBaseInstance = _db.CreateDBInstance(dbConf)

	_m.MainMenuInstace.ExecuteFunc(true)
}
