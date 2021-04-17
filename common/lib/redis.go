package lib

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/prometheus/common/log"
	"math/rand"
	"time"
)

//redis 工厂
//这里最好好实现负载均衡
func RedisConnFactory(name string)(redis.Conn,error)  {
	if ConfRedisMap !=nil&& ConfRedisMap.List !=nil {
		for confName,cfg := range ConfRedisMap.List{
			if name == confName {
				randHost := cfg.ProxyList[rand.Intn(len(cfg.ProxyList))]
				if cfg.ConnTimeout ==0 {
					cfg.ConnTimeout = 50
				}

				if cfg.ReadTimeout ==0 {
					cfg.ReadTimeout = 50
				}

				if cfg.WriteTimeout ==0 {
					cfg.WriteTimeout = 50
				}

				c,err := redis.Dial("tcp",
					randHost,
					redis.DialConnectTimeout(time.Duration(cfg.ConnTimeout)* time.Millisecond),
					redis.DialReadTimeout(time.Duration(cfg.ReadTimeout)*time.Millisecond),
					redis.DialWriteTimeout(time.Duration(cfg.WriteTimeout)*time.Millisecond))

				if err!=nil {
					return nil, err
				}

				if cfg.Password !="" {
					if _,err := c.Do("AUTH",cfg.Password);err!=nil {
						c.Close()
						return nil, err
					}
				}

				if cfg.Db !=0 {
					if _,err := c.Do("SELECT",cfg.Db);err!=nil {
						c.Close()
						return nil, err
					}
				}
				return c,nil
			}
		}
	}
	return  nil,errors.New("create redis conn fail")
}

func RedisConfDo(name string,commandName string,args ...interface{})(interface{},error) {
	c,err := RedisConnFactory(name)
	if err!=nil {
		log.Error(err)
		return nil, err
	}
	defer c.Close()
	startExecTime := time.Now()
	reply,err := c.Do(commandName,args...)
	log.Info(time.Since(startExecTime))
	if err!=nil {
		log.Error(err)
	}else {
		replyStr,_:=redis.String(reply,nil)
		log.Info(commandName,"do reply success:",replyStr)
	}
	return reply,err
}

func RedisConfFactory(name string) (addr,pwd string,err error){
	if ConfRedisMap==nil || ConfRedisMap.List==nil {
		return "","",errors.New("redis conf not found")
	}

	if cfg,ok := ConfRedisMap.List[name];ok{
		randHost := cfg.ProxyList[rand.Intn(len(cfg.ProxyList))]
		if randHost =="" {
			return "", "", errors.New("redis host not found")
		}

		return randHost,cfg.Password,nil
	}
	return "", "", errors.New("redis config error")
}

func RedisSetString(name,key,value string, expire string) error {
	_,err :=RedisConfDo(name,"SET",key,value,expire)
	return err
}