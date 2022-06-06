package list_objects

import (
	"encoding/xml"
	"s3-gateway/command/vars"
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

func GenListObjectsResult(dataset []string) ([]byte, error) {
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

	return xml.Marshal(ret)
}
