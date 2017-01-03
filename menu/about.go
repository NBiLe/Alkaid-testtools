package menu

import (
	"fmt"

	_L "github.com/fuyaocn/evaluatetools/log"
)

const (
	SA_INFO_MEMO = iota
)

// SoftwareAbout 关于
type SoftwareAbout struct {
	MenuSubItem
}

// InitMenu 初始化
func (ths *SoftwareAbout) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.MenuSubItem.InitMenu(parent, key)
	ths.Exec = ths.execute

	ths.title = []string{
		LEnglish: "About",
	}

	ths.infoStrings = []map[int]string{
		LEnglish: map[int]string{
			SA_INFO_MEMO: "\tSoftware Version : 1.0.0.0\r\n",
		},
	}
	return ths
}

func (ths *SoftwareAbout) execute(isSync bool) {
	_L.LoggerInstance.Info(" ** About show ** \r\n")
	fmt.Println("")
	fmt.Println(ths.infoStrings[ths.languageIndex][SA_INFO_MEMO])
	fmt.Println("")

	var input string

	fmt.Scanf("%s\n", &input)

	if !isSync {
		ths.ASyncChan <- 0
	}
}
