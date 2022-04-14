package main

import (
	"crypto"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"

	"github.com/elazarl/goproxy"
	"golang.org/x/crypto/pkcs12"

	"github.com/grt1st/wdproxy/g"
	"github.com/grt1st/wdproxy/handles"
)

const (
	PKCSPath     = "/Users/kintenroku/Desktop/charles-ssl-proxying.p12" // Charles: Help>SSL Proxying>Export Charles Root Certificate and Private Key...
	PKCSPassword = "123"                                                // the password to protect the exported certificate and key
)

func main() {
	// 初始化
	g.Init()
	handles.Init()

	runtime.GOMAXPROCS(runtime.NumCPU())
	SetCAPKCS() // 设置证书
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().HandleConnect(CustomMITM)

	// 添加方法
	//proxy.OnRequest().DoFunc(handles.SslStripRequest)
	//proxy.OnResponse().DoFunc(handles.SslStripResponse)
	//proxy.OnResponse().DoFunc(handles.HandleStaticResponse)
	//proxy.OnResponse().DoFunc(handles.HandleDocResponse)
	defer close(handles.ResultsChan)

	proxy.Verbose = false //记录日志

	// 开始监听
	l := "0.0.0.0:8880"
	log.Printf("[+] WdProxy start listening at %s", l)
	err := http.ListenAndServe(l, proxy)
	if err != nil {
		log.Fatalln(err)
	}
}

var CustomMITM goproxy.FuncHttpsHandler = func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	if strings.Contains(host, "example.cn") {
		// 放行不处理，白名单
		return goproxy.OkConnect, host
	}
	// 需要处理的
	return goproxy.MitmConnect, host
}

func SetCAPKCS() {
	p12Data, err := ioutil.ReadFile(PKCSPath)
	if err != nil {
		log.Fatal(err)
	}

	key, cert, err := pkcs12.Decode(p12Data, PKCSPassword)
	if err != nil {
		log.Fatal(err)
	}

	goProxyCA := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  key.(crypto.PrivateKey),
		Leaf:        cert,
	}

	goproxy.GoproxyCa = goProxyCA
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goProxyCA)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&goProxyCA)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&goProxyCA)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&goProxyCA)}
}
