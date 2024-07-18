package layer

import (
	"fmt"
	"gorm.io/gorm"
)

// 默认db
var DefaultDBClient *gorm.DB

// 可选的db集合
var NamedDBClient map[string]*gorm.DB

// 是否关闭自动读主库，默认开启
var closeAutoReadMaster bool

type IDao interface {
	ILayer
	DaoFunc()
	GetDB() *gorm.DB
	GetDBByName(name string) *gorm.DB
	SetDB(db *gorm.DB)
	ResetDB()
	ClearDB()
	SetTable(tableName string)
	GetTable() string
}

type Dao struct {
	Layer
	db         *gorm.DB
	defaultDB  *gorm.DB
	tableName  string
	partionNum int
}

func (entity *Dao) DaoFunc() {
	fmt.Print("this is dao func\n")
}

func (entity *Dao) OnCreate(param IFlowParam) {
	entity.Layer.OnCreate(param)
}

func (entity *Dao) GetDB() *gorm.DB {
	var db *gorm.DB
	if entity.db != nil {
		db = entity.db
	} else if entity.defaultDB != nil {
		db = entity.defaultDB.WithContext(entity.GetContext())
	} else if DefaultDBClient != nil {
		db = DefaultDBClient.WithContext(entity.GetContext())
	}
	if db != nil {
		db = db.Table(entity.GetTable())
	}
	return db
}

func (entity *Dao) GetDBByName(name string) *gorm.DB {
	var db *gorm.DB
	if entity.db != nil {
		db = entity.db
	} else {
		// 没有name，取默认的db
		if name == "" && DefaultDBClient != nil {
			db = DefaultDBClient.WithContext(entity.GetContext())
		} else if name != "" && NamedDBClient != nil {
			// 有name，尝试找对应的db
			if dbClient, exist := NamedDBClient[name]; exist {
				db = dbClient.WithContext(entity.GetContext())
			}
		}
	}
	if db != nil {
		db = db.Table(entity.GetTable())
	}
	return db
}

func (entity *Dao) SetDB(db *gorm.DB) {
	entity.db = db
}

func (entity *Dao) SetDefaultDB(db *gorm.DB) {
	entity.defaultDB = db
}

func (entity *Dao) ResetDB() {
	// 优先使用entity的defaultDB
	if entity.defaultDB != nil {
		entity.db = entity.defaultDB.WithContext(entity.GetContext())
	} else {
		entity.db = DefaultDBClient.WithContext(entity.GetContext())
	}
}

func (entity *Dao) ClearDB() {
	entity.db = nil
}

func (entity *Dao) SetTable(tableName string) {
	entity.tableName = tableName
}

func (entity *Dao) GetTable() string {
	return entity.tableName
}

func (entity *Dao) SetPartitionNum(num int) {
	entity.partionNum = num
}

func (entity *Dao) GetPartitionNum() int {
	return entity.partionNum
}

func (entity *Dao) GetPartitionTable(value int64) string {
	return fmt.Sprintf("%s%d", entity.GetTable(), value%int64(entity.partionNum))
}

func SetDefaultDBClient(db *gorm.DB) {
	DefaultDBClient = db
}

func SetNamedDBClient(namedDbs map[string]*gorm.DB) {
	NamedDBClient = namedDbs
}

func SetCloseAutoReadMaster(close bool) {
	closeAutoReadMaster = close
}
