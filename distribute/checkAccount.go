package distribute

import (
	"sync"

	_ac "github.com/fuyaocn/evaluatetools/appconf"
	_db "github.com/fuyaocn/evaluatetools/db"
	_L "github.com/fuyaocn/evaluatetools/log"
	_s "github.com/fuyaocn/evaluatetools/statics"
	_str "github.com/fuyaocn/evaluatetools/stellar"
)

// CAMainController 检查账户分发控制器
type CAMainController struct {
	AAMainController
	wt         *sync.WaitGroup
	Signatures []string
	cTrust     *_str.ChangeTrustList
}

// Init 初始化
func (ths *CAMainController) Init(l, bl, dl int, a *_str.AccountInfo) IController {
	ths.Executor.Init(l, bl, dl, a)
	ths.locker = &sync.Mutex{}
	ths.offset = 0
	ths.wt = new(sync.WaitGroup)
	ths.Signatures = make([]string, 0)
	ths.cTrust = new(_str.ChangeTrustList)
	return ths
}

// Execute 执行
func (ths *CAMainController) Execute() {
	// 先从数据库中读取需要激活的账户地址
	network := _ac.ConfigInstance.GetNetwork()
	_s.ActiveAccountStaticsInstance.Clear()
	index := 0
	for {
		accInfo := ths.GetRecordFromDB()
		if accInfo == nil {
			_L.LoggerInstance.InfoPrint(" > Read unactive account is complete \r\n")
			break
		}
		// _L.LoggerInstance.DebugPrint(" **** Get accinfo \r\n %+v\r\n", accInfo)
		// go cTrust.GetSignature(accInfo,network,ths.wt)
		for _, acc := range accInfo {
			ths.wt.Add(2)
			go ths.checkAccInfo(acc, network, ths.wt)
		}
		ths.wt.Wait()
		ths.wt.Add(2)
		go ths.getSignatures(index, accInfo, network, ths.wt)
		ths.wt.Wait()
		_L.LoggerInstance.InfoPrint(" > [%d] Sending change trust signatures... \r\n", index)
		ths.cTrust.SendTransaction(index*ths.BaseLevel, network, nil, ths.Signatures)
		ths.Signatures = make([]string, 0)
		ths.offset = 0
		index++
	}

	err := _s.ActiveAccountStaticsInstance.Update2DB()
	if err != nil {
		_L.LoggerInstance.ErrorPrint(" > Update statics data to database has err : \r\n %+v\r\n", err)
	} else {
		_L.LoggerInstance.InfoPrint(" > Update statics data to database complete!")
	}

	// _s.ActiveAccountStaticsInstance.Clear()
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
		}
		return
	}
	// wg.Done()
	_L.LoggerInstance.ErrorPrint(" Check account info [%s], has error \r\n %+v\r\n", acc.Address, err)
}

func (ths *CAMainController) getSignatures(idx int, accs []*_str.AccountInfo, flag string, wg *sync.WaitGroup) {
	defer wg.Done()
	b64s := ths.cTrust.GetSignature(accs, flag, wg)
	for i, b64 := range b64s {
		ths.Signatures = append(ths.Signatures, b64)
		index := i + idx*ths.BaseLevel
		_s.ActiveAccountStaticsInstance.Put(index, b64, "trustline")
	}
}

// func (ths *CAMainController) changeTrust(acc *_str.AccountInfo, flag string, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	tx := &_b.TransactionBuilder{}
// 	tx.Mutate(_b.SourceAccount{AddressOrSeed: acc.Address})
// 	tx.Mutate(_b.ChangeTrust(
// 		_b.Asset{
// 			Code:   _ac.ConfigInstance.GetMainCredit(),
// 			Issuer: _ac.ConfigInstance.GetMainIssuerID(),
// 			Native: false,
// 		},
// 		_b.MaxLimit,
// 		_b.SourceAccount{AddressOrSeed: acc.Address},
// 	))
// 	tx.Mutate(_b.Sequence{Sequence: acc.GetNextSequence()})
// 	if strings.ToLower(flag) == "live" {
// 		tx.Mutate(_str.NBiPublicNetwork)
// 	} else {
// 		tx.Mutate(_str.NBiTestNetwork)
// 	}
// 	tx.TX.Fee = xdr.Uint32(100)
// 	ret := tx.Sign(acc.Secret)
// 	base64, _ := ret.Base64()

// 	_L.LoggerInstance.DebugPrint("ChangeTrust Base64 \r\n > %s\r\n", base64)

// 	data := "tx=" + url.QueryEscape(base64)

// 	type Result struct {
// 		Hash string `json:"hash"`
// 	}
// 	result := &Result{}
// 	horizon := _ac.ConfigInstance.GetHorizon() + "/transactions"
// 	webhttp := &_h.WebController{}
// 	err := webhttp.HttpPostForm(horizon, data, result)
// 	if err != nil || len(result.Hash) == 0 {
// 		_L.LoggerInstance.ErrorPrint(" Account[%s] change trust has error \r\n %+v\r\n", acc.Address, err)
// 	}
// }
