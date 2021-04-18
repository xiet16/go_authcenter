package identityserver

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"github.com/xiet16/authcenter/controller"
	"github.com/xiet16/authcenter/middleware"
	"net/http"
	"time"
)



func Run() {
	router:=gin.Default()
	identityRouter := router.Group("/api")
	identityRouter.Use(middleware.LoggerToFile())
	{
      controller.IndentityServerRegister(identityRouter)
	}

	HttpSrvHandler := &http.Server{
		Addr:           ":9098",
		Handler:        router,
		ReadTimeout:    50 * time.Second,
		WriteTimeout:   50 * time.Second,
		MaxHeaderBytes: 1 << uint(20),
	}
	go func() {
		log.Info(" [INFO] HttpServerRun:%s\n", ":9098")
		if err := HttpSrvHandler.ListenAndServe(); err != nil {
			log.Fatalf(" [ERROR] HttpServerRun:%s err:%v\n", ":9098", err)
		}
	}()
}

