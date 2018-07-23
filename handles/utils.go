package handles

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/grt1st/wdproxy/g"
	"github.com/axgle/mahonia"
)

func RequestBody(res *http.Request) ([]byte, error) {
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body = ioutil.NopCloser(bytes.NewReader(buf))
	return buf, nil
}

func ResponseBody(res *http.Response) ([]byte, error) {
	if res == nil {
		return nil, nil
	}
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body = ioutil.NopCloser(bytes.NewReader(buf))
	return buf, nil
}

func GetContentType(HeradeCT string) string {
	ct := strings.Split(HeradeCT, "; ")[0]
	return ct
}

func GetExtension(path string) string {
	SlicePath := strings.Split(path, ".")
	if len(SlicePath) > 1 {
		return SlicePath[len(SlicePath)-1]
	}
	return ""
}

func (rr *RequestRecord) isStatic() bool {
	var mtype string
	if rr.ContentType != "" {
		mtype = strings.Split(rr.ContentType, "/")[0]
	}
	if ContainsString(static_ext, rr.Extension) {
		return true
	} else if ContainsString(static_types, rr.ContentType) {
		return true
	} else if ContainsString(media_types, mtype) {
		return true
	}
	return false
}

func ContainsString(sl []string, v string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func toJsonHeader(header http.Header) string {
	js, err := json.Marshal(header)
	if err != nil {
		log.Println(err)
	}
	return string(js)
}

type RequestRecord struct {
	// host、params、length、mime——type
	// 静态资源/cdn域名过滤
	Url            string      `json:"url"`
	Method         string      `json:"method"`
	Status         int         `json:"status"`
	Scheme         string      `json:"scheme"`
	Path           string      `json:"path"`
	ContentType    string      `json:"content_type"`
	ContentLength  uint        `json:"content_length"`
	RemoteAddr     string      `json:"remote_addr"`
	Host           string      `json:"host"`
	Port           string      `json:"port"`
	Extension      string      `json:"ext"`
	HeaderRequest  http.Header `json:"http_request"`
	HeaderResponse http.Header `json:"http_response"`
	BodyRequest    []byte      `json:"body_request"`
	BodyResponse   []byte      `json:"body_response"`
	TimeStart      time.Time   `json:"time_start"`
	TimeEnd        time.Time   `json:"time_end"`
}

func (rr *RequestRecord) Save() {
	if rr.ContentType == "application/ocsp-response" {
		//验证https证书
		//fmt.Println("get you !")
		return
	}
	var domain g.WdproxyDomain
	// 多个插入的一致性保证
	r := g.DB.Where(g.WdproxyDomain{Value: rr.Host}).FirstOrCreate(&domain)
	if r.Error != nil {
		fmt.Println("domain===", r.Error)
	}
	var bodyRequest, bodyResponse string
	if mahonia.GetCharset(string(rr.BodyRequest)) == nil{
		dec := mahonia.NewDecoder("utf8")
		bodyRequest = dec.ConvertString(string(rr.BodyRequest))
		bodyResponse = dec.ConvertString(string(rr.BodyResponse))
	}else {
		fmt.Println(rr.Url, rr.Status, rr.Extension, rr.HeaderResponse, rr.BodyRequest, rr.BodyRequest, rr.BodyResponse)
	}
	record := g.WdproxyRecord{
		Url: rr.Url,
		Method: rr.Method,
		Status: rr.Status,
		Scheme: rr.Scheme,
		Path: rr.Path,
		ContentType: rr.ContentType,
		RemoteAddr: rr.RemoteAddr,
		Host: rr.Host,
		Port: rr.Port,
		Extension: rr.Extension,
		HeaderRequest: toJsonHeader(rr.HeaderRequest),
		HeaderResponse: toJsonHeader(rr.HeaderResponse),
		BodyRequest: bodyRequest,
		BodyResponse: bodyResponse,
	}
	result := g.DB.Create(&record)
	if result.Error != nil {
		fmt.Println(result.Error)
		fmt.Println(rr.ContentType == "application/ocsp-response")
		fmt.Println(rr.ContentType)
		//fmt.Println(rr.Url, rr.Status, rr.Extension, rr.HeaderResponse, rr.BodyRequest, rr.BodyRequest, rr.BodyResponse)
	}
	g.DB.Model(&domain).Association("wdproxy_records").Append(&record)
}

func (rr *RequestRecord) String() string {
	return fmt.Sprintf(rr.Url, rr.Status, toJsonHeader(rr.HeaderRequest), toJsonHeader(rr.HeaderResponse), rr.BodyRequest, rr.BodyResponse, rr.TimeStart, rr.TimeEnd)
}

func NewRequestRecord(resp *http.Response, reqbody []byte, respbody []byte) *RequestRecord {

	var (
		ctype   string
		clength int
		StrHost string
		StrPort string
	)

	if len(resp.Header["Content-Type"]) >= 1 {
		ctype = GetContentType(resp.Header["Content-Type"][0])
	}

	if len(resp.Header["Content-Length"]) >= 1 {
		clength, _ = strconv.Atoi(resp.Header["Content-Length"][0])
	}

	SliceHost := strings.Split(resp.Request.URL.Host, ":")
	if len(SliceHost) > 1 {
		StrHost, StrPort = SliceHost[0], SliceHost[1]
	} else {
		StrHost = SliceHost[0]
		if resp.Request.URL.Scheme == "https" {
			StrPort = "443"
		} else {
			StrPort = "80"
		}
	}

	return &RequestRecord{
		Url:            resp.Request.URL.String(),
		Method:         resp.Request.Method,
		Status:         resp.StatusCode,
		Path:           resp.Request.URL.Path,
		Scheme:         resp.Request.URL.Scheme,
		ContentType:    string(ctype),
		ContentLength:  uint(clength),
		Host:           StrHost,
		Port:           StrPort,
		RemoteAddr:     resp.Request.RemoteAddr,
		Extension:      GetExtension(resp.Request.URL.Path),
		HeaderResponse: resp.Header,
		HeaderRequest:  resp.Request.Header,
		BodyResponse:   respbody,
		BodyRequest:    reqbody,
		TimeStart:      time.Now(),
		TimeEnd:        time.Now(),
	}
}
