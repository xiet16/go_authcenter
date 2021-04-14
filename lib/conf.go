package lib

import (
	"bytes"
	"database/sql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
)

type MysqlMapConf struct {
	List map[string]*MySQLConf `mapstructure:"list"`
}

type MySQLConf struct {
	DriverName      string `mapstructure:"driver_name"`
	DataSourceName  string `mapstructure:"data_source_name"`
	MaxOpenConn     int    `mapstructure:"max_open_conn"`
	MaxIdleConn     int    `mapstructure:"max_idle_conn"`
	MaxConnLifeTime int    `mapstructure:"max_conn_life_time"`
}

type RedisMapConf struct {
	List map[string]*RedisConf `mapstructure:"list"`
}

type RedisConf struct {
	ProxyList    []string `mapstructure:"proxy_list"`
	Password     string   `mapstructure:"password"`
	Db           int      `mapstructure:"db"`
	ConnTimeout  int      `mapstructure:"conn_timeout"`
	ReadTimeout  int      `mapstructure:"read_timeout"`
	WriteTimeout int      `mapstructure:"write_timeout"`
}

type ConfConnClientMap struct {
	List map[string]*ConnClientConf `mapstructure:"list`
}

type ConnClientConf struct {
	ID string `mapstructure:"id"`
	Secret string `mapstructure:"secret"`
	Name string `mapstructure:"name"`
	Domain string `mapstructure:"domain"`
	Scope []Scope `mapstructure:"scope"`
}

type Scope struct {
	ID string `mapstructure:"id"`
	Title string `mapstructure:"title"`
}

//

var DBMapPool map[string] *sql.DB
var GORMMapPool map[string]*gorm.DB
var DBDefaultPool *sql.DB
var GORMDefaultPool *gorm.DB
var ViperConfMap map[string]*viper.Viper

var ConfRedisMap *RedisMapConf
var ConfConnCientMap *ConfConnClientMap

//初始化配置文件
func InitViperConf() error {
	f, err := os.Open(ConfEnvPath + "/")
	if err != nil {
		return err
	}
	fileList, err := f.Readdir(1024)
	if err != nil {
		return err
	}
	for _, f0 := range fileList {
		if !f0.IsDir() {
			bts, err := ioutil.ReadFile(ConfEnvPath + "/" + f0.Name())
			if err != nil {
				return err
			}
			v := viper.New()
			v.SetConfigType("toml")
			v.ReadConfig(bytes.NewBuffer(bts))
			pathArr := strings.Split(f0.Name(), ".")
			if ViperConfMap == nil {
				ViperConfMap = make(map[string]*viper.Viper)
			}
			ViperConfMap[pathArr[0]] = v
		}
	}
	return nil
}

func InitRedisConf(path string) error {
    redisConf := &RedisMapConf{}
    err := ParseConfig(path,redisConf)
	if err!=nil {
		return err
	}

	ConfRedisMap = redisConf
	return nil
}

func InitConnClient(path string) error {
    clientConf :=&ConfConnClientMap{}
    err := ParseConfig(path,clientConf)
	if err !=nil {
		return err
	}
	ConfConnCientMap = clientConf
	return nil
}
