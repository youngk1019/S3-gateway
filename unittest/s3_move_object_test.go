package unittest

import (
	"crypto/md5"
	"github.com/appleboy/gofight/v2"
	"github.com/go-playground/assert/v2"
	"net/http"
	"s3-gateway/command/vars"
	"s3-gateway/log"
	"s3-gateway/routers"
	"s3-gateway/util"
	"testing"
)

func TestS3Handler_MoveObjects(t *testing.T) {
	vars.UnitTest = true
	vars.Debug = true
	log.InitLogger()
	defer log.Sync()
	routers.InitRouter()

	body := util.RandString(1024)
	md5Hash := md5.New()
	md5Hash.Write([]byte(body))

	r := gofight.New()
	r.PUT("/s3/test/aa").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetBody(body).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
		})

	r = gofight.New()
	r.GET("/s3/test/aa").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
			md5Hash2 := md5.New()
			md5Hash2.Write(r.Body.Bytes())
			assert.Equal(t, md5Hash2.Sum(nil), md5Hash.Sum(nil))
		})

	r = gofight.New()
	r.PUT("/s3/test/bb").
		SetHeader(gofight.H{
			vars.JWTHeader:      "Bearer " + jwtToken,
			"x-amz-copy-source": "test/aa",
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
		})

	r = gofight.New()
	r.GET("/s3/test/aa").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
			md5Hash2 := md5.New()
			md5Hash2.Write(r.Body.Bytes())
			assert.Equal(t, md5Hash2.Sum(nil), md5Hash.Sum(nil))
		})

	r = gofight.New()
	r.GET("/s3/test/bb").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
			md5Hash2 := md5.New()
			md5Hash2.Write(r.Body.Bytes())
			assert.Equal(t, md5Hash2.Sum(nil), md5Hash.Sum(nil))
		})

	r = gofight.New()
	r.PUT("/s3/test/cc").
		SetHeader(gofight.H{
			vars.JWTHeader:      "Bearer " + jwtToken,
			"x-amz-copy-source": "test/aa",
		}).
		SetQuery(gofight.H{
			"move": "",
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
		})

	r = gofight.New()
	r.GET("/s3/test/aa").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusNotFound)
		})

	r = gofight.New()
	r.GET("/s3/test/bb").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
			md5Hash2 := md5.New()
			md5Hash2.Write(r.Body.Bytes())
			assert.Equal(t, md5Hash2.Sum(nil), md5Hash.Sum(nil))
		})

	r = gofight.New()
	r.GET("/s3/test/cc").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
			md5Hash2 := md5.New()
			md5Hash2.Write(r.Body.Bytes())
			assert.Equal(t, md5Hash2.Sum(nil), md5Hash.Sum(nil))
		})

	r.DELETE("/s3/test/bb").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetQuery(gofight.H{
			"recursive": "",
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusNoContent)
		})

	r = gofight.New()
	r.DELETE("/s3/test/cc").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusNoContent)
		})

	r = gofight.New()
	r.GET("/s3/test/bb").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusNotFound)
		})

	r = gofight.New()
	r.GET("/s3/test/cc").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusNotFound)
		})
}

func TestS3Handler_RecursiveMoveObjects(t *testing.T) {
	vars.UnitTest = true
	vars.Debug = true
	log.InitLogger()
	defer log.Sync()
	routers.InitRouter()

	body := util.RandString(1024)
	md5Hash := md5.New()
	md5Hash.Write([]byte(body))

	r := gofight.New()
	r.DELETE("/s3/test/").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetQuery(gofight.H{
			"recursive": "",
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusNoContent)
		})

	r = gofight.New()
	r.DELETE("/s3/test2/").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetQuery(gofight.H{
			"recursive": "",
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusNoContent)
		})

	r = gofight.New()
	r.PUT("/s3/test/1.txt").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetBody(body).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
		})

	r = gofight.New()
	r.PUT("/s3/test/2.txt").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetBody(body).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
		})

	r = gofight.New()
	r.PUT("/s3/test2/3.txt").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetBody(body).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
		})

	r = gofight.New()
	r.GET("/s3/test/1.txt").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
			md5Hash2 := md5.New()
			md5Hash2.Write(r.Body.Bytes())
			assert.Equal(t, md5Hash2.Sum(nil), md5Hash.Sum(nil))
		})

	r = gofight.New()
	r.GET("/s3/test/2.txt").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
			md5Hash2 := md5.New()
			md5Hash2.Write(r.Body.Bytes())
			assert.Equal(t, md5Hash2.Sum(nil), md5Hash.Sum(nil))
		})

	r = gofight.New()
	r.GET("/s3/test2/3.txt").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
			md5Hash2 := md5.New()
			md5Hash2.Write(r.Body.Bytes())
			assert.Equal(t, md5Hash2.Sum(nil), md5Hash.Sum(nil))
		})

	r = gofight.New()
	r.PUT("/s3/test2/").
		SetHeader(gofight.H{
			vars.JWTHeader:      "Bearer " + jwtToken,
			"x-amz-copy-source": "test/",
		}).
		SetQuery(gofight.H{
			"recursive": "",
			"move":      "",
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
		})

	r = gofight.New()
	r.GET("/s3/test/").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			println(r.Body.String())
			assert.Equal(t, r.Code, http.StatusNotFound)
		})

	r = gofight.New()
	r.GET("/s3/test2/1.txt").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
			md5Hash2 := md5.New()
			md5Hash2.Write(r.Body.Bytes())
			assert.Equal(t, md5Hash2.Sum(nil), md5Hash.Sum(nil))
		})

	r = gofight.New()
	r.GET("/s3/test2/2.txt").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
			md5Hash2 := md5.New()
			md5Hash2.Write(r.Body.Bytes())
			assert.Equal(t, md5Hash2.Sum(nil), md5Hash.Sum(nil))
		})

	r = gofight.New()
	r.GET("/s3/test2/3.txt").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
			md5Hash2 := md5.New()
			md5Hash2.Write(r.Body.Bytes())
			assert.Equal(t, md5Hash2.Sum(nil), md5Hash.Sum(nil))
		})

	r = gofight.New()
	r.GET("/s3/test/1.txt").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusNotFound)
		})

	r = gofight.New()
	r.GET("/s3/test/2.txt").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusNotFound)
		})

	r = gofight.New()
	r.GET("/s3/test/").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusNotFound)
		})

	r = gofight.New()
	r.DELETE("/s3/test/").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetQuery(gofight.H{
			"recursive": "",
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusNoContent)
		})

	r = gofight.New()
	r.DELETE("/s3/test2/").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetQuery(gofight.H{
			"recursive": "",
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusNoContent)
		})

	r = gofight.New()
	r.GET("/s3/test/1.txt").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusNotFound)
		})

	r = gofight.New()
	r.GET("/s3/test/2.txt").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusNotFound)
		})

	r = gofight.New()
	r.GET("/s3/test2/1.txt").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusNotFound)
		})

	r = gofight.New()
	r.GET("/s3/test2/2.txt").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusNotFound)
		})

	r = gofight.New()
	r.GET("/s3/test2/3.txt").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusNotFound)
		})
}
