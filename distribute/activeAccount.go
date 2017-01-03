package distribute

import (
	"strconv"
	"sync"

	_ac "github.com/fuyaocn/evaluatetools/appconf"
	_db "github.com/fuyaocn/evaluatetools/db"
	_L "github.com/fuyaocn/evaluatetools/log"
	_str "github.com/fuyaocn/evaluatetools/stellar"
	_kp "github.com/stellar/go/keypair"
)

// AAMainController 激活账户分发控制器
type AAMainController struct {
	Executor
	locker     *sync.Mutex
	offset     int64
	firstRowID int64
}

// NewAAMainController 创建激活账户分发控制器
func NewAAMainController() *AAMainController {
	ret := new(AAMainController)
	sk := _ac.ConfigInstance.GetMainSecretKey()
	kp, err := _kp.Parse(sk)
	if err != nil {
		_L.LoggerInstance.ErrorPrint(" In 'skey.conf', secret key is not valid! \r\n err : %+v\r\n", err)
		return nil
	}
	acc := new(_str.AccountInfo)
	acc.Init(kp.Address(), sk)
	ret.Init(1, _ac.ConfigInstance.GetActiveLevel(), _ac.ConfigInstance.GetActiveDepth(), acc)
	return ret
}

// Init 初始化
func (ths *AAMainController) Init(l, bl, dl int, a *_str.AccountInfo) IController {
	ths.Executor.Init(l, bl, dl, a)
	ths.locker = &sync.Mutex{}
	ths.offset = 0
	return ths
}

// Execute 执行
func (ths *AAMainController) Execute() {
	var err error
	// 先从数据库中读取需要激活的账户地址

	atc := &_str.ActiveAccount{}

	network := _ac.ConfigInstance.GetNetwork()
	err = ths.BaseAccount.GetInfo(network, nil)
	if err != nil {
		_L.LoggerInstance.ErrorPrint(" > Get [%s] info has error : \r\n%+v\r\n", ths.BaseAccount.Address, err)
		return
	}
	atc.Init(ths.getAmount(), ths.BaseAccount.Secret, ths.BaseAccount.GetNextSequence())
	base64 := make([]string, 0)
	for {
		addr := ths.GetRecords()
		if addr == nil {
			_L.LoggerInstance.InfoPrint(" > Read unactive account is null, execute post!! \r\n")
			break
		}

		_L.LoggerInstance.DebugPrint("\r\n > Read local database address : \r\n  >> %+v\r\n", addr)

		b64 := atc.GetSignature(addr, network, nil)
		_L.LoggerInstance.DebugPrint("get base 64 : \r\n%s\r\n", b64)
		base64 = append(base64, b64)
	}

	atc.SendTransaction("", nil, base64)
}

// GetRecords 从数据库中读取需要激活的账户地址数组
func (ths *AAMainController) GetRecords() (ret []string) {
	ths.locker.Lock()
	defer ths.locker.Unlock()
	ths.SubAccounts = ths.getRecordFromDB()
	if ths.SubAccounts != nil {
		ret = make([]string, 0)
		for _, itm := range ths.SubAccounts {
			ret = append(ret, itm.Address)
		}
	}
	return
}

func (ths *AAMainController) getRecordFromDB() []*_str.AccountInfo {
	opera := _db.DataBaseInstance.GetAccountOperation()
	if ths.offset == 0 {
		ths.firstRowID = ths.getFirstRowID(opera)
	}
	ret, err := opera.GetCountRecords("F", int64(ths.BaseLevel), ths.offset)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("[AAMainController:getRecordFromDB] get records has error : \r\n%+v\r\n", err)
		return nil
	}
	size := len(ret)
	if size == 0 {
		return nil
	}
	accs := make([]*_str.AccountInfo, size)
	for idx := 0; idx < size; idx++ {
		accs[idx] = &_str.AccountInfo{}
		accs[idx].Init(string(ret[idx]["account_id"]), string(ret[idx]["secert_addr"]))
	}
	lastid, _ := strconv.ParseInt(string(ret[size-1]["id"]), 10, 64)
	ths.offset = lastid - ths.firstRowID + 1
	return accs
}

func (ths *AAMainController) getFirstRowID(o *_db.OperationAccount) int64 {
	ret, err := o.GetCountRecords("", 1, 0)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("[AAMainController:getFirstRowID] get first id has error : \r\n%+v\r\n", err)
		return -1
	}
	retid, _ := strconv.ParseInt(string(ret[0]["id"]), 10, 64)
	return retid
}

func (ths *AAMainController) getAmount() string {
	// if ths.Level == ths.DepthLevel-1 {
	return _ac.ConfigInstance.GetStartBalance()
	// }
	// saveBalance, _ := strconv.ParseFloat(_ac.ConfigInstance.GetStartBalance(), 64)
	// amount := (ths.BaseAccount.Balance - saveBalance) / float64(ths.BaseLevel)
	// return fmt.Sprintf("%f", amount)
}
