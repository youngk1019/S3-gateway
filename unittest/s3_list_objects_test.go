package unittest

import (
	"crypto/md5"
	"github.com/appleboy/gofight/v2"
	"github.com/go-playground/assert/v2"
	"net/http"
	"s3-gateway/command/vars"
	"s3-gateway/list_objects"
	"s3-gateway/log"
	"s3-gateway/routers"
	"s3-gateway/util"
	"strings"
	"testing"
)

func TestS3Handler_ListObject(t *testing.T) {
	vars.UnitTest = true
	vars.Debug = true
	log.InitLogger()
	defer log.Sync()
	routers.InitRouter()

	body := util.RandString(1024)
	md5Hash := md5.New()
	md5Hash.Write([]byte(body))

	r := gofight.New()
	r.PUT("/s3/test/aa/aa").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetBody(body).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
		})

	r = gofight.New()
	r.PUT("/s3/test/aa/bb").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetBody(body).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
		})

	r = gofight.New()
	r.PUT("/s3/test/bb").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetBody(body).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
		})

	r = gofight.New()
	r.PUT("/s3/test/cc/aa").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetBody(body).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusOK)
		})

	r = gofight.New()
	r.GET("/s3/").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetQuery(gofight.H{
			"list-type": "2",
			"delimiter": "/",
			"prefix":    "test/aa/",
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			ret, err := list_objects.UnmarshalListObjects(r.Body.Bytes())
			assert.Equal(t, err, nil)
			assert.Equal(t, strings.HasSuffix(ret.Prefix, "test/aa/"), true)
			assert.Equal(t, len(ret.Contents), 2)
			assert.Equal(t, ret.Delimiter, "/")
			assert.Equal(t, ret.ContinuationToken, "")
			assert.Equal(t, ret.IsTruncated, false)
			assert.Equal(t, strings.HasSuffix(ret.Contents[0].Key, "test/aa/aa"), true)
			assert.Equal(t, strings.HasSuffix(ret.Contents[1].Key, "test/aa/bb"), true)
		})

	r = gofight.New()
	r.GET("/s3/").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetQuery(gofight.H{
			"list-type":   "2",
			"delimiter":   "/",
			"prefix":      "test/aa/",
			"start-after": "test/aa/aa",
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			ret, err := list_objects.UnmarshalListObjects(r.Body.Bytes())
			assert.Equal(t, err, nil)
			assert.Equal(t, len(ret.Contents), 1)
			assert.Equal(t, strings.HasSuffix(ret.Prefix, "test/aa/"), true)
			assert.Equal(t, ret.Delimiter, "/")
			assert.Equal(t, ret.IsTruncated, false)
			assert.Equal(t, strings.HasSuffix(ret.StartAfter, "test/aa/aa"), true)
			assert.Equal(t, strings.HasSuffix(ret.Contents[0].Key, "test/aa/bb"), true)
		})

	r = gofight.New()
	r.GET("/s3/").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetQuery(gofight.H{
			"list-type": "2",
			"delimiter": "/",
			"prefix":    "test/aa/",
			"max-keys":  "1",
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			ret, err := list_objects.UnmarshalListObjects(r.Body.Bytes())
			assert.Equal(t, err, nil)
			assert.Equal(t, len(ret.Contents), 1)
			assert.Equal(t, strings.HasSuffix(ret.Prefix, "test/aa/"), true)
			assert.Equal(t, ret.Delimiter, "/")
			assert.Equal(t, ret.ContinuationToken, "")
			assert.Equal(t, ret.MaxKeys, 1)
			assert.Equal(t, ret.StartAfter, "")
			assert.Equal(t, ret.IsTruncated, true)
			assert.Equal(t, strings.HasSuffix(ret.Contents[0].Key, "test/aa/aa"), true)
		})

	r = gofight.New()
	r.GET("/s3/").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetQuery(gofight.H{
			"list-type":   "2",
			"delimiter":   "/",
			"prefix":      "test/",
			"start-after": "test/aa/",
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			ret, err := list_objects.UnmarshalListObjects(r.Body.Bytes())
			assert.Equal(t, err, nil)
			assert.Equal(t, len(ret.Contents), 1)
			assert.Equal(t, len(ret.CommonPrefixes), 1)
			assert.Equal(t, strings.HasSuffix(ret.Prefix, "test/"), true)
			assert.Equal(t, ret.Delimiter, "/")
			assert.Equal(t, ret.IsTruncated, false)
			assert.Equal(t, strings.HasSuffix(ret.Contents[0].Key, "test/bb"), true)
			assert.Equal(t, strings.HasSuffix(ret.CommonPrefixes[0].Prefix, "test/cc/"), true)
		})

	r = gofight.New()
	r.GET("/s3/").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetQuery(gofight.H{
			"list-type":   "2",
			"delimiter":   "/",
			"prefix":      "test/",
			"start-after": "test/bb",
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			ret, err := list_objects.UnmarshalListObjects(r.Body.Bytes())
			assert.Equal(t, err, nil)
			assert.Equal(t, len(ret.Contents), 0)
			assert.Equal(t, len(ret.CommonPrefixes), 1)
			assert.Equal(t, strings.HasSuffix(ret.Prefix, "test/"), true)
			assert.Equal(t, ret.Delimiter, "/")
			assert.Equal(t, ret.IsTruncated, false)
			assert.Equal(t, strings.HasSuffix(ret.CommonPrefixes[0].Prefix, "test/cc/"), true)
		})

	r = gofight.New()
	r.GET("/s3/").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetQuery(gofight.H{
			"list-type": "2",
			"delimiter": "/",
			"prefix":    "test/",
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			ret, err := list_objects.UnmarshalListObjects(r.Body.Bytes())
			assert.Equal(t, err, nil)
			assert.Equal(t, len(ret.Contents), 1)
			assert.Equal(t, len(ret.CommonPrefixes), 2)
			assert.Equal(t, strings.HasSuffix(ret.Prefix, "test/"), true)
			assert.Equal(t, ret.Delimiter, "/")
			assert.Equal(t, ret.IsTruncated, false)
			assert.Equal(t, strings.HasSuffix(ret.Contents[0].Key, "test/bb"), true)
			assert.Equal(t, strings.HasSuffix(ret.CommonPrefixes[0].Prefix, "test/aa/"), true)
			assert.Equal(t, strings.HasSuffix(ret.CommonPrefixes[1].Prefix, "test/cc/"), true)
		})

	continuationToken := ""
	r = gofight.New()
	r.GET("/s3/").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetQuery(gofight.H{
			"list-type": "2",
			"delimiter": "/",
			"prefix":    "test/",
			"max-keys":  "1",
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			ret, err := list_objects.UnmarshalListObjects(r.Body.Bytes())
			continuationToken = ret.NextContinuationToken
			assert.Equal(t, err, nil)
			assert.Equal(t, len(ret.Contents), 0)
			assert.Equal(t, len(ret.CommonPrefixes), 1)
			assert.Equal(t, strings.HasSuffix(ret.Prefix, "test/"), true)
			assert.Equal(t, ret.Delimiter, "/")
			assert.Equal(t, ret.IsTruncated, true)
			assert.Equal(t, ret.MaxKeys, 1)
			assert.Equal(t, strings.HasSuffix(ret.CommonPrefixes[0].Prefix, "test/aa/"), true)
		})

	r = gofight.New()
	r.GET("/s3/").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetQuery(gofight.H{
			"list-type":          "2",
			"delimiter":          "/",
			"prefix":             "test/",
			"continuation-token": continuationToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			ret, err := list_objects.UnmarshalListObjects(r.Body.Bytes())
			continuationToken = ret.ContinuationToken
			assert.Equal(t, err, nil)
			assert.Equal(t, len(ret.Contents), 1)
			assert.Equal(t, len(ret.CommonPrefixes), 1)
			assert.Equal(t, ret.ContinuationToken, continuationToken)
			assert.Equal(t, ret.NextContinuationToken, "")
			assert.Equal(t, strings.HasSuffix(ret.Prefix, "test/"), true)
			assert.Equal(t, ret.Delimiter, "/")
			assert.Equal(t, ret.IsTruncated, false)
			assert.Equal(t, strings.HasSuffix(ret.Contents[0].Key, "test/bb"), true)
			assert.Equal(t, strings.HasSuffix(ret.CommonPrefixes[0].Prefix, "test/cc/"), true)
		})

	r = gofight.New()
	r.GET("/s3/").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetQuery(gofight.H{
			"list-type":          "2",
			"delimiter":          "/",
			"prefix":             "test/",
			"continuation-token": continuationToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			ret, err := list_objects.UnmarshalListObjects(r.Body.Bytes())
			continuationToken = ret.ContinuationToken
			assert.Equal(t, err, nil)
			assert.Equal(t, len(ret.Contents), 1)
			assert.Equal(t, len(ret.CommonPrefixes), 1)
			assert.Equal(t, ret.ContinuationToken, continuationToken)
			assert.Equal(t, ret.NextContinuationToken, "")
			assert.Equal(t, strings.HasSuffix(ret.Prefix, "test/"), true)
			assert.Equal(t, ret.Delimiter, "/")
			assert.Equal(t, ret.IsTruncated, false)
			assert.Equal(t, strings.HasSuffix(ret.Contents[0].Key, "test/bb"), true)
			assert.Equal(t, strings.HasSuffix(ret.CommonPrefixes[0].Prefix, "test/cc/"), true)
		})

	r = gofight.New()
	r.DELETE("/s3/test/aa/aa").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusNoContent)
		})

	r = gofight.New()
	r.DELETE("/s3/test/aa/bb").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusNoContent)
		})

	r = gofight.New()
	r.DELETE("/s3/test/bb").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusNoContent)
		})

	r = gofight.New()
	r.DELETE("/s3/test/cc/aa").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, http.StatusNoContent)
		})

	r = gofight.New()
	r.GET("/s3/").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
			vars.PublicKey: "True",
		}).
		SetQuery(gofight.H{
			"prefix": "",
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			println(r.Body.String())
			ret, err := list_objects.UnmarshalListObjects(r.Body.Bytes())
			assert.Equal(t, err, nil)
			assert.Equal(t, len(ret.CommonPrefixes), 1)
			assert.Equal(t, ret.CommonPrefixes[0].Prefix, "public/BasicDataset/")
		})
}
