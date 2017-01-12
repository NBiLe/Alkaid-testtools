package statics

import (
	"sync"
	"time"

	"fmt"

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
	acc    map[int]*AccActiveStatic
	CTime  time.Time
	locker *sync.Mutex
}

// NewAccActiveController new account active controller
func NewAccActiveController() *AccActiveStaticController {
	ret := &AccActiveStaticController{}
	return ret.Init()
}

// Init init
func (ths *AccActiveStaticController) Init() *AccActiveStaticController {
	ths.acc = make(map[int]*AccActiveStatic)
	ths.locker = new(sync.Mutex)
	return ths
}

// Clear 清除统计结果
func (ths *AccActiveStaticController) Clear() {
	ths.Init()
}

// Put signed Base64 string
func (ths *AccActiveStaticController) Put(idx int, signer string, action string) {
	ths.locker.Lock()
	defer ths.locker.Unlock()
	lenAcc := len(ths.acc)
	if lenAcc == 0 {
		ths.CTime = time.Now()
	}
	accact := NewAccActiveStatic()
	accact.SetSignature(signer)
	accact.Index = idx
	accact.Counter = lenAcc + 1
	accact.CreateTime = ths.CTime
	accact.Action = action
	ths.acc[idx] = accact
}

func (ths *AccActiveStaticController) PrintAccs() {
	for i := 0; i < len(ths.acc); i++ {
		fmt.Printf("[%d] %+v\r\n", i, ths.acc[i])
	}
	// fmt.Printf("Accounts = \r\n%+v\r\n", ths.acc)
}

// SetTimeTicker set time ticker
func (ths *AccActiveStaticController) SetTimeTicker(index int, t int64, tick TimeTicker) {
	ths.locker.Lock()
	defer ths.locker.Unlock()
	if index >= 1000000 {
		acc, ok := ths.acc[index]
		if !ok {
			acc, _ = ths.acc[index%1000000]
			ths.Put(index, acc.SignerB64, acc.Action)
		}
		acc.SetTimeTick(tick, t)
	}
	ths.acc[index].SetTimeTick(tick, t)
}

// SetResult set result
func (ths *AccActiveStaticController) SetResult(index int, b bool) {
	ths.locker.Lock()
	defer ths.locker.Unlock()
	if index >= 1000000 {
		acc, ok := ths.acc[index]
		if !ok {
			acc, _ = ths.acc[index%1000000]
			ths.Put(index, acc.SignerB64, acc.Action)
		}
	}
	if b {
		ths.acc[index].Success = "T"
	} else {
		ths.acc[index].Success = "F"
	}
}

// Update2DB 统计结果入库
func (ths *AccActiveStaticController) Update2DB() error {
	ths.locker.Lock()
	defer ths.locker.Unlock()
	lenAcc := len(ths.acc)
	save := make([]*_db.TStaticsActiveAccount, lenAcc)
	idx := 0
	for _, itm := range ths.acc {
		save[idx] = &(itm.TStaticsActiveAccount)
		idx++
	}
	return _db.DataBaseInstance.Quary(_db.KeyAAStatics, _db.QtAddRecords, save)
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
