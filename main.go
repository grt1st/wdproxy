package main

import (
	"runtime"
	"github.com/elazarl/goproxy"
	"net/http"
	"log"
	"github.com/grt1st/wdproxy/handles"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	proxy := goproxy.NewProxyHttpServer()

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(handles.HandleRequest)
	proxy.OnResponse().DoFunc(handles.HandleResponse)

	proxy.Verbose = true

	err := http.ListenAndServe("127.0.0.1:1080", proxy)
	if err != nil {
		log.Fatalln(err)
	}
}
