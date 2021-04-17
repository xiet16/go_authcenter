package main

import (
	"github.com/xiet16/authcenter/authserver"
	"github.com/xiet16/authcenter/common/lib"
)

func main()  {
	lib.InitModule("./conf/dev/",[]string{"auth_scope","mysql","redis"})
	lib.InitSession()

	 authserver.Run()
}