package main

import (
	"github.com/xiet16/authcenter/common/lib"
	"github.com/xiet16/authcenter/server/identityserver"
	"os"
	"os/signal"
	"syscall"
)

func main()  {
	lib.InitModule("./conf/dev/",[]string{"auth_scope","mysql","redis"})
	lib.InitSession()
	identityserver.Run()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}