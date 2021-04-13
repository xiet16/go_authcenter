package main

import (
	"github.com/xiet16/authcenter/authserver"
	"github.com/xiet16/authcenter/lib"
)

func main()  {
	lib.InitModule("/conf/dev/",[]string{"mysql","redis"})
	lib.InitSession()

	 authserver.Run()
}