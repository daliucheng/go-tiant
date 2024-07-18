package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/tiant-developer/go-tiant/utils"
	"github.com/tiant-developer/go-tiant/zlog"
	ormUtil "gorm.io/gorm/utils"
	"gorm.io/plugin/dbresolver"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	driver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const prefix = "@@mysql."

type MysqlConf struct {
	Service         string        `yaml:"service"`
	DataBase        string        `yaml:"database"`
	Addr            string        `yaml:"addr"`
	SlaveAddrs      []string      `yaml:"slaveAddrs"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Charset         string        `yaml:"charset"`
	MaxIdleConns    int           `yaml:"maxidleconns"`
	MaxOpenConns    int           `yaml:"maxopenconns"`
	ConnMaxIdlTime  time.Duration `yaml:"maxIdleTime"`
	ConnMaxLifeTime time.Duration `yaml:"connMaxLifeTime"`
	ConnTimeOut     time.Duration `yaml:"connTimeOut"`
	WriteTimeOut    time.Duration `yaml:"writeTimeOut"`
	ReadTimeOut     time.Duration `yaml:"readTimeOut"`

	// sql 字段最大长度
	MaxSqlLen int `yaml:"maxSqlLen"`
}

func (conf *MysqlConf) checkConf() {

	if conf.MaxIdleConns == 0 {
		conf.MaxIdleConns = 50
	}
	if conf.MaxOpenConns == 0 {
		conf.MaxOpenConns = 50
	}
	if conf.ConnMaxIdlTime == 0 {
		conf.ConnMaxIdlTime = 5 * time.Minute
	}
	if conf.ConnMaxLifeTime == 0 {
		conf.ConnMaxLifeTime = 10 * time.Minute
	}
	if conf.ConnTimeOut == 0 {
		conf.ConnTimeOut = 3 * time.Second
	}
	if conf.WriteTimeOut == 0 {
		conf.WriteTimeOut = 1200 * time.Millisecond
	}
	if conf.ReadTimeOut == 0 {
		conf.ReadTimeOut = 1200 * time.Millisecond
	}
	if conf.MaxSqlLen == 0 {
		// 日志中sql字段长度：
		// 如果不指定使用默认2048；如果<0表示不展示sql语句；否则使用用户指定的长度，过长会被截断
		conf.MaxSqlLen = 2048
	}
}

func InitMysqlClient(conf MysqlConf) (client *gorm.DB, err error) {
	conf.checkConf()

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=%s&readTimeout=%s&writeTimeout=%s&parseTime=True&loc=Asia%%2FShanghai",
		conf.User,
		conf.Password,
		conf.Addr,
		conf.DataBase,
		conf.ConnTimeOut,
		conf.ReadTimeOut,
		conf.WriteTimeOut)
	dsnArr := []string{}
	for _, s := range conf.SlaveAddrs {
		dsn2 := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=%s&readTimeout=%s&writeTimeout=%s&parseTime=True&loc=Asia%%2FShanghai",
			conf.User,
			conf.Password,
			s,
			conf.DataBase,
			conf.ConnTimeOut,
			conf.ReadTimeOut,
			conf.WriteTimeOut)
		dsnArr = append(dsnArr, dsn2)
	}
	dsnArr = append(dsnArr, dsn)
	if conf.Charset != "" {
		dsn = dsn + "&charset=" + conf.Charset
	}

	l := newLogger(&conf)
	_ = driver.SetLogger(l)

	c := &gorm.Config{
		SkipDefaultTransaction:                   true,
		NamingStrategy:                           nil,
		FullSaveAssociations:                     false,
		Logger:                                   l,
		NowFunc:                                  nil,
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		AllowGlobalUpdate:                        false,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  nil,
	}

	client, err = gorm.Open(mysql.Open(dsn), c)
	if err != nil {
		return client, err
	}
	var p []gorm.Dialector
	for _, s := range dsnArr {
		p = append(p, mysql.Open(s))
	}

	client.Use(dbresolver.Register(dbresolver.Config{
		// `db2` 作为 sources，`db3`、`db4` 作为 replicas
		Sources:  []gorm.Dialector{mysql.Open(dsn)},
		Replicas: p,
		// sources/replicas 负载均衡策略
		Policy: dbresolver.RandomPolicy{},
	}))

	sqlDB, err := client.DB()
	if err != nil {
		return client, err
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)

	// SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)

	// SetConnMaxLifetime 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(conf.ConnMaxLifeTime)

	// only for go version >= 1.15 设置最大空闲连接时间
	sqlDB.SetConnMaxIdleTime(conf.ConnMaxIdlTime)

	return client, nil
}

type ormLogger struct {
	Service   string
	Database  string
	MaxSqlLen int
	logger    *zlog.Logger
}

func newLogger(conf *MysqlConf) *ormLogger {
	s := conf.Service
	if conf.Service == "" {
		s = conf.DataBase
	}

	return &ormLogger{
		Service:   s,
		Database:  conf.DataBase,
		MaxSqlLen: conf.MaxSqlLen,
		logger:    zlog.ZapLogger.WithOptions(zlog.AddCallerSkip(2)),
	}
}

// go-sql-driver error log
func (l *ormLogger) Print(args ...interface{}) {
	l.logger.Error(fmt.Sprint(args...), l.commonFields(nil)...)
}

// LogMode log mode
func (l *ormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

// Info print info
func (l ormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	m := fmt.Sprintf(msg, append([]interface{}{ormUtil.FileWithLineNum()}, data...)...)
	// 非trace日志改为debug级别输出
	l.logger.Debug(m, l.commonFields(ctx)...)
}

// Warn print warn messages
func (l ormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	m := fmt.Sprintf(msg, append([]interface{}{ormUtil.FileWithLineNum()}, data...)...)
	l.logger.Warn(m, l.commonFields(ctx)...)
}

// Error print error messages
func (l ormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	m := fmt.Sprintf(msg, append([]interface{}{ormUtil.FileWithLineNum()}, data...)...)
	l.logger.Error(m, l.commonFields(ctx)...)
}

// Trace print sql message
func (l ormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	end := time.Now()
	elapsed := end.Sub(begin)
	cost := float64(elapsed.Nanoseconds()/1e4) / 100.0

	// 请求是否成功
	msg := "mysql do success"
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 没有找到记录不统计在请求错误中
		msg = err.Error()
	}

	sql, rows := fc()
	if l.MaxSqlLen < 0 {
		sql = ""
	} else if len(sql) > l.MaxSqlLen {
		sql = sql[:l.MaxSqlLen]
	}

	fields := l.commonFields(ctx)
	fields = append(fields,
		zlog.Int64("affectedrow", rows),
		zlog.String("reqEndTime", utils.GetFormatRequestTime(end)),
		zlog.String("reqStartTime", utils.GetFormatRequestTime(begin)),
		zlog.Float64("cost", cost),
		zlog.String("sql", sql),
	)

	l.logger.Info(msg, fields...)
}

func (l ormLogger) commonFields(ctx context.Context) []zlog.Field {
	var requestID string
	if c, ok := ctx.(*gin.Context); (ok && c != nil) || (!ok && !IsNil(ctx)) {
		requestID, _ = ctx.Value(zlog.ContextKeyRequestID).(string)
	}

	fields := []zlog.Field{
		zlog.String("requestId", requestID),
	}
	return fields
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
