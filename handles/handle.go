package handles

import (
	"net/http"
	"log"
	"github.com/elazarl/goproxy"
	"github.com/grt1st/wdproxy/g"
	"fmt"
)

var (

	record_static = true

	// http static resource file extension
	static_ext []string = []string{
		"js",
		"css",
		"ico",
	}

	// media resource files type
	media_types []string = []string{
		"image",
		"video",
		"audio",
	}

	// http static resource files
	static_types []string = []string{
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
	RespCapture := NewParserHTTP(resp, reqbody, respbody).Parser()

	// Saving to MYSQL with a goroutine.
	go func() {
		var static_resource int = 0
		static := NewResType(
			RespCapture.Extension,
			RespCapture.ContentType).isStatic()
		if static {
			if record_static {
				static_resource = 1
				RespCapture.Body = []byte(nil)

				fmt.Println(RespCapture.ContentLength, static_resource, RespCapture.Extension, RespCapture.URL, RespCapture.Status, RespCapture.Host, RespCapture.Port, RespCapture.Body, toJsonHeader(RespCapture.Header), RespCapture.ContentType, RespCapture.Path, RespCapture.Scheme, RespCapture.Method, RespCapture.RequestBody, toJsonHeader(RespCapture.RequestHeader), RespCapture.DateStart, RespCapture.DateEnd)
			}
		} else {
			log.Println(RespCapture.ContentLength, static_resource, RespCapture.Extension, RespCapture.URL, RespCapture.Status, RespCapture.Host, RespCapture.Port, RespCapture.Body, toJsonHeader(RespCapture.Header), RespCapture.ContentType, RespCapture.Path, RespCapture.Scheme, RespCapture.Method, RespCapture.RequestBody, toJsonHeader(RespCapture.RequestHeader), RespCapture.DateStart, RespCapture.DateEnd)
		}
	}()

	return resp
}