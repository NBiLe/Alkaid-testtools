package menu

import (
	"fmt"
	"strings"
	"sync"

	_ac "github.com/fuyaocn/evaluatetools/appconf"
	_db "github.com/fuyaocn/evaluatetools/db"
	_dis "github.com/fuyaocn/evaluatetools/distribute"
	_L "github.com/fuyaocn/evaluatetools/log"
	_str "github.com/fuyaocn/evaluatetools/stellar"
	_kp "github.com/stellar/go/keypair"
)

// ActiveAccount 激活本地账户
type ActiveAccount struct {
	MenuSubItem
	baseAcc *_str.AccountInfo
	factive int64
}

// InitMenu 初始化
func (ths *ActiveAccount) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.MenuSubItem.InitMenu(parent, key)
	ths.Exec = ths.execute

	ths.title = []string{
		LEnglish: "Active local stellar accounts",
	}
	return ths
}

func (ths *ActiveAccount) execute(isSync bool) {
	_L.LoggerInstance.Info(" ** Active stellar accounts ** \r\n")
	_L.LoggerInstance.InfoPrint(" > Current 'active.conf' -> start balance = %s\r\n", _ac.ConfigInstance.GetStartBalance())
	var input string
	ths.factive = 0
	err := _db.DataBaseInstance.GetAccountCount("F", false, &ths.factive)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("Get local stellar account count from database has error :\r\n%+v\r\n", err)
	} else {
		_L.LoggerInstance.InfoPrint(" > Unactive account count = %d\r\n", ths.factive)
		err = ths.getBaseAccountInfo()
		if err == nil {
			fmt.Printf("\r\n Are you sure to active those accounts?(y/n) : ")
			fmt.Scanf("%s\n", &input)
			if strings.ToLower(input) == "y" {
				ths.activeAccount()
			}
		}
	}

	if !isSync {
		ths.ASyncChan <- 0
	}
}

func (ths *ActiveAccount) getBaseAccountInfo() error {
	sk := _ac.ConfigInstance.GetMainSecretKey()
	kp, err := _kp.Parse(sk)
	if err != nil {
		_L.LoggerInstance.ErrorPrint(" In 'skey.conf', secret key is not valid! \r\n err : %+v\r\n", err)
		return err
	}
	ths.baseAcc = new(_str.AccountInfo)
	ths.baseAcc.Init(kp.Address(), sk)

	wt := new(sync.WaitGroup)

	fmt.Printf(" > reading base account info, please wait ...")
	wt.Add(1)
	go ths.baseAcc.GetInfo(_ac.ConfigInstance.GetNetwork(), wt)
	wt.Wait()
	fmt.Print("\r")
	_L.LoggerInstance.InfoPrint(" > Base account balance = %f\r\n", ths.baseAcc.Balance)
	return nil
}

func (ths *ActiveAccount) activeAccount() (err error) {
	ret := _dis.NewAAMainController()
	ret.Execute()
	return
}
