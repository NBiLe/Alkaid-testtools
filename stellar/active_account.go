package stellar

import (
	"net/url"
	"strings"
	"sync"

	_ac "github.com/fuyaocn/evaluatetools/appconf"
	_web "github.com/fuyaocn/evaluatetools/http"
	_L "github.com/fuyaocn/evaluatetools/log"
	_s "github.com/fuyaocn/evaluatetools/statics"
	_b "github.com/stellar/go/build"
	_kp "github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
)

// ActiveAccount 激活
type ActiveAccount struct {
	Address  []string
	key      string
	amount   string
	sequence uint64
}

// Init 初始化
func (ths *ActiveAccount) Init(amount, skey string, sequence uint64) {
	ths.key = skey
	ths.amount = amount
	ths.sequence = sequence
}

// GetSignature get signature
func (ths *ActiveAccount) GetSignature(addr []string, flag string, wt *sync.WaitGroup) string {
	if wt != nil {
		defer wt.Done()
	}

	kp, _ := _kp.Parse(ths.key)

	tx := &_b.TransactionBuilder{}
	tx.Mutate(_b.SourceAccount{AddressOrSeed: kp.Address()})
	for _, itm := range addr {
		tx.Mutate(_b.CreateAccount(
			_b.Destination{AddressOrSeed: itm},
			_b.NativeAmount{Amount: ths.amount},
			_b.SourceAccount{AddressOrSeed: kp.Address()},
		))
	}
	tx.Mutate(_b.Sequence{Sequence: ths.sequence})
	_L.LoggerInstance.DebugPrint("Signer sequence : %d\r\n", ths.sequence)
	ths.sequence++
	if strings.ToLower(flag) == "live" {
		tx.Mutate(NBiPublicNetwork)
	} else {
		tx.Mutate(NBiTestNetwork)
	}

	tx.TX.Fee = xdr.Uint32(100 * len(addr))
	ret := tx.Sign(ths.key)
	base64, _ := ret.Base64()
	return base64
}

// SendTransaction send transaction
func (ths *ActiveAccount) SendTransaction(index int, flag string, wt *sync.WaitGroup, b64 []string) (err error) {
	if wt != nil {
		defer wt.Done()
	}

	type ExtrasData struct {
		ResultXdr string `json:"result_xdr"`
	}

	type Result struct {
		Hash string `json:"hash"`
		// Detail string `json:"detail"`
		Extras *ExtrasData `json:"extras"`
	}

	lenb64 := len(b64)
	slicehttp := make([]*_web.SocketHttp, lenb64)
	horizon := _ac.ConfigInstance.GetHorizon() + "/transactions"

	gw := new(sync.WaitGroup)

	for idx := 0; idx < lenb64; idx++ {
		data := "tx=" + url.QueryEscape(b64[idx])
		// http := &_web.WebController{}
		slicehttp[idx] = new(_web.SocketHttp)
		err = slicehttp[idx].Init(horizon)
		if err != nil {
			return err
		}
		// err := http.HttpPostForm(horizon, data, ret)
		err = slicehttp[idx].PostForm(data, _ac.ConfigInstance.GetHorizonHeader())
		if err != nil {
			return err
		}
		_s.ActiveAccountStaticsInstance.SetTimeTicker(index, slicehttp[idx].StartSend, _s.StartTime)
		gw.Add(1)
		go func(w *sync.WaitGroup, subIndex int, http *_web.SocketHttp) {
			defer w.Done()
			ret := &Result{}
			err = http.Response(ret)
			_s.ActiveAccountStaticsInstance.SetTimeTicker(subIndex, http.CompleteSend, _s.EndTime)
			if err == nil {
				if ret.Extras == nil && ret.Hash != "" {
					http.Result = "Success"
					_s.ActiveAccountStaticsInstance.SetResult(subIndex, true)
					return
				}
			}
			http.Result = "Failure"
			_s.ActiveAccountStaticsInstance.SetResult(subIndex, false)
			if ret.Extras != nil {
				tret := &xdr.TransactionResult{}
				tret.Scan(ret.Extras.ResultXdr)
				ret.Extras.ResultXdr = tret.Result.Code.String()
				_L.LoggerInstance.ErrorPrint(" ##[%d]## Create account transaction is fail!! ###\r\n ### error : %v\r\n ### Detail : %s\r\n", subIndex, err, ret.Extras.ResultXdr)
			} else {
				_L.LoggerInstance.ErrorPrint(" ##[%d]## Create account transaction is fail!! ###\r\n ### error : %v\r\n ### Detail : Timeout\r\n", subIndex, err)
			}
		}(gw, index+idx*1000000, slicehttp[idx])
	}
	gw.Wait()

	gw.Add(1)
	go ths.SaveStatic(slicehttp, gw)
	gw.Wait()
	return
}

func (ths *ActiveAccount) SaveStatic(h []*_web.SocketHttp, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	for idx, itm := range h {
		st := float64(itm.CompleteSend-itm.StartSend) * 1.0 / 1000000000.0
		_L.LoggerInstance.InfoPrint("[%05d]\t[%s]\t[TimeSpan:%.5f s]\r\n",
			idx, itm.Result, st)
	}
}
