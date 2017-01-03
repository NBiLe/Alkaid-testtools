package menu

import (
	"fmt"
	"sync"
	"time"

	_db "github.com/fuyaocn/evaluatetools/db"
	_L "github.com/fuyaocn/evaluatetools/log"
	_kp "github.com/stellar/go/keypair"
)

const (
	NumberOfLocalAccount = iota
	InputNumberError
)

// CreateAccount 创建账户
type CreateAccount struct {
	MenuSubItem
}

// InitMenu 初始化
func (ths *CreateAccount) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.MenuSubItem.InitMenu(parent, key)
	ths.Exec = ths.execute

	ths.title = []string{
		LEnglish: "Create local stellar account",
	}
	ths.infoStrings = []map[int]string{
		LEnglish: map[int]string{
			NumberOfLocalAccount: " > Please enter the number of local accounts that need to be created (If you enter a number less than or equal to 0, it will end this operation) :",
		},
	}
	return ths
}

func (ths *CreateAccount) execute(isSync bool) {
	_L.LoggerInstance.Info(" ** Create local stellar account ** \r\n")
	var input string
	fmt.Printf("\r\n" + ths.infoStrings[ths.languageIndex][NumberOfLocalAccount])
	fmt.Scanf("%s\n", &input)

	num, bNum := IsNumber(input)

	if bNum {
		if num > 0 {
			ths.run(num)
		}
	}
	if !isSync {
		ths.ASyncChan <- 0
	}
}

// BaseCount 基准数
var BaseCount int = 100

// run execute process
func (ths *CreateAccount) run(count int) {
	_L.LoggerInstance.InfoPrint("Create stellar account into local database BEGIN ...\r\n")
	_L.LoggerInstance.InfoPrint("\tcreate count = %d\r\n", count)
	_L.LoggerInstance.InfoPrint("\tcreate start time = %s\r\n", time.Now().Format("2006-01-02 15:04:05.000"))

	times := count/BaseCount + 1
	leftTimes := count % BaseCount
	wt := &sync.WaitGroup{}
	for {
		if times == 0 {
			break
		}
		times--
		accSize := BaseCount
		if times == 0 {
			accSize = leftTimes
		}

		if accSize == 0 {
			continue
		}

		wt.Add(1)
		go ths.add2DB(accSize, wt)

		wt.Wait()
	}
	_L.LoggerInstance.InfoPrint("\tcreate end time = %s\r\n", time.Now().Format("2006-01-02 15:04:05.000"))
	_L.LoggerInstance.InfoPrint("Create stellar account into local database Complete! \r\n")
}

func (ths *CreateAccount) add2DB(sz int, w *sync.WaitGroup) {
	defer w.Done()
	accs := make([]*_db.TAccount, sz)
	for idx := 0; idx < sz; idx++ {
		accs[idx] = ths.newAcc()
		if accs[idx] == nil {
			panic(1)
		}
	}

	err := _db.DataBaseInstance.AddGroupStellarAccount(accs)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("[CreateAccount:Execute] Add keypair to database has error \r\n%+v\r\n", err)
	}
}

func (ths *CreateAccount) newAcc() *_db.TAccount {
	full, err := _kp.Random()
	if err != nil {
		_L.LoggerInstance.ErrorPrint("[CreateAccount:Execute] Create keypair has error \r\n%+v\r\n", err)
		return nil
	}
	t := time.Now()
	return &_db.TAccount{
		AccountID:      full.Address(),
		SecertAddr:     full.Seed(),
		CreateTime:     t,
		CreateTimeUnix: t.UnixNano(),
		LastUpdateTime: t,
		UpdateTimeUnix: t.UnixNano(),
		Active:         "F",
	}
}
