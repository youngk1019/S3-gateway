package unittest

import (
	"crypto/md5"
	"github.com/appleboy/gofight/v2"
	"github.com/go-playground/assert/v2"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"s3-gateway/command/vars"
	"s3-gateway/list_objects"
	"s3-gateway/log"
	"s3-gateway/routers"
	"s3-gateway/util"
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
			ret := &list_objects.ListResult{}
			err := jsoniter.Unmarshal(r.Body.Bytes(), ret)
			assert.Equal(t, err, nil)
			assert.Equal(t, ret.Prefix, "test/aa/")
			assert.Equal(t, ret.NextContinuationToken, "")
			assert.Equal(t, ret.IsTruncated, false)
			assert.Equal(t, ret.Delimiter, "/")
			assert.Equal(t, len(ret.Objects), 2)
			assert.Equal(t, ret.KeyCount, 2)
			assert.Equal(t, ret.IsPublic, false)
			lists := []list_objects.Object{
				{FullPath: "test/aa/aa", Name: "aa", IsDirectory: false},
				{FullPath: "test/aa/bb", Name: "bb", IsDirectory: false},
			}
			for i, u := range lists {
				assert.Equal(t, ret.Objects[i].FullPath, u.FullPath)
				assert.Equal(t, ret.Objects[i].Name, u.Name)
				assert.Equal(t, ret.Objects[i].IsDirectory, u.IsDirectory)
			}
			//ret, err := list_objects.UnmarshalListObjects(r.Body.Bytes())
			//assert.Equal(t, err, nil)
			//assert.Equal(t, strings.HasSuffix(ret.Prefix, "test/aa/"), true)
			//assert.Equal(t, len(ret.Contents), 2)
			//assert.Equal(t, ret.Delimiter, "/")
			//assert.Equal(t, ret.ContinuationToken, "")
			//assert.Equal(t, ret.IsTruncated, false)
			//assert.Equal(t, strings.HasSuffix(ret.Contents[0].Key, "test/aa/aa"), true)
			//assert.Equal(t, strings.HasSuffix(ret.Contents[1].Key, "test/aa/bb"), true)
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
			ret := &list_objects.ListResult{}
			err := jsoniter.Unmarshal(r.Body.Bytes(), ret)
			assert.Equal(t, err, nil)
			assert.Equal(t, ret.Prefix, "test/aa/")
			assert.Equal(t, ret.NextContinuationToken, "")
			assert.Equal(t, ret.IsTruncated, false)
			assert.Equal(t, ret.Delimiter, "/")
			assert.Equal(t, ret.StartAfter, "test/aa/aa")
			assert.Equal(t, len(ret.Objects), 1)
			assert.Equal(t, ret.KeyCount, 1)
			assert.Equal(t, ret.IsPublic, false)
			lists := []list_objects.Object{
				{FullPath: "test/aa/bb", Name: "bb", IsDirectory: false},
			}
			for i, u := range lists {
				assert.Equal(t, ret.Objects[i].FullPath, u.FullPath)
				assert.Equal(t, ret.Objects[i].Name, u.Name)
				assert.Equal(t, ret.Objects[i].IsDirectory, u.IsDirectory)
			}
			//ret, err := list_objects.UnmarshalListObjects(r.Body.Bytes())
			//assert.Equal(t, err, nil)
			//assert.Equal(t, len(ret.Contents), 1)
			//assert.Equal(t, strings.HasSuffix(ret.Prefix, "test/aa/"), true)
			//assert.Equal(t, ret.Delimiter, "/")
			//assert.Equal(t, ret.IsTruncated, false)
			//assert.Equal(t, strings.HasSuffix(ret.StartAfter, "test/aa/aa"), true)
			//assert.Equal(t, strings.HasSuffix(ret.Contents[0].Key, "test/aa/bb"), true)
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
			"max-keys":  "2",
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			ret := &list_objects.ListResult{}
			err := jsoniter.Unmarshal(r.Body.Bytes(), ret)
			assert.Equal(t, err, nil)
			assert.Equal(t, ret.Prefix, "test/aa/")
			assert.Equal(t, ret.NextContinuationToken != "", true)
			assert.Equal(t, ret.IsTruncated, true)
			assert.Equal(t, ret.Delimiter, "/")
			assert.Equal(t, ret.StartAfter, "")
			assert.Equal(t, len(ret.Objects), 1)
			assert.Equal(t, ret.KeyCount, 1)
			assert.Equal(t, ret.MaxKeys, 2)
			assert.Equal(t, ret.IsPublic, false)
			lists := []list_objects.Object{
				{FullPath: "test/aa/aa", Name: "aa", IsDirectory: false},
			}
			for i, u := range lists {
				assert.Equal(t, ret.Objects[i].FullPath, u.FullPath)
				assert.Equal(t, ret.Objects[i].Name, u.Name)
				assert.Equal(t, ret.Objects[i].IsDirectory, u.IsDirectory)
			}
			//ret, err := list_objects.UnmarshalListObjects(r.Body.Bytes())
			//assert.Equal(t, err, nil)
			//assert.Equal(t, len(ret.Contents), 1)
			//assert.Equal(t, strings.HasSuffix(ret.Prefix, "test/aa/"), true)
			//assert.Equal(t, ret.Delimiter, "/")
			//assert.Equal(t, ret.ContinuationToken, "")
			//assert.Equal(t, ret.MaxKeys, 1)
			//assert.Equal(t, ret.StartAfter, "")
			//assert.Equal(t, ret.IsTruncated, true)
			//assert.Equal(t, strings.HasSuffix(ret.Contents[0].Key, "test/aa/aa"), true)
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
			ret := &list_objects.ListResult{}
			err := jsoniter.Unmarshal(r.Body.Bytes(), ret)
			assert.Equal(t, err, nil)
			assert.Equal(t, ret.Prefix, "test/")
			assert.Equal(t, ret.NextContinuationToken, "")
			assert.Equal(t, ret.IsTruncated, false)
			assert.Equal(t, ret.Delimiter, "/")
			assert.Equal(t, ret.StartAfter, "test/aa/")
			assert.Equal(t, len(ret.Objects), 2)
			assert.Equal(t, ret.KeyCount, 2)
			assert.Equal(t, ret.IsPublic, false)
			lists := []list_objects.Object{
				{FullPath: "test/bb", Name: "bb", IsDirectory: false},
				{FullPath: "test/cc/", Name: "cc", IsDirectory: true},
			}
			for i, u := range lists {
				assert.Equal(t, ret.Objects[i].FullPath, u.FullPath)
				assert.Equal(t, ret.Objects[i].Name, u.Name)
				assert.Equal(t, ret.Objects[i].IsDirectory, u.IsDirectory)
			}
			//ret, err := list_objects.UnmarshalListObjects(r.Body.Bytes())
			//assert.Equal(t, err, nil)
			//assert.Equal(t, len(ret.Contents), 1)
			//assert.Equal(t, len(ret.CommonPrefixes), 1)
			//assert.Equal(t, strings.HasSuffix(ret.Prefix, "test/"), true)
			//assert.Equal(t, ret.Delimiter, "/")
			//assert.Equal(t, ret.IsTruncated, false)
			//assert.Equal(t, strings.HasSuffix(ret.Contents[0].Key, "test/bb"), true)
			//assert.Equal(t, strings.HasSuffix(ret.CommonPrefixes[0].Prefix, "test/cc/"), true)
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
			ret := &list_objects.ListResult{}
			err := jsoniter.Unmarshal(r.Body.Bytes(), ret)
			assert.Equal(t, err, nil)
			assert.Equal(t, ret.Prefix, "test/")
			assert.Equal(t, ret.NextContinuationToken, "")
			assert.Equal(t, ret.IsTruncated, false)
			assert.Equal(t, ret.Delimiter, "/")
			assert.Equal(t, ret.StartAfter, "test/bb")
			assert.Equal(t, len(ret.Objects), 1)
			assert.Equal(t, ret.KeyCount, 1)
			assert.Equal(t, ret.IsPublic, false)
			lists := []list_objects.Object{
				{FullPath: "test/cc/", Name: "cc", IsDirectory: true},
			}
			for i, u := range lists {
				assert.Equal(t, ret.Objects[i].FullPath, u.FullPath)
				assert.Equal(t, ret.Objects[i].Name, u.Name)
				assert.Equal(t, ret.Objects[i].IsDirectory, u.IsDirectory)
			}
			//ret, err := list_objects.UnmarshalListObjects(r.Body.Bytes())
			//assert.Equal(t, err, nil)
			//assert.Equal(t, len(ret.Contents), 0)
			//assert.Equal(t, len(ret.CommonPrefixes), 1)
			//assert.Equal(t, strings.HasSuffix(ret.Prefix, "test/"), true)
			//assert.Equal(t, ret.Delimiter, "/")
			//assert.Equal(t, ret.IsTruncated, false)
			//assert.Equal(t, strings.HasSuffix(ret.CommonPrefixes[0].Prefix, "test/cc/"), true)
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
			ret := &list_objects.ListResult{}
			err := jsoniter.Unmarshal(r.Body.Bytes(), ret)
			assert.Equal(t, err, nil)
			assert.Equal(t, ret.Prefix, "test/")
			assert.Equal(t, ret.NextContinuationToken, "")
			assert.Equal(t, ret.IsTruncated, false)
			assert.Equal(t, ret.Delimiter, "/")
			assert.Equal(t, len(ret.Objects), 3)
			assert.Equal(t, ret.KeyCount, 3)
			assert.Equal(t, ret.IsPublic, false)
			lists := []list_objects.Object{
				{FullPath: "test/bb", Name: "bb", IsDirectory: false},
				{FullPath: "test/aa/", Name: "aa", IsDirectory: true},
				{FullPath: "test/cc/", Name: "cc", IsDirectory: true},
			}
			for i, u := range lists {
				assert.Equal(t, ret.Objects[i].FullPath, u.FullPath)
				assert.Equal(t, ret.Objects[i].Name, u.Name)
				assert.Equal(t, ret.Objects[i].IsDirectory, u.IsDirectory)
			}
			//ret, err := list_objects.UnmarshalListObjects(r.Body.Bytes())
			//assert.Equal(t, err, nil)
			//assert.Equal(t, len(ret.Contents), 1)
			//assert.Equal(t, len(ret.CommonPrefixes), 2)
			//assert.Equal(t, strings.HasSuffix(ret.Prefix, "test/"), true)
			//assert.Equal(t, ret.Delimiter, "/")
			//assert.Equal(t, ret.IsTruncated, false)
			//assert.Equal(t, strings.HasSuffix(ret.Contents[0].Key, "test/bb"), true)
			//assert.Equal(t, strings.HasSuffix(ret.CommonPrefixes[0].Prefix, "test/aa/"), true)
			//assert.Equal(t, strings.HasSuffix(ret.CommonPrefixes[1].Prefix, "test/cc/"), true)
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
			"max-keys":  "2",
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			ret := &list_objects.ListResult{}
			err := jsoniter.Unmarshal(r.Body.Bytes(), ret)
			assert.Equal(t, err, nil)
			continuationToken = ret.NextContinuationToken
			assert.Equal(t, ret.Prefix, "test/")
			assert.Equal(t, ret.NextContinuationToken != "", true)
			assert.Equal(t, ret.IsTruncated, true)
			assert.Equal(t, ret.Delimiter, "/")
			assert.Equal(t, len(ret.Objects), 1)
			assert.Equal(t, ret.KeyCount, 1)
			assert.Equal(t, ret.IsPublic, false)
			lists := []list_objects.Object{
				{FullPath: "test/aa/", Name: "aa", IsDirectory: true},
			}
			for i, u := range lists {
				assert.Equal(t, ret.Objects[i].FullPath, u.FullPath)
				assert.Equal(t, ret.Objects[i].Name, u.Name)
				assert.Equal(t, ret.Objects[i].IsDirectory, u.IsDirectory)
			}
			//ret, err := list_objects.UnmarshalListObjects(r.Body.Bytes())
			//continuationToken = ret.NextContinuationToken
			//assert.Equal(t, err, nil)
			//assert.Equal(t, len(ret.Contents), 0)
			//assert.Equal(t, len(ret.CommonPrefixes), 1)
			//assert.Equal(t, strings.HasSuffix(ret.Prefix, "test/"), true)
			//assert.Equal(t, ret.Delimiter, "/")
			//assert.Equal(t, ret.IsTruncated, true)
			//assert.Equal(t, ret.MaxKeys, 1)
			//assert.Equal(t, strings.HasSuffix(ret.CommonPrefixes[0].Prefix, "test/aa/"), true)
		})

	r = gofight.New()
	r.GET("/s3/").
		SetHeader(gofight.H{
			vars.JWTHeader:    "Bearer " + jwtToken,
			"Accept-Encoding": "gzip, deflate, br",
		}).
		SetQuery(gofight.H{
			"list-type":          "2",
			"delimiter":          "/",
			"prefix":             "test/",
			"continuation-token": continuationToken,
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			ret := &list_objects.ListResult{}
			err := jsoniter.Unmarshal(r.Body.Bytes(), ret)
			assert.Equal(t, err, nil)
			assert.Equal(t, ret.Prefix, "test/")
			assert.Equal(t, ret.ContinuationToken, continuationToken)
			assert.Equal(t, ret.NextContinuationToken, "")
			assert.Equal(t, ret.IsTruncated, false)
			assert.Equal(t, ret.Delimiter, "/")
			assert.Equal(t, len(ret.Objects), 2)
			assert.Equal(t, ret.KeyCount, 2)
			assert.Equal(t, ret.IsPublic, false)
			lists := []list_objects.Object{
				{FullPath: "test/bb", Name: "bb", IsDirectory: false},
				{FullPath: "test/cc/", Name: "cc", IsDirectory: true},
			}
			for i, u := range lists {
				assert.Equal(t, ret.Objects[i].FullPath, u.FullPath)
				assert.Equal(t, ret.Objects[i].Name, u.Name)
				assert.Equal(t, ret.Objects[i].IsDirectory, u.IsDirectory)
			}
			//ret, err := list_objects.UnmarshalListObjects(r.Body.Bytes())
			//continuationToken = ret.ContinuationToken
			//assert.Equal(t, err, nil)
			//assert.Equal(t, len(ret.Contents), 1)
			//assert.Equal(t, len(ret.CommonPrefixes), 1)
			//assert.Equal(t, ret.ContinuationToken, continuationToken)
			//assert.Equal(t, ret.NextContinuationToken, "")
			//assert.Equal(t, strings.HasSuffix(ret.Prefix, "test/"), true)
			//assert.Equal(t, ret.Delimiter, "/")
			//assert.Equal(t, ret.IsTruncated, false)
			//assert.Equal(t, strings.HasSuffix(ret.Contents[0].Key, "test/bb"), true)
			//assert.Equal(t, strings.HasSuffix(ret.CommonPrefixes[0].Prefix, "test/cc/"), true)
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
			ret := &list_objects.ListResult{}
			err := jsoniter.Unmarshal(r.Body.Bytes(), ret)
			assert.Equal(t, err, nil)
			assert.Equal(t, ret.Prefix, "test/")
			assert.Equal(t, ret.ContinuationToken, continuationToken)
			assert.Equal(t, ret.NextContinuationToken, "")
			assert.Equal(t, ret.IsTruncated, false)
			assert.Equal(t, ret.Delimiter, "/")
			assert.Equal(t, len(ret.Objects), 2)
			assert.Equal(t, ret.KeyCount, 2)
			assert.Equal(t, ret.IsPublic, false)
			lists := []list_objects.Object{
				{FullPath: "test/bb", Name: "bb", IsDirectory: false},
				{FullPath: "test/cc/", Name: "cc", IsDirectory: true},
			}
			for i, u := range lists {
				assert.Equal(t, ret.Objects[i].FullPath, u.FullPath)
				assert.Equal(t, ret.Objects[i].Name, u.Name)
				assert.Equal(t, ret.Objects[i].IsDirectory, u.IsDirectory)
			}
			//ret, err := list_objects.UnmarshalListObjects(r.Body.Bytes())
			//continuationToken = ret.ContinuationToken
			//assert.Equal(t, err, nil)
			//assert.Equal(t, len(ret.Contents), 1)
			//assert.Equal(t, len(ret.CommonPrefixes), 1)
			//assert.Equal(t, ret.ContinuationToken, continuationToken)
			//assert.Equal(t, ret.NextContinuationToken, "")
			//assert.Equal(t, strings.HasSuffix(ret.Prefix, "test/"), true)
			//assert.Equal(t, ret.Delimiter, "/")
			//assert.Equal(t, ret.IsTruncated, false)
			//assert.Equal(t, strings.HasSuffix(ret.Contents[0].Key, "test/bb"), true)
			//assert.Equal(t, strings.HasSuffix(ret.CommonPrefixes[0].Prefix, "test/cc/"), true)
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
			ret := &list_objects.ListResult{}
			err := jsoniter.Unmarshal(r.Body.Bytes(), ret)
			assert.Equal(t, err, nil)
			assert.Equal(t, ret.Prefix, "")
			assert.Equal(t, ret.ContinuationToken, "")
			assert.Equal(t, ret.NextContinuationToken, "")
			assert.Equal(t, ret.IsTruncated, false)
			assert.Equal(t, ret.Delimiter, "/")
			assert.Equal(t, len(ret.Objects), 1)
			assert.Equal(t, ret.KeyCount, 1)
			assert.Equal(t, ret.IsPublic, true)
			lists := []list_objects.Object{
				{FullPath: "BasicDatasets/", Name: "BasicDatasets", IsDirectory: true},
			}
			for i, u := range lists {
				assert.Equal(t, ret.Objects[i].FullPath, u.FullPath)
				assert.Equal(t, ret.Objects[i].Name, u.Name)
				assert.Equal(t, ret.Objects[i].IsDirectory, u.IsDirectory)
			}
			//println(r.Body.String())
			//ret, err := list_objects.UnmarshalListObjects(r.Body.Bytes())
			//assert.Equal(t, err, nil)
			//assert.Equal(t, len(ret.CommonPrefixes), 1)
			//assert.Equal(t, ret.CommonPrefixes[0].Prefix, "public/BasicDatasets/")
		})
}

func TestS3Handler_ObjectTree(t *testing.T) {
	vars.UnitTest = true
	vars.Debug = true
	log.InitLogger()
	defer log.Sync()
	routers.InitRouter()

	r := gofight.New()
	r.GET("/s3/").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetQuery(gofight.H{
			"tree":      "",
			"prefix":    "",
			"delimiter": "/",
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			println(r.Body.String())
		})

	r = gofight.New()
	r.GET("/s3/").
		SetHeader(gofight.H{
			vars.JWTHeader: "Bearer " + jwtToken,
		}).
		SetQuery(gofight.H{
			"tree":      "",
			"prefix":    "aaa/",
			"delimiter": "/",
		}).
		Run(routers.GetRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			println(r.Body.String())
		})
}
