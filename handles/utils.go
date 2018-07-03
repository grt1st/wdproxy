package handles

import (
	"io/ioutil"
	"bytes"
	"net/http"
	"time"
	"strconv"
	"strings"
	"log"
	"encoding/json"
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
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body = ioutil.NopCloser(bytes.NewReader(buf))
	return buf, nil
}

func NewParserHTTP(resp *http.Response, reqbody []byte, respbody []byte) *ParserHTTP {
	return &ParserHTTP{r: resp, reqbody: reqbody, respbody: respbody, s: time.Now()}
}

type ParserHTTP struct {
	r        *http.Response
	reqbody  []byte
	respbody []byte
	s        time.Time
}

func (parser *ParserHTTP) Parser() Response {

	var (
		ctype   string
		clength int
		StrHost string
		StrPort string
	)

	if len(parser.r.Header["Content-Type"]) >= 1 {
		ctype = GetContentType(parser.r.Header["Content-Type"][0])
	}

	if len(parser.r.Header["Content-Length"]) >= 1 {
		clength, _ = strconv.Atoi(parser.r.Header["Content-Length"][0])
	}

	SliceHost := strings.Split(parser.r.Request.URL.Host, ":")
	if len(SliceHost) > 1 {
		StrHost, StrPort = SliceHost[0], SliceHost[1]
	} else {
		StrHost = SliceHost[0]
		if parser.r.Request.URL.Scheme == "https" {
			StrPort = "443"
		} else {
			StrPort = "80"
		}
	}

	now := time.Now()

	r := Response{
		Origin:        parser.r.Request.RemoteAddr,
		Method:        parser.r.Request.Method,
		Status:        parser.r.StatusCode,
		ContentType:   string(ctype),
		ContentLength: uint(clength),
		Host:          StrHost,
		Port:          StrPort,
		URL:           parser.r.Request.URL.String(),
		Scheme:        parser.r.Request.URL.Scheme,
		Path:          parser.r.Request.URL.Path,
		Extension:     GetExtension(parser.r.Request.URL.Path),
		Header:        parser.r.Header,
		Body:          parser.respbody,
		RequestHeader: parser.r.Request.Header,
		RequestBody:   parser.reqbody,
		DateStart:     parser.s,
		DateEnd:       now,
	}
	return r
}

type Response struct {
	Origin        string      `json:"origin" db:",json"`
	Method        string      `json:"method" db:",json"`
	Status        int         `json:"status" db:",json"`
	ContentType   string      `json:"content_type" db:",json"`
	ContentLength uint        `json:"content_length" db:",json"`
	Host          string      `json:"host" db:",json"`
	Port          string      `json:"port" db:",json"`
	URL           string      `json:"url" db:",json"`
	Scheme        string      `json:"scheme" db:",json"`
	Path          string      `json:"path" db:",path"`
	Extension     string      `json:"ext" db:",path"`
	Header        http.Header `json:"header,omitempty" db:",json"`
	Body          []byte      `json:"body,omitempty" db:",json"`
	RequestHeader http.Header `json:"request_header,omitempty" db:",json"`
	RequestBody   []byte      `json:"request_body,omitempty" db:",json"`
	DateStart     time.Time   `json:"date_start" db:",json"`
	DateEnd       time.Time   `json:"date_end" db:",json"`
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

type ResType struct {
	ext   string
	ctype string
	mtype string
}

func NewResType(ext string, ctype string) *ResType {
	var mtype string
	if ctype != "" {
		mtype = strings.Split(ctype, "/")[0]
	}
	return &ResType{ext, ctype, mtype}
}

func (r *ResType) isStatic() bool {
	if ContainsString(static_ext, r.ext) {
		return true
	} else if ContainsString(static_types, r.ctype) {
		return true
	} else if ContainsString(media_types, r.mtype) {
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
