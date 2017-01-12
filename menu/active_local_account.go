package menu

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	_ac "github.com/fuyaocn/evaluatetools/appconf"
	_db "github.com/fuyaocn/evaluatetools/db"
	_L "github.com/fuyaocn/evaluatetools/log"
	_s "github.com/fuyaocn/evaluatetools/statics"
	_str "github.com/fuyaocn/evaluatetools/stellar"
)

// ActiveAccount 激活本地账户
type ActiveAccount struct {
	MenuSubItem
	baseAcc     []*_str.AccountInfo
	SubAccounts []*_str.AccountInfo
	factive     int64
	locker      *sync.Mutex
	offset      int64
	firstRowID  int64
	Level       int // 当前等级
	BaseLevel   int // 基础等级 第一级 之后将按照这个等级扩展
	DepthLevel  int // 深度 做多扩展几级
}

// InitMenu 初始化
func (ths *ActiveAccount) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.MenuSubItem.InitMenu(parent, key)
	ths.Exec = ths.execute
	ths.locker = new(sync.Mutex)

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
			ths.Level = 1
			ths.BaseLevel = _ac.ConfigInstance.GetActiveLevel()
			ths.DepthLevel = _ac.ConfigInstance.GetActiveDepth()
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
	retAcc := make([]*_db.TMainAccount, 0)
	opera := _db.DataBaseInstance.GetOperation(_db.KeyMainAccount).(*_db.MainAccountOperation)
	err := opera.GetEngine().Where("`index`<>?", 0).Asc("index").Find(&retAcc)
	if err != nil {
		return err
	}
	length := len(retAcc)
	ths.baseAcc = make([]*_str.AccountInfo, length)
	// network := _ac.ConfigInstance.GetNetwork()
	for i := 0; i < length; i++ {
		ths.baseAcc[i] = new(_str.AccountInfo)
		ths.baseAcc[i].Init(retAcc[i].AccountID, retAcc[i].SecertAddr)
	}
	return nil
}

func (ths *ActiveAccount) activeAccount() (err error) {
	_s.ActiveAccountStaticsInstance.Clear()
	network := _ac.ConfigInstance.GetNetwork()
	amount := _ac.ConfigInstance.GetStartBalance()

	wg := new(sync.WaitGroup)
	for idx, itm := range ths.baseAcc {
		wg.Add(1)
		go func(index int, baseacc *_str.AccountInfo, amt string, wt *sync.WaitGroup) {
			defer wt.Done()
			addr := ths.GetRecords()
			if addr == nil {
				return
			}
			atc := &_str.ActiveAccount{}
			err := baseacc.GetInfo(network, nil)
			if err != nil {
				_L.LoggerInstance.ErrorPrint(" > Get main account[%s] info has error : \r\n%+v\r\n", baseacc.Address, err)
				return
			}
			atc.Init(amt, baseacc.Secret, baseacc.GetNextSequence())
			b64 := atc.GetSignature(addr, network, nil)
			_s.ActiveAccountStaticsInstance.Put(index, b64, "active")
			base64 := make([]string, 0)
			base64 = append(base64, b64)
			atc.SendTransaction(index, "", nil, base64)
		}(idx, itm, amount, wg)
	}
	wg.Wait()

	err = _s.ActiveAccountStaticsInstance.Update2DB()
	if err != nil {
		_L.LoggerInstance.ErrorPrint(" > Update statics data to database has err : \r\n %+v\r\n", err)
	} else {
		_L.LoggerInstance.InfoPrint(" > Update statics data to database complete!")
	}
	return
}

// GetRecords 从数据库中读取需要激活的账户地址数组
func (ths *ActiveAccount) GetRecords() (ret []string) {
	ths.locker.Lock()
	defer ths.locker.Unlock()
	ths.SubAccounts = ths.GetRecordFromDB()
	if ths.SubAccounts != nil {
		ret = make([]string, 0)
		for _, itm := range ths.SubAccounts {
			ret = append(ret, itm.Address)
		}
	}
	return
}

func (ths *ActiveAccount) GetRecordFromDB() []*_str.AccountInfo {
	opera := _db.DataBaseInstance.GetAccountOperation()
	if ths.offset == 0 {
		ths.firstRowID = ths.getFirstRowID(opera)
		if ths.firstRowID == -1 {
			return nil
		}
		_L.LoggerInstance.DebugPrint("[AAMainController:GetRecordFromDB] get records firstRowID = %d\r\n", ths.firstRowID)
	}

	_L.LoggerInstance.DebugPrint("[AAMainController:GetRecordFromDB] get records offset = %d\r\n", ths.offset)

	ret, err := opera.GetCountRecords("F", int64(ths.BaseLevel), ths.offset)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("[AAMainController:GetRecordFromDB] get records has error : \r\n%+v\r\n", err)
		return nil
	}
	size := len(ret)
	if size == 0 {
		_L.LoggerInstance.DebugPrint("[AAMainController:GetRecordFromDB] get records length = 0\r\n")
		return nil
	}
	accs := make([]*_str.AccountInfo, size)
	for idx := 0; idx < size; idx++ {
		accs[idx] = &_str.AccountInfo{}
		accs[idx].Init(string(ret[idx]["account_id"]), string(ret[idx]["secert_addr"]))
	}
	lastid, _ := strconv.ParseInt(string(ret[size-1]["id"]), 10, 64)
	_L.LoggerInstance.DebugPrint("[AAMainController:GetRecordFromDB] get records LastRowID = %d\r\n", lastid)
	ths.offset = lastid - ths.firstRowID + 1
	return accs
}

func (ths *ActiveAccount) getFirstRowID(o *_db.OperationAccount) int64 {
	ret, err := o.GetCountRecords("F", 1, 0)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("[AAMainController:getFirstRowID] get first id has error : \r\n%+v\r\n", err)
		return -1
	}
	if len(ret) == 0 {
		return -1
	}
	retid, _ := strconv.ParseInt(string(ret[0]["id"]), 10, 64)
	return retid
}
