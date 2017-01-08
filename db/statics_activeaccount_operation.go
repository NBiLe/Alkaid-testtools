package db

import (
	"fmt"
	"sync"

	_L "github.com/fuyaocn/evaluatetools/log"
	"github.com/go-xorm/xorm"
)

// OperationActiveAccStatic operation for account active statics
type OperationActiveAccStatic struct {
	engine *xorm.Engine
	locker *sync.Mutex
}

// Init 初始化
func (ths *OperationActiveAccStatic) Init(e *xorm.Engine) {
	ths.locker = &sync.Mutex{}
	ths.engine = e
}

// GetKey get key string
func (ths *OperationActiveAccStatic) GetKey() string {
	return KeyAAStatics
}

// Quary quary exeute
func (ths *OperationActiveAccStatic) Quary(qtype int, v ...interface{}) (err error) {
	if qtype != QtClearAllRecord && (v == nil || len(v) < 1) {
		return fmt.Errorf("[OperationActiveAccStatic:Quary] Quary parameter 'v' is not be null")
	}

	switch qtype {
	case QtAddRecord:
		ths.locker.Lock()
		_, err = ths.engine.InsertOne(v[0])
		ths.locker.Unlock()
	case QtAddRecords:
		ths.locker.Lock()
		_, err = ths.engine.Insert(v...)
		ths.locker.Unlock()
	case QtClearAllRecord:
		ths.locker.Lock()
		// err = ths.clearStellarAccount()
		ths.locker.Unlock()
	default:
		_L.LoggerInstance.ErrorPrint("Unknown quary type [%d]\r\n", qtype)
	}

	return
}
