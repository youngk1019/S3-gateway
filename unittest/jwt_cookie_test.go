package unittest

import (
	"github.com/appleboy/gofight/v2"
	"github.com/go-playground/assert/v2"
	"net/http"
	"s3-gateway/command/vars"
	"s3-gateway/log"
	"s3-gateway/routers"
	"testing"
)

func TestS3Handler_JWTCookie(t *testing.T) {
	vars.UnitTest = true
	vars.Debug = true
	log.InitLogger()
	defer log.Sync()
	r := gofight.New()
	routers.InitRouter()

	r.GET("/s3/").
		SetCookie(gofight.H{
			vars.JWTCookie: jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
		})
}
