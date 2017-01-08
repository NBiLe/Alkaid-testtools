package statics

import (
	"time"

	_db "github.com/fuyaocn/evaluatetools/db"
)

const (
	StartTime TimeTicker = 1
	EndTime   TimeTicker = 2
)

// TimeTicker statics time tick define
type TimeTicker int

// ActiveAccountStaticsInstance 唯一实例
var ActiveAccountStaticsInstance *AccActiveStaticController

// AccActiveStaticController 激活账户全局唯一实例
type AccActiveStaticController struct {
	acc   []*AccActiveStatic
	CTime time.Time
}

// NewAccActiveController new account active controller
func NewAccActiveController() *AccActiveStaticController {
	ret := &AccActiveStaticController{}
	return ret.Init()
}

// Init init
func (ths *AccActiveStaticController) Init() *AccActiveStaticController {
	ths.acc = make([]*AccActiveStatic, 0)
	return ths
}

// Clear 清除统计结果
func (ths *AccActiveStaticController) Clear() {
	ths.Init()
}

// Put signed Base64 string
func (ths *AccActiveStaticController) Put(signer string) {
	lenAcc := len(ths.acc)
	if lenAcc == 0 {
		ths.CTime = time.Now()
	}
	accact := NewAccActiveStatic()
	accact.SetSignature(signer)
	accact.Counter = lenAcc + 1
	accact.CreateTime = ths.CTime
	ths.acc = append(ths.acc, accact)
}

// SetTimeTicker set time ticker
func (ths *AccActiveStaticController) SetTimeTicker(index int, t int64, tick TimeTicker) {
	ths.acc[index].SetTimeTick(tick, t)
}

// SetResult set result
func (ths *AccActiveStaticController) SetResult(index int, b bool) {
	if b {
		ths.acc[index].Success = "T"
	} else {
		ths.acc[index].Success = "F"
	}
}

// Update2DB 统计结果入库
func (ths *AccActiveStaticController) Update2DB() error {
	lenAcc := len(ths.acc)
	save := make([]*_db.TStaticsActiveAccount, lenAcc)
	for i := 0; i < lenAcc; i++ {
		save[i] = &(ths.acc[i].TStaticsActiveAccount)
	}
	err := _db.DataBaseInstance.Quary(_db.KeyAAStatics, _db.QtAddRecords, save)
	return err
}

// AccActiveStatic account active statics
type AccActiveStatic struct {
	_db.TStaticsActiveAccount
}

// NewAccActiveStatic new account active statics
func NewAccActiveStatic() *AccActiveStatic {
	ret := &AccActiveStatic{}
	return ret.Init()
}

// Init init
func (ths *AccActiveStatic) Init() *AccActiveStatic {
	ths.StartTimeTick = -1
	ths.EndTimeTick = -1
	return ths
}

// SetTimeTick set time ticker
func (ths *AccActiveStatic) SetTimeTick(tick TimeTicker, t int64) {
	if tick == StartTime {
		ths.StartTimeTick = t
	} else if tick == EndTime {
		ths.EndTimeTick = t
	}
}

// SetSignature set signature
func (ths *AccActiveStatic) SetSignature(s string) {
	ths.SignerB64 = s
}
