package routers

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"s3-gateway/api"
	"s3-gateway/command/vars"
	"s3-gateway/log"
)

var router *gin.Engine

func InitRouter() {
	if !vars.Debug {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = ioutil.Discard
	}

	r := gin.Default()
	r.MaxMultipartMemory = 1 << 30

	//api handler
	r.Any("/s3/*object", JWT(), Logger(), api.S3Handler)

	router = r
}

func GetRouter() http.Handler {
	return router
}

func Run() {
	err := router.Run(":" + vars.Port)
	if err != nil {
		log.Errorw("gin router init failed", "error", err.Error())
	}
}
