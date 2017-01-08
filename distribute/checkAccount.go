package distribute

import (
	"net/url"
	"strings"
	"sync"

	_ac "github.com/fuyaocn/evaluatetools/appconf"
	_db "github.com/fuyaocn/evaluatetools/db"
	_h "github.com/fuyaocn/evaluatetools/http"
	_L "github.com/fuyaocn/evaluatetools/log"
	_str "github.com/fuyaocn/evaluatetools/stellar"
	_b "github.com/stellar/go/build"
	"github.com/stellar/go/xdr"
)

// CAMainController 检查账户分发控制器
type CAMainController struct {
	AAMainController
	wt *sync.WaitGroup
}

// Init 初始化
func (ths *CAMainController) Init(l, bl, dl int, a *_str.AccountInfo) IController {
	ths.Executor.Init(l, 100, dl, a)
	ths.locker = &sync.Mutex{}
	ths.offset = 0
	ths.wt = new(sync.WaitGroup)
	return ths
}

// Execute 执行
func (ths *CAMainController) Execute() {
	// 先从数据库中读取需要激活的账户地址
	network := _ac.ConfigInstance.GetNetwork()
	for {
		accInfo := ths.GetRecordFromDB()
		if accInfo == nil {
			_L.LoggerInstance.InfoPrint(" > Read unactive account is complete \r\n")
			break
		}
		_L.LoggerInstance.DebugPrint(" **** Get accinfo \r\n %+v\r\n", accInfo)
		for _, acc := range accInfo {
			ths.wt.Add(3)
			go ths.checkAccInfo(acc, network, ths.wt)
		}
		ths.wt.Wait()
	}
}

func (ths *CAMainController) checkAccInfo(acc *_str.AccountInfo, flag string, wg *sync.WaitGroup) {
	err := acc.GetInfo(flag, wg)
	defer wg.Done()
	if err == nil {
		operaAcc := _db.DataBaseInstance.GetAccountOperation()
		sqlcmd := "update `t_account` set `active`=? where `account_id`=?"
		err = operaAcc.CmdExec(sqlcmd, "T", acc.Address)
		if err != nil {
			_L.LoggerInstance.ErrorPrint(" ### Update account info to database has error :\r\n > %+v\r\n", err)
		} else {
			ths.changeTrust(acc, flag, wg)
		}
		return
	}
	wg.Done()
	_L.LoggerInstance.ErrorPrint(" Check account info [%s], has error \r\n %+v\r\n", acc.Address, err)
}

func (ths *CAMainController) changeTrust(acc *_str.AccountInfo, flag string, wg *sync.WaitGroup) {
	defer wg.Done()
	tx := &_b.TransactionBuilder{}
	tx.Mutate(_b.SourceAccount{AddressOrSeed: acc.Address})
	tx.Mutate(_b.ChangeTrust(
		_b.Asset{
			Code:   _ac.ConfigInstance.GetMainCredit(),
			Issuer: _ac.ConfigInstance.GetMainIssuerID(),
			Native: false,
		},
		_b.MaxLimit,
		_b.SourceAccount{AddressOrSeed: acc.Address},
	))
	tx.Mutate(_b.Sequence{Sequence: acc.GetNextSequence()})
	if strings.ToLower(flag) == "live" {
		tx.Mutate(_str.NBiPublicNetwork)
	} else {
		tx.Mutate(_str.NBiTestNetwork)
	}
	tx.TX.Fee = xdr.Uint32(100)
	ret := tx.Sign(acc.Secret)
	base64, _ := ret.Base64()

	_L.LoggerInstance.DebugPrint("ChangeTrust Base64 \r\n > %s\r\n", base64)

	data := "tx=" + url.QueryEscape(base64)

	type Result struct {
		Hash string `json:"hash"`
	}
	result := &Result{}
	horizon := _ac.ConfigInstance.GetHorizon() + "/transactions"
	webhttp := &_h.WebController{}
	err := webhttp.HttpPostForm(horizon, data, result)
	if err != nil || len(result.Hash) == 0 {
		_L.LoggerInstance.ErrorPrint(" Account[%s] change trust has error \r\n %+v\r\n", acc.Address, err)
	}
}
