package menu

import (
	"fmt"

	_L "github.com/fuyaocn/evaluatetools/log"
)

// LocalAccount 本地账户菜单
type LocalAccount struct {
	MenuSubItem
}

// InitMenu 初始化菜单
func (ths *LocalAccount) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.MenuSubItem.InitMenu(parent, key)
	ths.Exec = ths.execute

	ths.title = []string{
		LEnglish: "Local Stellar Account",
	}

	ths.infoStrings = []map[int]string{
		LEnglish: map[int]string{},
	}

	count := &CountAccount{}
	ths.AddSubItem(count.InitMenu(ths, "1"))

	ca := &CreateAccount{}
	ths.AddSubItem(ca.InitMenu(ths, "1"))

	aa := &ActiveAccount{}
	ths.AddSubItem(aa.InitMenu(ths, "1"))

	cha := &CheckAccounts{}
	ths.AddSubItem(cha.InitMenu(ths, "1"))

	clear := &ClearAccounts{}
	ths.AddSubItem(clear.InitMenu(ths, "1"))

	rp := &ReturnParentMenu{}
	ths.AddSubItem(rp.InitMenu(ths, "1"))

	ea := &ExitApp{}
	ths.AddSubItem(ea.InitMenu(ths, "1"))

	return ths
}

func (ths *LocalAccount) execute(isSync bool) {
	_L.LoggerInstance.Info(" ** Local account menu ** \r\n")
	for {
		fmt.Printf("\n\n%s\r\n\n", ths.GetTitlePath(ths.languageIndex))
		ths.PrintSubmenu()
		fmt.Printf("\n %s", ths.GetInputMemo(ths.languageIndex))

		var input string

		_, err := fmt.Scanf("%s\n", &input)
		if err == nil {
			selectIndex, b := IsNumber(input)
			if b {
				if selectIndex <= len(ths.subItems) && selectIndex >= 0 {
					ths.subItems[selectIndex-1].ExecuteFunc(false)
					ret := ths.subItems[selectIndex-1].ExecFlag()
					if ret == BACK_TO_MENU_FLAG {
						break
					}
				}
			}
		} else {
			fmt.Println(err)
		}
	}

	if !isSync {
		ths.ASyncChan <- 1
	}

}
