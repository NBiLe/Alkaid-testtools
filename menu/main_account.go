package menu

import (
	"fmt"

	_L "github.com/fuyaocn/evaluatetools/log"
)

// MainAccount 本地主账户菜单
type MainAccount struct {
	MenuSubItem
}

// InitMenu 初始化菜单
func (ths *MainAccount) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.MenuSubItem.InitMenu(parent, key)
	ths.Exec = ths.execute

	ths.title = []string{
		LEnglish: "Main Stellar Account",
	}

	ths.infoStrings = []map[int]string{
		LEnglish: map[int]string{},
	}

	ma := &CreateMainAccount{}
	ths.AddSubItem(ma.InitMenu(ths, "1"))

	rp := &ReturnParentMenu{}
	ths.AddSubItem(rp.InitMenu(ths, "1"))

	ea := &ExitApp{}
	ths.AddSubItem(ea.InitMenu(ths, "1"))

	return ths
}

func (ths *MainAccount) execute(isSync bool) {
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
