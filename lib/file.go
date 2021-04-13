package lib

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

func ParseConfig(path string,conf interface{}) error {
	file, err := os.Open(path)
	if err!=nil {
		return fmt.Errorf("Open config #{path} fail ,#{err}")
	}

	data ,err :=ioutil.ReadAll(file)
	if err!=nil {
		return fmt.Errorf("Read config fail, #{err}")
	}

	v := viper.New()
	v.SetConfigType("toml")
	v.ReadConfig(bytes.NewBuffer(data))
	if err:=v.Unmarshal(conf);err!=nil {
		return fmt.Errorf("Parse config fail, config:%v, err:%v", string(data), err)
	}

	return nil
}