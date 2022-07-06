package api

import "C"
import (
	"bytes"
	"compress/gzip"
	"context"
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
	"s3-gateway/full_path"
	"s3-gateway/list_objects"
	"s3-gateway/log"
	"s3-gateway/util"
	"strings"
	"time"
)

func S3Handler(c *gin.Context) {
	method := c.Request.Method
	params := c.Request.URL.Query()
	object := c.Param("object")[1:]
	fp := object

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
	belong := ""
	isPublic := false
	if strings.ToLower(c.Request.Header.Get(vars.PublicKey)) == "true" {
		isPublic = true
		if method != http.MethodGet {
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

		belong = "public/"

		dataSets, err := edu_backend.GetDataSet(c, c.Value(vars.UIDKey).(string))
		if err != nil {
			log.Errorw("get data set", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
			c.String(http.StatusBadGateway, "gateway get data set")
			return
		}

		if object == "" && params.Get("prefix") == "" {
			listRet, err := list_objects.GenListObjectsResult(dataSets)
			if err != nil {
				log.Errorw("get list objects result", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
				c.String(http.StatusBadGateway, "gateway get data set")
			}
			json, err := listRet.GenJson(belong, isPublic)
			if err != nil {
				log.Errorw("gen list objects json", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
				c.String(http.StatusBadGateway, "gen list objects json")
			}
			c.Data(http.StatusOK, "application/json", json)
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
		}

	} else {
		belong = "workplace/" + c.Value(vars.UIDKey).(string) + "/"
	}
	c.Request.Header.Del(vars.PublicKey)

	if object != "" {
		object = belong + object
	}
	if params.Has("prefix") {
		prefix := belong + params.Get("prefix")
		params.Del("prefix")
		params.Set("prefix", prefix)
	}
	if params.Has("start-after") {
		startAfter := belong + params.Get("start-after")
		params.Del("start-after")
		params.Set("start-after", startAfter)
	}
	copySrc := ""
	if c.Request.Header.Get("x-amz-copy-source") != "" {
		if params.Has("move") || params.Has("recursive") {
			copySrc = belong + c.Request.Header.Get("x-amz-copy-source")
		}
		copySource := vars.Bucket + "/" + belong + c.Request.Header.Get("x-amz-copy-source")
		c.Request.Header.Del("x-amz-copy-source")
		c.Request.Header.Set("x-amz-copy-source", copySource)
	}

	if method == http.MethodGet && params.Has("is-exist") {
		_, err := client.StatObject(c, vars.Bucket, object, minio.StatObjectOptions{})
		if err != nil {
			if strings.HasSuffix(object, "/") {
				for _ = range client.ListObjects(c, vars.Bucket, minio.ListObjectsOptions{
					Prefix:  object,
					MaxKeys: 50,
				}) {
					c.String(http.StatusOK, "true")
					return
				}
			}
			c.String(http.StatusOK, "false")
		} else {
			c.String(http.StatusOK, "true")
		}
		return
	}

	if method == http.MethodGet && params.Has("prefix") && params.Has("tree") {
		delimiter := params.Get("delimiter")
		prefix := params.Get("prefix")
		prefix2 := prefix[len(belong):]

		objectsCh := make(chan minio.ObjectInfo)
		err := error(nil)
		go func() {
			defer close(objectsCh)
			// List all objects from a bucket-name with a matching prefix.
			for object := range client.ListObjects(c, vars.Bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: true}) {
				if object.Err != nil {
					log.Errorw("s3 gateway recursive delete list objects error", vars.UUIDKey, c.Value(vars.UUIDKey), "error", object.Err.Error())
					err = object.Err
					return
				}
				objectsCh <- object
			}
		}()

		names := strings.Split(prefix2, delimiter)
		name := ""
		if names[len(names)-1] == "" && len(names) >= 2 {
			name = names[len(names)-2]
		} else {
			name = names[len(names)-1]
		}

		trie := util.NewTrie[string](name)

		for obj := range objectsCh {
			key := strings.Split(obj.Key[len(prefix):], delimiter)
			ret := make([]string, 0)
			for _, u := range key {
				if u != "" {
					ret = append(ret, u)
				}
			}
			trie.Insert(key[len(key)-1] != "", ret...)
		}

		if err != nil {
			log.Errorw("list objects", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
			c.String(http.StatusBadGateway, "gateway list objects")
		}

		if strings.HasSuffix(prefix2, name) {
			prefix2 = prefix2[:len(prefix2)-len(name)]
		} else if len(prefix2)-len(name)-len(delimiter) >= 0 {
			prefix2 = prefix2[:len(prefix2)-len(name)-len(delimiter)]
		}

		ret := list_objects.BuildObjectTree(trie, prefix2, delimiter)
		c.JSON(
			http.StatusOK,
			gin.H{
				"data": ret,
			},
		)
		return
	}

	if method == http.MethodGet && params.Has("calc-sum") {
		objectsCh := make(chan minio.ObjectInfo)
		err := error(nil)
		if object == "" {
			object = belong
		}

		// Send object names that are needed to be removed to objectsCh
		go func() {
			defer close(objectsCh)
			// List all objects from a bucket-name with a matching prefix.
			for object := range client.ListObjects(c, vars.Bucket, minio.ListObjectsOptions{Prefix: object, Recursive: true}) {
				if object.Err != nil {
					log.Errorw("s3 gateway recursive delete list objects error", vars.UUIDKey, c.Value(vars.UUIDKey), "error", object.Err.Error())
					err = object.Err
					return
				}
				objectsCh <- object
			}
		}()

		fileSize := int64(0)
		fileNum := int64(0)
		dirNum := int64(0)

		for obj := range objectsCh {
			if strings.HasSuffix(obj.Key, "/") && obj.Size == 0 {
				dirNum++
			} else {
				fileNum++
			}
			fileSize += obj.Size
		}

		if err != nil {
			log.Errorw("list objects", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
			c.String(http.StatusBadGateway, "gateway list objects")
		}

		c.JSON(http.StatusOK,
			gin.H{
				"FileSize": fileSize,
				"FileNum":  fileNum,
				"DirNum":   dirNum,
			})
		return
	}

	if (method == http.MethodPut && !params.Has("uploadId")) || (method == http.MethodPost && params.Has("uploadId")) {
		dirs := full_path.SplitFullPath(fp)
		for i, u := range dirs {
			_, err := client.StatObject(c, vars.Bucket, belong+u, minio.StatObjectOptions{})
			if err == nil {
				err := client.RemoveObject(c, vars.Bucket, belong+u, minio.RemoveObjectOptions{GovernanceBypass: true})
				if err != nil {
					log.Errorw("delete object", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
					c.String(http.StatusBadGateway, "gateway delete object")
					return
				}
			}

			if i != len(dirs)-1 {
				_, err := client.StatObject(c, vars.Bucket, belong+u+"/", minio.StatObjectOptions{})
				if err != nil {
					_, err := client.PutObject(c, vars.Bucket, belong+u+"/", nil, 0, minio.PutObjectOptions{})
					if err != nil {
						log.Errorw("put object", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
						c.String(http.StatusBadGateway, "gateway put object")
						return
					}
				}
			}
		}

		err := client.RemoveObject(c, vars.Bucket, belong+fp, minio.RemoveObjectOptions{GovernanceBypass: true})
		if err != nil {
			log.Errorw("delete object", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
			c.String(http.StatusBadGateway, "gateway delete object")
			return
		}

		if !strings.HasSuffix(fp, "/") {
			err := recursiveDelete(c, client, belong+fp+"/", "/")
			if err != nil {
				c.String(http.StatusBadGateway, "gateway recursive delete error")
				return
			}
		}
	}

	if params.Has("recursive") && method == http.MethodDelete {
		if object == "" {
			object = belong
		}

		err := recursiveDelete(c, client, object, params.Get("delimiter"))
		if err != nil {
			c.String(http.StatusBadGateway, "gateway recursive delete error")
			return
		}

		c.Status(http.StatusNoContent)
		return
	}

	if params.Has("recursive") && method == http.MethodPut && c.Request.Header.Get("x-amz-copy-source") != "" {
		objectsCh := make(chan minio.ObjectInfo)
		err := error(nil)
		if object == "" {
			object = belong
		}

		// Send object names that are needed to be removed to objectsCh
		go func() {
			defer close(objectsCh)
			// List all objects from a bucket-name with a matching prefix.
			for object := range client.ListObjects(c, vars.Bucket, minio.ListObjectsOptions{Prefix: copySrc, Recursive: true}) {
				if object.Err != nil {
					log.Errorw("s3 gateway recursive delete list objects error", vars.UUIDKey, c.Value(vars.UUIDKey), "error", object.Err.Error())
					err = object.Err
					return
				}
				objectsCh <- object
			}
		}()

		for obj := range objectsCh {
			if !strings.HasPrefix(obj.Key, copySrc) {
				log.Errorw("s3 gateway copy error prefix", vars.UUIDKey, c.Value(vars.UUIDKey))
				c.String(http.StatusBadGateway, "gateway copy error")
				return
			}
			suffix := obj.Key[len(copySrc):]
			_, err := client.CopyObject(c, minio.CopyDestOptions{Bucket: vars.Bucket, Object: object + suffix}, minio.CopySrcOptions{Bucket: vars.Bucket, Object: obj.Key})
			if err != nil {
				log.Errorw("s3 gateway copy error", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
				c.String(http.StatusBadGateway, "gateway copy error")
				return
			}
		}

		if err != nil {
			log.Errorw("s3 gateway recursive copy", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
			c.String(http.StatusBadGateway, "gateway recursive copy error")
			return
		}

		if params.Has("move") {
			err := recursiveDelete(c, client, copySrc, params.Get("delimiter"))
			if err != nil {
				c.String(http.StatusBadGateway, "gateway recursive delete error")
				return
			}
		}

		if strings.HasSuffix(object, "/") {
			_, err := client.StatObject(c, vars.Bucket, object, minio.StatObjectOptions{})
			if err != nil {
				_, err := client.PutObject(c, vars.Bucket, object, nil, 0, minio.PutObjectOptions{})
				if err != nil {
					log.Errorw("put object", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
					c.String(http.StatusBadGateway, "gateway put object")
					return
				}
			}
		}

		c.Status(http.StatusOK)
		return
	}

	params.Del("move")
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

	if copySrc != "" && resp.StatusCode == http.StatusOK && method == http.MethodPut {
		err := client.RemoveObject(c, vars.Bucket, copySrc, minio.RemoveObjectOptions{GovernanceBypass: true})
		if err != nil {
			log.Errorw("remove object", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
			c.String(http.StatusBadGateway, "gateway move copy object")
			return
		}
	}

	if method == http.MethodGet && params.Has("prefix") && params.Get("list-type") == "2" && resp.StatusCode == http.StatusOK {
		var reader io.Reader
		if resp.Header.Get("Content-Encoding") == "gzip" {
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				log.Errorw("new gzip reader", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
				c.String(http.StatusBadGateway, "new gzip reader")
				return
			}
		} else {
			reader = resp.Body
		}

		b, err := io.ReadAll(reader)
		if err != nil {
			log.Errorw("read http body", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
			c.String(http.StatusBadGateway, "gateway read http body")
			return
		}
		objects, err := list_objects.UnmarshalListObjects(b)
		if err != nil {
			log.Errorw("unmarshal list objects xml", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
			c.String(http.StatusBadGateway, "unmarshal list objects xml")
			return
		}
		json, err := objects.GenJson(belong, isPublic)
		if err != nil {
			log.Errorw("gen list objects json", vars.UUIDKey, c.Value(vars.UUIDKey), "error", err.Error())
			c.String(http.StatusBadGateway, "gen list objects json")
			return
		}
		c.Data(http.StatusOK, "application/json", json)
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

func recursiveDelete(ctx context.Context, client *minio.Client, prefix string, delimiter string) error {
	if strings.HasSuffix(prefix, delimiter) {
		err := client.RemoveObject(ctx, vars.Bucket, prefix[:len(prefix)-1-len(delimiter)], minio.RemoveObjectOptions{GovernanceBypass: true})
		if err != nil {
			log.Errorw("s3 gateway delete same name object", vars.UUIDKey, ctx.Value(vars.UUIDKey), "error", err.Error())
			return err
		}
	}

	objectsCh := make(chan minio.ObjectInfo)
	err := error(nil)

	// Send object names that are needed to be removed to objectsCh
	go func() {
		defer close(objectsCh)
		// List all objects from a bucket-name with a matching prefix.
		for object := range client.ListObjects(ctx, vars.Bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: true}) {
			if object.Err != nil {
				log.Errorw("s3 gateway recursive delete list objects error", vars.UUIDKey, ctx.Value(vars.UUIDKey), "error", object.Err.Error())
				err = object.Err
				return
			}
			objectsCh <- object
		}
	}()

	for rErr := range client.RemoveObjects(ctx, vars.Bucket, objectsCh, minio.RemoveObjectsOptions{GovernanceBypass: true}) {
		log.Errorw("s3 gateway detected during deletion", vars.UUIDKey, ctx.Value(vars.UUIDKey), "error", rErr.Err.Error())
		return rErr.Err
	}

	if err != nil {
		log.Errorw("s3 gateway recursive delete list objects error", vars.UUIDKey, ctx.Value(vars.UUIDKey), "error", err.Error())
		return err
	}

	return nil
}
