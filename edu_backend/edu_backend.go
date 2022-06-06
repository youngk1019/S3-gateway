package edu_backend

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"s3-gateway/command/vars"
	"s3-gateway/log"
)

type dataset struct {
	Msg  string   `json:"msg"`
	Code int      `json:"code"`
	Data []string `json:"data"`
}

func GetDataSet(ctx context.Context, uid string) ([]string, error) {
	resp, err := http.Get("http://" + vars.EduBackend + "/dataset/datasettype/selectDatasetTypeNameByUserId/" + uid)
	if err != nil {
		log.Errorw("get data set", vars.UUIDKey, ctx.Value(vars.UUIDKey), "error", err.Error())
		return nil, err
	}

	httpBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorw("get http body error", vars.UUIDKey, ctx.Value(vars.UUIDKey), "error", err.Error())
		return nil, err
	}

	info := &dataset{}
	err = jsoniter.Unmarshal(httpBody, info)
	if err != nil {
		log.Errorw("json unmarshal", vars.UUIDKey, ctx.Value(vars.UUIDKey), "error", err.Error())
		return nil, err
	}

	return info.Data, nil
}
