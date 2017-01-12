package menu

import (
	"fmt"

	_L "github.com/fuyaocn/evaluatetools/log"
)

// MenuInfo 菜单
type MenuInfo struct {
	MenuSubItem
	currentLevel  int
	WelcomeString []string
}

// MainMenuInstace 菜单唯一实例
var MainMenuInstace *MenuInfo

func init() {
	MainMenuInstace = new(MenuInfo)
	MainMenuInstace.MenuSubItem.InitMenu(nil, "0")
	MainMenuInstace.currentLevel = 0
	MainMenuInstace.title = []string{
		LEnglish: "Menu",
	}
	MainMenuInstace.WelcomeString = []string{
		LEnglish: " ##      please choose the function you need       ##\n",
	}
	MainMenuInstace.Exec = MainMenuInstace.execute

	ma := &MainAccount{}
	MainMenuInstace.AddSubItem(ma.InitMenu(MainMenuInstace, "0"))

	la := &LocalAccount{}
	MainMenuInstace.AddSubItem(la.InitMenu(MainMenuInstace, "0"))

	// accInfo := &AccountInfo{}
	// accInfo.InitAccInfo(MainMenuInstace, "0")
	// MainMenuInstace.AddSubItem(accInfo)

	// mergeAcc := &MergeAccount{}
	// mergeAcc.InitMerge(MainMenuInstace, "0")
	// MainMenuInstace.AddSubItem(mergeAcc)

	sa := &SoftwareAbout{}
	MainMenuInstace.AddSubItem(sa.InitMenu(MainMenuInstace, "0"))

	ea := &ExitApp{}
	MainMenuInstace.AddSubItem(ea.InitMenu(MainMenuInstace, "0"))
}

// Execute 执行函数
func (ths *MenuInfo) execute(isSync bool) {
	for {
		_L.LoggerInstance.Info(" ** Step in main menu ** \r\n")
		fmt.Println("\r\n******************************************************")
		fmt.Println(ths.getWelcomeString(ths.languageIndex))
		fmt.Println(" " + ths.GetTitle(ths.languageIndex))
		ths.PrintSubmenu()
		fmt.Printf("\n %s", ths.GetInputMemo(ths.languageIndex))

		var input string

		_, err := fmt.Scanf("%s\n", &input)
		if err == nil {
			selectIndex, b := IsNumber(input)
			if b {
				if selectIndex <= len(ths.subItems) && selectIndex > 0 {
					ths.subItems[selectIndex-1].ExecuteFunc(false)
					ths.subItems[selectIndex-1].ExecFlag()
				}
			}
		}
	}
}

func (ths *MenuInfo) getWelcomeString(langType int) string {
	return ths.WelcomeString[langType]
}
