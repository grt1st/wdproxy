package handles

import (
	"log"
	"net/http"

	"github.com/elazarl/goproxy"

	"github.com/grt1st/wdproxy/cache"
	"github.com/grt1st/wdproxy/g"
	"github.com/grt1st/wdproxy/utils"
)

var ResultsChan = make(chan *RequestRecord, 100)

func Init() {
	staticCache = cache.NewStorage("file")
	docCache = cache.NewStorage("doc")

	go func() {
		for {
			select {
			case rr := <-ResultsChan:
				if !g.Conf.NeedSave {
					continue
				}
				rr.Save()
			}
		}
	}()

}

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
	_, ok := g.C.Get(ctx.Session)
	if ok == false {
		//log.Println(err)
		return resp
	}
	respbody, err := ResponseBody(resp)
	if err != nil || respbody == nil {
		if err != nil {
			log.Println(err)
		}
		return resp
	}
	g.C.Del(ctx.Session)

	// Attaching capture tool.
	// reqbody
	RespCapture := NewRequestRecord(resp, nil, respbody)

	// Saving to MYSQL with a goroutine.
	go func() {
		//var static_resource = 0
		if utils.IsStatic(RespCapture.ContentType, RespCapture.Extension) {
			if g.Conf.RecordStatic {
				//static_resource = 1
				RespCapture.BodyResponse = []byte(nil)

				//fmt.Println(RespCapture.ContentLength, static_resource, RespCapture.Extension, RespCapture.Url, RespCapture.Status, RespCapture.Host, RespCapture.Port, RespCapture.BodyResponse, toJsonHeader(RespCapture.HeaderResponse), RespCapture.ContentType, RespCapture.Path, RespCapture.Scheme, RespCapture.Method, RespCapture.BodyRequest, toJsonHeader(RespCapture.HeaderRequest))
			}
		} //else {
		//	log.Println(RespCapture.ContentLength, static_resource, RespCapture.Extension, RespCapture.Url, RespCapture.Status, RespCapture.Host, RespCapture.Port, RespCapture.BodyResponse, toJsonHeader(RespCapture.HeaderResponse), RespCapture.ContentType, RespCapture.Path, RespCapture.Scheme, RespCapture.Method, RespCapture.BodyRequest, toJsonHeader(RespCapture.HeaderRequest))
		//}

		//RespCapture.Save()
		ResultsChan <- RespCapture
	}()

	return resp
}
