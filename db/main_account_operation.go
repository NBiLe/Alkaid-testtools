package db

import (
	"fmt"
	"sync"

	_L "github.com/fuyaocn/evaluatetools/log"
	"github.com/go-xorm/xorm"
)

// MainAccountOperation 主账号操作
type MainAccountOperation struct {
	engine *xorm.Engine
	locker *sync.Mutex
}

// Init 初始化
func (ths *MainAccountOperation) Init(e *xorm.Engine) {
	ths.locker = &sync.Mutex{}
	ths.engine = e
}

// GetKey get key string
func (ths *MainAccountOperation) GetKey() string {
	return KeyMainAccount
}

// Quary quary exeute
func (ths *MainAccountOperation) Quary(qtype int, v ...interface{}) (err error) {
	if qtype != QtClearAllRecord && (v == nil || len(v) < 1) {
		return fmt.Errorf("[MainAccountOperation:Quary] Quary parameter 'v' is not be null")
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
	case QtGetCount:
		ths.locker.Lock()
		c := v[0].(*int64)
		err = ths.getCount(c)
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

func (ths *MainAccountOperation) getCount(c *int64) error {
	temp := new(TMainAccount)
	cnt, err := ths.engine.Count(temp)
	*c = cnt
	return err
}

func (ths *MainAccountOperation) CmdExec(c string, v ...interface{}) error {
	ths.locker.Lock()
	defer ths.locker.Unlock()
	_, err := ths.engine.Exec(c, v...)
	if err != nil {
		err = fmt.Errorf("[MainAccountOperation:CmdExec] %v", err)
	}
	return err
}

func (ths *MainAccountOperation) CmdQuery(c string, v ...interface{}) (ret []map[string][]byte, err error) {
	ths.locker.Lock()
	defer ths.locker.Unlock()
	return ths.engine.Query(c, v...)
}

func (ths *MainAccountOperation) GetEngine() *xorm.Engine {
	return ths.engine
}
