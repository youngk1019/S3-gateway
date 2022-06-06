package api

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"math/rand"
	"net"
	"net/http"
	"s3-gateway/command/vars"
	"s3-gateway/edu_backend"
	"s3-gateway/list_objects"
	"s3-gateway/log"
	"strings"
	"time"
)

func S3Handler(c *gin.Context) {
	method := c.Request.Method
	params := c.Request.URL.Query()
	object := c.Param("object")[1:]

	endpointSet := make(map[string]bool)
	endpoint := vars.EndpointList[rand.Intn(len(vars.EndpointList))]
	_, err := net.DialTimeout("tcp", endpoint, 1*time.Second)
	endpointSet[endpoint] = true
	for err != nil && len(endpointSet) < len(vars.Endpoint) {
		endpoint = vars.EndpointList[rand.Intn(len(vars.EndpointList))]
		for endpointSet[endpoint] {
			endpoint = vars.EndpointList[rand.Intn(len(vars.EndpointList))]
		}
		_, err = net.DialTimeout("tcp", endpoint, 1*time.Second)
		endpointSet[endpoint] = true
	}
	if err != nil {
		log.Errorw("s3 server error", vars.UUIDKey, c.Value(vars.UUIDKey))
		c.String(http.StatusInternalServerError, "s3 server error")
		return
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(vars.AccessKey, vars.SecretKey, ""),
		Secure: false,
		Region: "us-east-1",
	})
	if err != nil {
		log.Errorw("new minio client failed", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
		c.String(http.StatusBadGateway, "gateway new minio client failed")
		return
	}

	//var XAmzHeaders http.Header
	//for k, values := range c.Request.Header {
	//	if strings.HasPrefix(k, "x-amz-") {
	//		for _, v := range values {
	//			XAmzHeaders.Add(k, v)
	//		}
	//	}
	//}
	if strings.ToLower(c.Request.Header.Get(vars.PublicKey)) == "true" {
		if method == http.MethodPut || method == http.MethodPost {
			uname := c.Value(vars.UNAMEKey).(string)
			check := false
			for _, u := range vars.AdminList {
				if u == uname {
					check = true
					break
				}
			}
			if check == false {
				c.String(http.StatusUnauthorized, "Do not have permission for put object")
				return
			}
		}

		dataSets, err := edu_backend.GetDataSet(c, c.Value(vars.UIDKey).(string))
		if err != nil {
			log.Errorw("get data set", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
			c.String(http.StatusBadGateway, "gateway get data set")
			return
		}

		if object == "" && params.Get("prefix") == "" {
			ret, err := list_objects.GenListObjectsResult(dataSets)
			if err != nil {
				log.Errorw("get list objects result", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
				c.String(http.StatusBadGateway, "gateway get data set")
			}
			c.Data(http.StatusOK, "text/xml", ret)
			return
		}

		ret := strings.Split(object, "/")
		if ret[0] != "" {
			check := false
			for _, u := range dataSets {
				if u == ret[0] {
					check = true
					break
				}
			}
			if check == false {
				c.String(http.StatusUnauthorized, "Do not have permission for this dataset: "+ret[0])
				return
			}
			object = "public/" + object
		}

		if params.Has("prefix") {
			ret = strings.Split(params.Get("prefix"), "/")
			if ret[0] != "" {
				check := false
				for _, u := range dataSets {
					if u == ret[0] {
						check = true
						break
					}
				}
				if check == false {
					c.String(http.StatusUnauthorized, "Do not have permission for this dataset: "+ret[0])
					return
				}
			}
			prefix := "public/" + params.Get("prefix")
			params.Del("prefix")
			params.Set("prefix", prefix)
		}

		if params.Has("start-after") {
			ret = strings.Split(params.Get("start-after"), "/")
			if ret[0] != "" {
				check := false
				for _, u := range dataSets {
					if u == ret[0] {
						check = true
						break
					}
				}
				if check == false {
					c.String(http.StatusUnauthorized, "Do not have permission for this dataset: "+ret[0])
					return
				}
			}
			startAfter := "public/" + params.Get("start-after")
			params.Del("start-after")
			params.Set("start-after", startAfter)
		}

		if c.Request.Header.Get("x-amz-copy-source") != "" {
			copySource := "public/" + c.Request.Header.Get("x-amz-copy-source")
			c.Request.Header.Del("x-amz-copy-source")
			c.Request.Header.Set("x-amz-copy-source", copySource)
		}

	} else {
		if object != "" {
			object = "workplace/" + c.Value(vars.UIDKey).(string) + "/" + object
		}
		if params.Has("prefix") {
			prefix := "workplace/" + c.Value(vars.UIDKey).(string) + "/" + params.Get("prefix")
			params.Del("prefix")
			params.Set("prefix", prefix)
		}
		if params.Has("start-after") {
			startAfter := "workplace/" + c.Value(vars.UIDKey).(string) + "/" + params.Get("start-after")
			params.Del("start-after")
			params.Set("start-after", startAfter)
		}
		if c.Request.Header.Get("x-amz-copy-source") != "" {
			copySource := vars.Bucket + "/workplace/" + c.Value(vars.UIDKey).(string) + "/" + c.Request.Header.Get("x-amz-copy-source")
			c.Request.Header.Del("x-amz-copy-source")
			c.Request.Header.Set("x-amz-copy-source", copySource)
		}
	}
	c.Request.Header.Del(vars.PublicKey)

	url, err := client.PresignHeader(c, method, vars.Bucket, object, 86400*time.Second, params, c.Request.Header)
	if err != nil {
		log.Errorw("presign failed", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
		c.String(http.StatusBadGateway, "gateway presign failed")
		return
	}

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, c.Request.Body)
	if err != nil {
		log.Errorw("io copy", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
		c.String(http.StatusBadGateway, "gateway io copy")
		return
	}

	req, err := http.NewRequest(method, url.String(), buf)
	if err != nil {
		log.Errorw("new request failed", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
		c.String(http.StatusBadGateway, "gateway new request failed")
		return
	}
	for k, values := range c.Request.Header {
		for _, v := range values {
			req.Header.Add(k, v)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorw("do request", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
		c.String(http.StatusBadGateway, "gateway do request")
		return
	}

	//c.Status(resp.StatusCode)
	//_, err = io.Copy(c.Writer, resp.Body)
	//if err != nil {
	//	c.String(http.StatusBadGateway, "gateway io copy")
	//	return
	//}
	//
	//for k, values := range resp.Header {
	//	for _, v := range values {
	//		c.Writer.Header().Add(k, v)
	//	}
	//}

	header := make(map[string]string)
	for k, _ := range resp.Header {
		header[k] = resp.Header.Get(k)
	}

	c.Render(resp.StatusCode, render.Reader{
		Headers:       header,
		ContentLength: resp.ContentLength,
		Reader:        resp.Body,
	})
	return
}
