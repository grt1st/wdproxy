package main

import (
	"runtime"
	"github.com/grt1st/wdproxy/handles"
	"net/http"
	"log"
	"fmt"
	"github.com/elazarl/goproxy"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	proxy := goproxy.NewProxyHttpServer()

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(handles.HandleRequest)
	proxy.OnResponse().DoFunc(handles.HandleResponse)
	defer close(handles.ResultsChan)
	proxy.Verbose = false //true

	err := http.ListenAndServe("127.0.0.1:1080", proxy)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("[+] WdProxy start listening at http://127.0.0.1:1080")
	//platform.RunWebPlatform()
}
