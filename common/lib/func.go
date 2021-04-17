package lib

import (
	"fmt"
	"strings"
	"time"
)

var ConfEnvPath string //配置文件夹
var ConfEnv string     //配置环境名 比如：dev prod test

// 解析配置文件目录
//
// 配置文件必须放到一个文件夹中
// 如：config=conf/dev/base.json 	ConfEnvPath=conf/dev	ConfEnv=dev
// 如：config=conf/base.json		ConfEnvPath=conf		ConfEnv=conf
func ParseConfPath(config string) error {
	path := strings.Split(config, "/")
	prefix := strings.Join(path[:len(path)-1], "/")
	ConfEnvPath = prefix
	ConfEnv = path[len(path)-2]
	return nil
}

//
func InitModule(configPath string,modules []string) error {
	if err := ParseConfPath(configPath);err!=nil{
		return err
	}

	//初始化配置文件
	if err := InitViperConf(); err != nil {
		return err
	}

	if InArrayString("auth_scope",modules) {
		if err:= InitConnClient(GetConfPath("auth_scope"));err!=nil {
			fmt.Printf("[ERROR] %s%s\n",time.Now().Format(TimeFormat),"InitConnClient:"+err.Error())
		}
	}

	if InArrayString("redis",modules) {
		if err:= InitRedisConf(GetConfPath("redis_map"));err!=nil{
          fmt.Printf("[ERROR] %s%s\n",time.Now().Format(TimeFormat),"InitRedisConf:"+err.Error())
		}
	}

	// 加载mysql配置并初始化实例
	if InArrayString("mysql", modules) {
		if err := InitDBPool(GetConfPath("mysql_map")); err != nil {
			fmt.Printf("[ERROR] %s%s\n", time.Now().Format(TimeFormat), " InitDBPool:"+err.Error())
		}
	}

	return nil
}

func InArrayString(s string,arr[] string) bool {
	for _,i :=range arr{
		if i==s{
			return true
		}
	}
	return false
}

func GetConfPath(fileName string) string {
	return ConfEnvPath + "/" + fileName + ".toml"
}