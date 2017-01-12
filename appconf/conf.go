package appconf

import (
	"strconv"
	"strings"

	_c "github.com/astaxie/beego/config"
	_db "github.com/fuyaocn/evaluatetools/db"
	_L "github.com/fuyaocn/evaluatetools/log"
)

// ConfigInstance 配置文件唯一实例
var ConfigInstance *ConfigController

// ConfigController 配置文件读取控制器
type ConfigController struct {
	appConfig   _c.Configer
	skeyConfig  _c.Configer
	servConfig  _c.Configer
	activConfig _c.Configer
}

// NewConfigController new ConfigController
func NewConfigController() *ConfigController {
	ret := new(ConfigController)
	ret.Init()
	return ret
}

// Init 初始化
func (ths *ConfigController) Init() {
	var err error
	ths.appConfig, err = _c.NewConfig("json", "config/app.conf")
	if err != nil {
		_L.LoggerInstance.ErrorPrint("Read 'app.conf' has error : \r\n%v\r\n", err)
		return
	}

	ths.skeyConfig, err = _c.NewConfig("json", "config/skey.conf")
	if err != nil {
		_L.LoggerInstance.ErrorPrint("Read 'skey.conf' has error : \r\n%v\r\n", err)
		return
	}

	ths.servConfig, err = _c.NewConfig("json", "config/server.conf")
	if err != nil {
		_L.LoggerInstance.ErrorPrint("Read 'server.conf' has error : \r\n%v\r\n", err)
		return
	}

	ths.activConfig, err = _c.NewConfig("ini", "config/active.conf")
	if err != nil {
		_L.LoggerInstance.ErrorPrint("Read 'active.conf' has error : \r\n%v\r\n", err)
	}
}

// GetDateBaseConf 获取数据库相关配置
func (ths *ConfigController) GetDateBaseConf() (ret *_db.DatabaseInfo) {

	ret = new(_db.DatabaseInfo)
	ret.DbType = _db.DatabaseType(ths.appConfig.String("dbtype"))
	ret.AliasName = ths.appConfig.String("aliasname")
	ret.Host = ths.appConfig.String("dbHost")
	ret.Port = ths.appConfig.String("dbPort")
	ret.UserName = ths.appConfig.String("dbUserName")
	ret.Password = ths.appConfig.String("dbPassword")
	ret.IsDebug = ths.appConfig.String("debug") == "true"
	return
}

// GetRootAccountID 获取母ID配置
func (ths *ConfigController) GetRootAccountID() string {
	return ths.skeyConfig.String("RootAccountID")
}

// GetRootSecretKey 获取母key配置
func (ths *ConfigController) GetRootSecretKey() string {
	return ths.skeyConfig.String("RootSecretKey")
}

// GetMainIssuerID get IssuerID
func (ths *ConfigController) GetMainIssuerID() string {
	return ths.skeyConfig.String("IssuerID")
}

// GetMainCredit 获取Credit
func (ths *ConfigController) GetMainCredit() string {
	return ths.skeyConfig.String("Credit")
}

// GetHorizonTest get test horizon url
func (ths *ConfigController) GetHorizonTest() string {
	return ths.servConfig.String("horizon_test")
}

// GetHorizonLive get live horizon url
func (ths *ConfigController) GetHorizonLive() string {
	return ths.servConfig.String("horizon_live")
}

// GetNetwork 得到当前配置使用网络
func (ths *ConfigController) GetNetwork() string {
	return ths.servConfig.String("current_network")
}

// GetHorizon 根据网络得到Horizon
func (ths *ConfigController) GetHorizon() string {
	if strings.ToLower(ths.GetNetwork()) == "live" {
		return ths.GetHorizonLive()
	}
	return ths.GetHorizonTest()
}

// GetCore get core url
func (ths *ConfigController) GetCore() string {
	return ths.servConfig.String("stellar-core")
}

// GetStartBalance get active start balance
func (ths *ConfigController) GetStartBalance() string {
	return ths.activConfig.String("startingBalance")
}

// GetMainStartBalance get main account active start balance
func (ths *ConfigController) GetMainStartBalance() string {
	return ths.activConfig.String("mainStartingBalance")
}

// GetActiveDepth get active depth
func (ths *ConfigController) GetActiveDepth() int {
	ret, _ := strconv.Atoi(ths.activConfig.String("Depth"))
	return ret
}

// GetActiveLevel get active level
func (ths *ConfigController) GetActiveLevel() int {
	ret, _ := strconv.Atoi(ths.activConfig.String("Level"))
	return ret
}

// GetHorizonHeader get horizon http header config
func (ths *ConfigController) GetHorizonHeader() map[string]string {
	headers, err := ths.activConfig.GetSection("horizon_header")
	if err != nil {
		panic(err)
	}
	return headers
}
