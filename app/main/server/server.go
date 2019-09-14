package server

import (
	"battery-analysis-platform/app/main/conf"
	"battery-analysis-platform/app/main/middleware"
	"github.com/gin-gonic/gin"
)

func Start() error {
	ginConf := &conf.App.Gin
	gin.SetMode(ginConf.RunMode)
	r := gin.Default()
	r.Use(middleware.Session(ginConf.SecretKey))
	register(r)
	return r.Run(conf.App.Gin.HttpAddr)
}
