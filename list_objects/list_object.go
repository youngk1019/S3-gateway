package list_objects

import (
	"encoding/xml"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"s3-gateway/command/vars"
	"strings"
)

type ListBucketResult struct {
	XMLName               xml.Name         `xml:"ListBucketResult"`
	Xmlns                 string           `xml:"xmlns,attr"`
	Name                  string           `xml:"Name"`
	CommonPrefixes        []CommonPrefixes `xml:"CommonPrefixes"`
	EncodingType          string           `xml:"EncodingType"`
	Prefix                string           `xml:"Prefix"`
	StartAfter            string           `xml:"StartAfter"`
	KeyCount              int              `xml:"KeyCount"`
	MaxKeys               int              `xml:"MaxKeys"`
	Delimiter             string           `xml:"Delimiter"`
	ContinuationToken     string           `xml:"ContinuationToken"`
	NextContinuationToken string           `xml:"NextContinuationToken"`
	IsTruncated           bool             `xml:"IsTruncated"`
	Contents              []Contents       `xml:"Contents"`
}

type CommonPrefixes struct {
	Prefix string `xml:"Prefix"`
}

type Contents struct {
	Key          string `xml:"Key"`
	LastModified string `xml:"LastModified"`
	ETag         string `xml:"ETag"`
	Size         int    `xml:"Size"`
	Owner        Owner  `xml:"Owner"`
	StorageClass string `xml:"StorageClass"`
}

type Owner struct {
	ID          string `xml:"ID"`
	DisplayName string `xml:"DisplayName"`
}

func UnmarshalListObjects(b []byte) (*ListBucketResult, error) {
	ret := &ListBucketResult{}
	err := xml.Unmarshal(b, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func GenListObjectsResult(dataset []string) (*ListBucketResult, error) {
	ret := &ListBucketResult{}
	ret.XMLName.Space = "http://s3.amazonaws.com/doc/2006-03-01/"
	ret.XMLName.Local = "ListBucketResult"
	ret.Xmlns = "http://s3.amazonaws.com/doc/2006-03-01/"
	ret.Name = vars.Bucket
	ret.Prefix = "public/"
	ret.KeyCount = len(dataset)
	ret.MaxKeys = 1000
	ret.Delimiter = "/"
	ret.IsTruncated = false
	commonPrefixes := make([]CommonPrefixes, 0)
	for _, u := range dataset {
		var commonPrefix CommonPrefixes
		commonPrefix.Prefix = "public/" + u + "/"
		commonPrefixes = append(commonPrefixes, commonPrefix)
	}
	ret.CommonPrefixes = commonPrefixes
	return ret, nil
}

type ListResult struct {
	IsTruncated           bool     `json:"IsTruncated"`
	KeyCount              int      `json:"KeyCount"`
	ContinuationToken     string   `json:"ContinuationToken"`
	NextContinuationToken string   `json:"NextContinuationToken"`
	Prefix                string   `json:"Prefix"`
	Delimiter             string   `json:"Delimiter"`
	StartAfter            string   `json:"StartAfter"`
	IsPublic              bool     `json:"IsPublic"`
	MaxKeys               int      `json:"MaxKeys"`
	Objects               []Object `json:"Objects"`
}

type Object struct {
	FullPath     string `json:"FullPath"`
	Name         string `json:"Name"`
	LastModified string `json:"LastModified"`
	Size         int    `json:"Size"`
	IsDirectory  bool   `json:"IsDirectory"`
}

var BelongErr = fmt.Errorf("list object belong match failed")

func (list *ListBucketResult) GenJson(belong string, isPublic bool) ([]byte, error) {
	ret := &ListResult{}
	if strings.HasPrefix(list.Prefix, belong) {
		ret.Prefix = list.Prefix[len(belong):]
	} else {
		return nil, BelongErr
	}
	ret.IsPublic = isPublic
	ret.ContinuationToken = list.ContinuationToken
	ret.NextContinuationToken = list.NextContinuationToken
	ret.Delimiter = list.Delimiter
	ret.MaxKeys = list.MaxKeys
	if list.StartAfter != "" {
		if strings.HasPrefix(list.StartAfter, belong) {
			ret.StartAfter = list.StartAfter[len(belong):]
		} else {
			return nil, BelongErr
		}
	}
	ret.IsTruncated = list.IsTruncated
	objects := make([]Object, 0)
	for _, u := range list.Contents {
		var obj Object
		if strings.HasPrefix(u.Key, belong) {
			obj.FullPath = u.Key[len(belong):]
		} else {
			return nil, BelongErr
		}
		if obj.FullPath == ret.Prefix {
			continue
		}
		obj.Size = u.Size
		obj.LastModified = u.LastModified
		obj.IsDirectory = false
		names := strings.Split(obj.FullPath, ret.Delimiter)
		obj.Name = names[len(names)-1]
		objects = append(objects, obj)
		ret.KeyCount++
	}
	for _, u := range list.CommonPrefixes {
		var obj Object
		if strings.HasPrefix(u.Prefix, belong) {
			obj.FullPath = u.Prefix[len(belong):]
		} else {
			return nil, BelongErr
		}
		obj.Size = 0
		obj.IsDirectory = true
		names := strings.Split(obj.FullPath, ret.Delimiter)
		obj.Name = names[len(names)-2]
		objects = append(objects, obj)
		ret.KeyCount++
	}
	ret.Objects = objects
	return jsoniter.Marshal(ret)
}
