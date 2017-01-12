package menu

import (
	"fmt"
	"net/url"
	"strings"
	"sync"

	"strconv"

	_ac "github.com/fuyaocn/evaluatetools/appconf"
	_db "github.com/fuyaocn/evaluatetools/db"
	_h "github.com/fuyaocn/evaluatetools/http"
	_L "github.com/fuyaocn/evaluatetools/log"
	_str "github.com/fuyaocn/evaluatetools/stellar"
	_b "github.com/stellar/go/build"
	_kp "github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
)

// CreateMainAccount 创建主账户
type CreateMainAccount struct {
	MenuSubItem
}

// InitMenu 初始化
func (ths *CreateMainAccount) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.MenuSubItem.InitMenu(parent, key)
	ths.Exec = ths.execute

	ths.title = []string{
		LEnglish: "Create main stellar account",
	}
	ths.infoStrings = []map[int]string{
		LEnglish: map[int]string{
			NumberOfLocalAccount: " > Please enter the number of main accounts that need to be created (If you enter a number less than or equal to 0, it will end this operation) :",
		},
	}
	return ths
}

func (ths *CreateMainAccount) execute(isSync bool) {
	_L.LoggerInstance.Info(" ** Create main stellar account ** \r\n")
	// 先得到数据库里面有多少主账户
	opera := _db.DataBaseInstance.GetOperation(_db.KeyMainAccount)
	mainCnt := int64(0)
	err := opera.Quary(_db.QtGetCount, &mainCnt)
	if err == nil {
		root := ths.getRootAccInfo()
		if root != nil {
			_L.LoggerInstance.InfoPrint(" > Root account in balance : %f\r\n", root.Balance)
			_L.LoggerInstance.InfoPrint(" > Main account in database : %d\r\n", mainCnt)
			var input string
			fmt.Printf("\r\n" + ths.infoStrings[ths.languageIndex][NumberOfLocalAccount])
			fmt.Scanf("%s\n", &input)

			num, bNum := IsNumber(input)

			if bNum {
				if num > 0 {
					ths.run(num, root)
				}
			}
		}
	} else {
		_L.LoggerInstance.ErrorPrint("Get main account from db has error : \r\n %+v\r\n", err)
	}

	if !isSync {
		ths.ASyncChan <- 0
	}
}

func (ths *CreateMainAccount) run(n int, r *_str.AccountInfo) {
	wait := new(sync.WaitGroup)
	group := n / 100
	left := n % 100
	current := 0
	sb := _ac.ConfigInstance.GetMainStartBalance()
	gindex := ths.getMaxIndexFromDB()
	if gindex == -1 {
		return
	}
	gindex++

	for {
		if left > 0 {
			current = left
			left = 0
		} else if group > 0 {
			current = 100
			group--
		} else {
			break
		}
		wait.Add(1)
		ths.createGroup(sb, r, current, gindex, wait)
		wait.Wait()

		n = n - current
		gindex += current
		if n <= 0 {
			break
		}
	}
}

func (ths *CreateMainAccount) getMaxIndexFromDB() int {
	opera := _db.DataBaseInstance.GetOperation(_db.KeyMainAccount).(*_db.MainAccountOperation)
	sql := "select max(index) from t_main_account"
	ret, err := opera.CmdQuery(sql)
	if err == nil {
		val, ok := ret[0]["max"]
		if ok {
			ret, err := strconv.Atoi(string(val))
			if err == nil {
				return ret
			}
		} else {
			return 0
		}
	}
	_L.LoggerInstance.ErrorPrint(" Get main account Index(max) has error : \r\n%+v\r\n", err)
	return -1
}

func (ths *CreateMainAccount) getRootAccInfo() *_str.AccountInfo {
	ret := &_str.AccountInfo{}
	ret.Init(_ac.ConfigInstance.GetRootAccountID(), _ac.ConfigInstance.GetRootSecretKey())
	err := ret.GetInfo(_ac.ConfigInstance.GetNetwork(), nil)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("Get root account info has error :\r\n%+v\r\n", err)
		return nil
	}
	return ret
}

