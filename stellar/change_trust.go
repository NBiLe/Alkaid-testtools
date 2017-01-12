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
	"github.com/stellar/go/xdr"
)

// ChangeTrustList 检查有效性并且信任Issuer
type ChangeTrustList struct {
	Address  []string
	key      string
	sequence uint64
}

// Init 初始化
func (ths *ChangeTrustList) Init() {
}

// GetSignature get signature
func (ths *ChangeTrustList) GetSignature(accinfos []*AccountInfo, flag string, wt *sync.WaitGroup) (retB64 []string) {
	if wt != nil {
		defer wt.Done()
	}
	assetCode := _ac.ConfigInstance.GetMainCredit()
	assetIssuer := _ac.ConfigInstance.GetMainIssuerID()
	lenAcc := len(accinfos)
	retB64 = make([]string, 0)
	for i := 0; i < lenAcc; i++ {
		acc := accinfos[i]
		if acc.Status != 0 {
			_L.LoggerInstance.ErrorPrint(" > Account [%s] is not actived, can not set trust line.\r\n", acc.Address)
			continue
		}
		tx := &_b.TransactionBuilder{}
		tx.Mutate(_b.SourceAccount{AddressOrSeed: acc.Address})
		tx.Mutate(_b.ChangeTrust(
			_b.Asset{
				Code:   assetCode,
				Issuer: assetIssuer,
				Native: false,
			},
			_b.MaxLimit,
			_b.SourceAccount{AddressOrSeed: acc.Address},
		))
		tx.Mutate(_b.Sequence{Sequence: acc.GetNextSequence()})
		if strings.ToLower(flag) == "live" {
			tx.Mutate(NBiPublicNetwork)
		} else {
			tx.Mutate(NBiTestNetwork)
		}
		tx.TX.Fee = xdr.Uint32(100)
		ret := tx.Sign(acc.Secret)
		base64, _ := ret.Base64()
		retB64 = append(retB64, base64)
	}
	return
}

// SendTransaction send transaction
func (ths *ChangeTrustList) SendTransaction(index int, flag string, wt *sync.WaitGroup, b64 []string) (err error) {
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
		// _s.ActiveAccountStaticsInstance.PrintAccs()
		// fmt.Printf("index + idx = %d\r\n", index+idx)
		_s.ActiveAccountStaticsInstance.SetTimeTicker(index+idx, slicehttp[idx].StartSend, _s.StartTime)
		gw.Add(1)
		go func(w *sync.WaitGroup, index int, http *_web.SocketHttp) {
			defer w.Done()
			ret := &Result{}
			err = http.Response(ret)
			_s.ActiveAccountStaticsInstance.SetTimeTicker(index, http.CompleteSend, _s.EndTime)
			if err == nil {
				if ret.Extras == nil && ret.Hash != "" {
					http.Result = "Success"
					_s.ActiveAccountStaticsInstance.SetResult(index, true)
					return
				}
			}
			http.Result = "Failure"
			_s.ActiveAccountStaticsInstance.SetResult(index, false)
			if ret.Extras != nil {
				tret := &xdr.TransactionResult{}
				tret.Scan(ret.Extras.ResultXdr)
				ret.Extras.ResultXdr = tret.Result.Code.String()
				_L.LoggerInstance.ErrorPrint(" ### Create change trust transaction is fail!! ###\r\n ### error : %v\r\n ### Detail : %s\r\n", err, ret.Extras.ResultXdr)
			} else {
				_L.LoggerInstance.ErrorPrint(" ### Create change trust transaction is fail!! ###\r\n ### error : %v\r\n ### Detail : Timeout\r\n", err)
			}
		}(gw, index+idx, slicehttp[idx])
	}
	gw.Wait()

	gw.Add(1)
	go ths.SaveStatic(slicehttp, gw)
	gw.Wait()
	return
}

func (ths *ChangeTrustList) SaveStatic(h []*_web.SocketHttp, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	for idx, itm := range h {
		st := float64(itm.CompleteSend-itm.StartSend) * 1.0 / 1000000000.0
		_L.LoggerInstance.InfoPrint("[%05d]\t[%s]\t[TimeSpan:%.5f s]\r\n",
			idx, itm.Result, st)
	}
}
