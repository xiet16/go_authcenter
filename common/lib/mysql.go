package lib

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

var TimeFormat = "2006-01-02 15:04:05"
//初始化线程池
func InitDBPool(path string) error {
   DbConfMap := &MysqlMapConf{}
   err := ParseConfig(path,DbConfMap)
	if err!=nil {
		fmt.Printf("[INFO] %s%s\n", time.Now().Format(TimeFormat), " empty mysql config.")
	}

	DBMapPool = map[string]*sql.DB{}
    GORMMapPool = map[string]*gorm.DB{}

	for configName,dbConf := range DbConfMap.List{
		dbpool,err :=sql.Open("mysql",dbConf.DataSourceName) // create db pool
		if err!=nil {
			return err
		}

		dbpool.SetMaxOpenConns(dbConf.MaxOpenConn)
		dbpool.SetMaxIdleConns(dbConf.MaxIdleConn)
		dbpool.SetConnMaxLifetime(time.Duration(dbConf.MaxConnLifeTime)*time.Second)
		err=dbpool.Ping()
		if err!=nil {
			return err
		}

		//gorm
		dbgrom,err := gorm.Open("mysql",dbConf.DataSourceName)
		//dbgrom.SingularTable(true) //表明后面s设置
		if err!=nil {
			return err
		}
		dbgrom.SingularTable(true)
		err=dbgrom.DB().Ping()
		if err != nil {
            return err
		}
		dbgrom.LogMode(true)
		//dbgorm.LogCtx(true)
		//dbgrom.SetLogger(&mysqlg)
		dbgrom.DB().SetMaxIdleConns(dbConf.MaxIdleConn)
		dbgrom.DB().SetMaxOpenConns(dbConf.MaxOpenConn)
		dbgrom.DB().SetConnMaxLifetime(time.Duration(dbConf.MaxConnLifeTime)*time.Second)
		dbgrom.Callback().Create().Replace("gorm:update_time_stamp",updateTimeStampForCreateCallback)
		dbgrom.Callback().Update().Register("gorm:update_time_stamp",updateTimeStampForUpdateCallback)
		DBMapPool[configName] = dbpool
		GORMMapPool[configName] = dbgrom
	}

	if dbpool,err := GetDBPool("default");err!=nil {
		DBDefaultPool = dbpool
	}

	if dbpool,err := GetGormPool("default");err!=nil {
       GORMDefaultPool = dbpool
	}

	return nil
}

func GetDBPool(name string) (*sql.DB,error) {
	if dbpool,ok := DBMapPool[name];ok {
		return dbpool,nil
	}

	return nil,errors.New("get pool error")
}

func GetGormPool(name string) (*gorm.DB,error) {
	if dbpool,ok := GORMMapPool[name];ok {
		return dbpool,nil
	}

	return nil,errors.New("get pool error")
}

func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreateTime"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("UpdateTime"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}
// 注册更新钩子在持久化之前
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("UpdateTime", time.Now().Unix())
	}
}