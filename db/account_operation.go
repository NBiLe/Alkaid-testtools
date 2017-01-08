package db

import (
	"fmt"
	"sync"

	_L "github.com/fuyaocn/evaluatetools/log"
	"github.com/go-xorm/xorm"
)

// OperationAccount operation for account
type OperationAccount struct {
	engine *xorm.Engine
	locker *sync.Mutex
}

// Init 初始化
func (ths *OperationAccount) Init(e *xorm.Engine) {
	ths.locker = &sync.Mutex{}
	ths.engine = e
}

// GetKey get key string
func (ths *OperationAccount) GetKey() string {
	return KeyAccount
}

// Quary quary exeute
func (ths *OperationAccount) Quary(qtype int, v ...interface{}) (err error) {
	if qtype != QtClearAllRecord && (v == nil || len(v) < 1) {
		return fmt.Errorf("[OperationAccount:Quary] Quary parameter 'v' is not be null")
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
		err = ths.clearStellarAccount()
		ths.locker.Unlock()
	case QtGetCount:
		ths.locker.Lock()
		act := v[0].(string)
		all := v[1].(bool)
		c := v[2].(*int64)
		err = ths.getCount(act, all, c)
		ths.locker.Unlock()
	default:
		_L.LoggerInstance.ErrorPrint("Unknown quary type [%d]\r\n", qtype)
	}

	return
}

func (ths *OperationAccount) clearStellarAccount() (err error) {
	sqlcmd := "truncate table t_account;"
	_, err = ths.engine.Query(sqlcmd)
	return err
}

func (ths *OperationAccount) getCount(active string, all bool, c *int64) error {
	temp := new(TAccount)
	if !all {
		temp.Active = active
	}
	cnt, err := ths.engine.Count(temp)
	*c = cnt
	return err
}

func (ths *OperationAccount) getAccountIter(active string, callback xorm.IterFunc) error {
	temp := &TAccount{
		Active: active,
	}
	return ths.engine.Iterate(temp, callback)
}

// GetCountRecords get count of record
func (ths *OperationAccount) GetCountRecords(active string, cnt, offset int64) (ret []map[string][]byte, err error) {
	ths.locker.Lock()
	defer ths.locker.Unlock()
	return ths.getCountRecords(active, cnt, offset)
}

func (ths *OperationAccount) getCountRecords(active string, cnt, offset int64) (ret []map[string][]byte, err error) {
	if active == "" {
		sql := "SELECT * FROM `t_account` order by id limit ? offset ?"
		return ths.engine.Query(sql, cnt, offset)
	}
	sql := "SELECT * FROM `t_account` WHERE `active`=? order by id limit ? offset ?"
	return ths.engine.Query(sql, active, cnt, offset)
}

func (ths *OperationAccount) CmdExec(c string, v ...interface{}) error {
	ths.locker.Lock()
	defer ths.locker.Unlock()
	_, err := ths.engine.Exec(c, v...)
	if err != nil {
		err = fmt.Errorf("[OperationAccount:CmdExec] %v", err)
	}
	return err
}
