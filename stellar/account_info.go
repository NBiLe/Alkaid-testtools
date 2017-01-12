package stellar

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"fmt"

	_ac "github.com/fuyaocn/evaluatetools/appconf"
	_L "github.com/fuyaocn/evaluatetools/log"
)

// AssetInfo stellar asset info
type AssetInfo struct {
	Type    string `json:"asset_type"`
	Balance string `json:"balance"`
	Code    string `json:"asset_code"`
	Issuer  string `json:"issuer"`
}

// AccountInfo stellar account info
type AccountInfo struct {
	ID       string      `json:"id"`
	Sequence string      `json:"sequence"`
	Assets   []AssetInfo `json:"balances"`
	Address  string      `json:"address"`
	Status   int         `json:"status"`
	Balance  float64
	sequence uint64
	Secret   string
}

// Init set address
func (ths *AccountInfo) Init(id, secret string) {
	ths.Address = id
	ths.Secret = secret
}

// GetInfo get base information
func (ths *AccountInfo) GetInfo(flag string, wt *sync.WaitGroup) error {
	if wt != nil {
		defer wt.Done()
	}
	addr := ""
	if strings.ToLower(flag) == "live" {
		addr = _ac.ConfigInstance.GetHorizonLive()
	} else {
		addr = _ac.ConfigInstance.GetHorizonTest()
	}

	addr = addr + "/accounts/" + ths.Address

	resp, err := http.Get(addr)
	if err != nil {
		_L.LoggerInstance.ErrorPrint(" Http get '%s' has error : \r\n %+v\r\n", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		_L.LoggerInstance.ErrorPrint(" ** HTTP Response ERROR\r\n\tReadAll: %+v\r\n", err)
		return err
	}

	err = json.Unmarshal(body, ths)
	if err != nil {
		_L.LoggerInstance.ErrorPrint(" Unmarshal body has error : %+v\r\n", err)
		return err
	}
	if ths.Status == 0 {
		for _, itm := range ths.Assets {
			if itm.Type == "native" {
				ths.Balance, _ = strconv.ParseFloat(itm.Balance, 64)
				break
			}
		}

		ths.sequence, err = strconv.ParseUint(ths.Sequence, 10, 64)
		_L.LoggerInstance.DebugPrint(" Current Account info : %+v\r\n", ths)
		return err
	}
	return fmt.Errorf("Account is not exist")
}

// GetNextSequence get next sequence
func (ths *AccountInfo) GetNextSequence() uint64 {
	ths.sequence++
	return ths.sequence
}

// GetCurrentSequence get currnt sequence
func (ths *AccountInfo) GetCurrentSequence() uint64 {
	return ths.sequence
}

// GetResetSequence reset currnt sequence
func (ths *AccountInfo) GetResetSequence() {
	ths.sequence--
}
