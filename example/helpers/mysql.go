package helpers

import (
	"github.com/tiant-developer/golib/example/conf"
	"github.com/tiant-developer/golib/orm"
	"gorm.io/gorm"
)

var (
	MysqlClient *gorm.DB
)

func InitMysql() {
	var err error
	for name, dbConf := range conf.GlobalC.Mysql {
		switch name {
		case "default":
			MysqlClient, err = orm.InitMysqlClient(dbConf)
		}
		if err != nil {
			panic("mysql connect error: %v" + err.Error())
		}
	}
}

func CloseMysql() {
}
