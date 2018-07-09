package handles

import (
	"net/http"
	"log"
	"github.com/elazarl/goproxy"
	"github.com/grt1st/wdproxy/g"
	"fmt"
)

var ResultsChan = make(chan *RequestRecord)

func init() {

	go func() {
		for{
			select {
			case rr := <- ResultsChan:
				fmt.Println(rr.String())
			}
		}
	}()

}

var (

	record_static = true

	// http static resource file extension
	static_ext = []string{
		"js",
		"css",
		"ico",
	}

	// media resource files type
	media_types = []string{
		"image",
		"video",
		"audio",
	}

	// http static resource files
	static_types = []string{
		"text/css",
		// "application/javascript",
		// "application/x-javascript",
		"application/msword",
		"application/vnd.ms-excel",
		"application/vnd.ms-powerpoint",
		"application/x-ms-wmd",
		"application/x-shockwave-flash",
	}
)

func HandleRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	reqbody, err := RequestBody(req)
	if err != nil {
		log.Println(err)
		return req, nil
	}
	g.C.Add(ctx.Session, reqbody)
	return req, nil
}


func HandleResponse(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	reqbody, ok := g.C.Get(ctx.Session)
	if ok == false {
		//log.Println(err)
		return resp
	}
	respbody, err := ResponseBody(resp)
	if err != nil {
		log.Println(err)
		return resp
	}
	g.C.Del(ctx.Session)

	// Attaching capture tool.
	RespCapture := NewRequestRecord(resp, reqbody, respbody)

	// Saving to MYSQL with a goroutine.
	go func() {
		var static_resource = 0
		if RespCapture.isStatic() {
			if record_static {
				static_resource = 1
				RespCapture.BodyResponse = []byte(nil)

				fmt.Println(RespCapture.ContentLength, static_resource, RespCapture.Extension, RespCapture.Url, RespCapture.Status, RespCapture.Host, RespCapture.Port, RespCapture.BodyResponse, toJsonHeader(RespCapture.HeaderResponse), RespCapture.ContentType, RespCapture.Path, RespCapture.Scheme, RespCapture.Method, RespCapture.BodyRequest, toJsonHeader(RespCapture.HeaderRequest))
			}
		} else {
			log.Println(RespCapture.ContentLength, static_resource, RespCapture.Extension, RespCapture.Url, RespCapture.Status, RespCapture.Host, RespCapture.Port, RespCapture.BodyResponse, toJsonHeader(RespCapture.HeaderResponse), RespCapture.ContentType, RespCapture.Path, RespCapture.Scheme, RespCapture.Method, RespCapture.BodyRequest, toJsonHeader(RespCapture.HeaderRequest))
		}
	}()

	return resp
}