func (ths *CreateMainAccount) createGroup(sbalance string, r *_str.AccountInfo, cnt int, gindex int, wg *sync.WaitGroup) {
	defer wg.Done()

	accs := ths.getNewAccs(cnt, gindex)
	if accs == nil {
		return
	}

	tx := &_b.TransactionBuilder{}
	tx.Mutate(_b.SourceAccount{AddressOrSeed: r.Address})
	for _, itm := range accs {
		tx.Mutate(_b.CreateAccount(
			_b.Destination{AddressOrSeed: itm.AccountID},
			_b.NativeAmount{Amount: sbalance},
			_b.SourceAccount{AddressOrSeed: r.Address},
		))
	}
	tx.Mutate(_b.Sequence{Sequence: r.GetNextSequence()})
	if strings.ToLower(_ac.ConfigInstance.GetNetwork()) == "live" {
		tx.Mutate(_str.NBiPublicNetwork)
	} else {
		tx.Mutate(_str.NBiTestNetwork)
	}

	tx.TX.Fee = xdr.Uint32(100 * cnt)
	ret := tx.Sign(r.Secret)
	base64, _ := ret.Base64()

	type ExtrasData struct {
		ResultXdr string `json:"result_xdr"`
	}

	type Result struct {
		Hash   string      `json:"hash"`
		Extras *ExtrasData `json:"extras"`
	}

	web := &_h.WebController{}
	horizon := _ac.ConfigInstance.GetHorizon() + "/transactions"
	data := "tx=" + url.QueryEscape(base64)
	srlt := &Result{}
	err := web.HttpPostForm(horizon, data, srlt)
	if err == nil {
		if srlt.Extras != nil {
			tret := &xdr.TransactionResult{}
			tret.Scan(srlt.Extras.ResultXdr)
			srlt.Extras.ResultXdr = tret.Result.Code.String()
			ths.setResultFalse(accs)
			_L.LoggerInstance.ErrorPrint(" ### Create main account transaction is fail!! ###\r\n ### error : %v\r\n ### Detail : %s\r\n", err, srlt.Extras.ResultXdr)
			return
		} else {
			b, _ := strconv.ParseFloat(sbalance, 64)
			ths.setResultBalance(accs, b)
			_L.LoggerInstance.InfoPrint(" ### Create main account [%d] success!\r\n", cnt)
		}
	} else {
		ths.setResultFalse(accs)
		_L.LoggerInstance.ErrorPrint(" ### Create main account transaction is fail!! ###\r\n ### error : %v\r\n ### Detail : Timeout\r\n", err)
	}
	ths.add2DB(accs, nil)
}

func (ths *CreateMainAccount) setResultBalance(src []*_db.TMainAccount, b float64) {
	for _, itm := range src {
		itm.Balance = b
	}
}

func (ths *CreateMainAccount) setResultFalse(src []*_db.TMainAccount) {
	for _, itm := range src {
		itm.Success = "F"
	}
}

func (ths *CreateMainAccount) add2DB(src []*_db.TMainAccount, w *sync.WaitGroup) {
	if w != nil {
		defer w.Done()
	}

	opera := _db.DataBaseInstance.GetOperation(_db.KeyMainAccount)
	err := opera.Quary(_db.QtAddRecords, src)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("[CreateMainAccount:Execute] Add main account to database has error \r\n%+v\r\n", err)
	}
}

func (ths *CreateMainAccount) getNewAccs(n int, gindex int) (ret []*_db.TMainAccount) {
	if n <= 0 {
		return nil
	}
	ret = make([]*_db.TMainAccount, n)
	for i := 0; i < n; i++ {
		ret[i] = ths.newAcc(i + gindex)
		_L.LoggerInstance.DebugPrint(" > Addr = %s\r\n > Skey = %s\r\n", ret[i].AccountID, ret[i].SecertAddr)
	}
	return
}

func (ths *CreateMainAccount) newAcc(idx int) *_db.TMainAccount {
	full, err := _kp.Random()
	if err != nil {
		_L.LoggerInstance.ErrorPrint("[CreateMainAccount:Execute] Create keypair has error \r\n%+v\r\n", err)
		return nil
	}
	return &_db.TMainAccount{
		Index:      idx,
		AccountID:  full.Address(),
		SecertAddr: full.Seed(),
		Balance:    0,
		Success:    "T",
	}
}
