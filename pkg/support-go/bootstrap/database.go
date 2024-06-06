package bootstrap

import (
	"errors"
	"fmt"
	"github.com/JasonMetal/simple-develop-template/pkg/support-go/helper/config"
	yCfg "github.com/olebedev/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

type MysqlInstance struct {
	DSN string
	DB  *gorm.DB
}

var mysqlDbList = make(map[string]MysqlInstance)

func InitMysql() {
	dbList := getDbNames("mysql")

	for _, dbname := range dbList {
		instance, err := initDbConn(dbname)
		if err == nil {
			mysqlDbList[dbname] = instance
		}
	}
}

func mysqlLogger() logger.Interface {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,   // 慢 SQL 阈值
			LogLevel:                  logger.Silent, // 日志级别
			IgnoreRecordNotFoundError: true,          // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,         // 禁用彩色打印
		},
	)

	return newLogger
}

func initDbConn(dbName string) (MysqlInstance, error) {
	path := fmt.Sprintf("%smanifest/config/%s/mysql.yml", ProjectPath(), DevEnv)

	cfg, err := config.GetConfig(path)

	maxOpenConns, _ := cfg.Int("mysql." + dbName + ".maxOpenConns")
	maxIdleConns, _ := cfg.Int("mysql." + dbName + ".maxIdleConns")
	maxLifetime, _ := cfg.Int("mysql." + dbName + ".maxLifetime")
	tablePrefix, _ := cfg.String("mysql." + dbName + ".tablePrefix")
	debug, _ := cfg.Bool("mysql." + dbName + ".debug")
	charset, _ := cfg.String("mysql." + dbName + ".charset")

	servers, err := cfg.List("mysql." + dbName + ".servers")
	if err != nil || len(servers) < 1 {
		return MysqlInstance{}, err
	}

	host, _ := yCfg.Get(servers[0], "host")
	port, _ := yCfg.Get(servers[0], "port")
	name, _ := yCfg.Get(servers[0], "db")
	user, _ := yCfg.Get(servers[0], "user")
	passwd, _ := yCfg.Get(servers[0], "passwd")

	addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		//addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		user.(string), passwd.(string), host.(string), port.(int), name.(string), charset)

	if connTimeout, err := yCfg.Get(servers[0], "connTimeout"); err == nil {
		addr += fmt.Sprintf("&timeout=%dms", connTimeout.(int))
	}
	if readTimeout, err := yCfg.Get(servers[0], "readTimeout"); err == nil {
		addr += fmt.Sprintf("&readTimeout=%dms", readTimeout.(int))
	}
	if writeTimeout, err := yCfg.Get(servers[0], "writeTimeout"); err == nil {
		addr += fmt.Sprintf("&writeTimeout=%dms", writeTimeout.(int))
	}
	gormConfig := &gorm.Config{
		Logger: mysqlLogger(),
	}

	if tablePrefix != "" {
		gormConfig.NamingStrategy = schema.NamingStrategy{
			TablePrefix: tablePrefix,
		}
	}

	db, err := gorm.Open(mysql.Open(addr), &gorm.Config{
		Logger: mysqlLogger(),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix,
			SingularTable: true,
		},
	})

	if err != nil {
		//return MysqlInstance{}, errors.New("connection is not exist")
		return MysqlInstance{}, err
	}

	if debug {
		db = db.Debug()
	}

	DB, err := db.DB()

	err = DB.Ping()
	if err != nil {
		return MysqlInstance{}, err
	}

	if maxLifetime > 0 {
		DB.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)
	}

	DB.SetMaxIdleConns(maxIdleConns)
	DB.SetMaxOpenConns(maxOpenConns)

	return MysqlInstance{fmt.Sprintf("%s:%d/%s", host.(string), port.(int), name.(string)), db}, nil
}

func GetMysqlInstance(dbName string) (MysqlInstance, error) {
	if instance, ok := mysqlDbList[dbName]; ok {
		return instance, nil
	} else {
		return instance, errors.New(dbName + " db is null")
	}

}
