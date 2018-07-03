# wyproxy2
wyproxy: golang version
   
# 帮助
WYproxy Python Version   
https://github.com/ring04h/wyproxy
   
# 说明
学习Golang，基于goproxy库，顺便造了一个稳定的轮子。

## 使用帮助

```bash
yum install golang
cd /root/
git clone https://github.com/ring04h/wyproxy2
mv ./wyproxy2/ ./golang/
export GOPATH="/root/golang"
go build wyproxy.go
```
## 安装信任证书
iPhone使用safari浏览器打开   
https://raw.githubusercontent.com/ring04h/wyproxy2/master/src/goproxy/ca.pem   
   
## 启动代理服务器
```bash
cd /root/golang/
# 使用自定义的数据库配置信息
export WYDSN='root:@tcp(localhost:3306)/wyproxy?charset=utf8'
./wyproxy -addr 0.0.0.0:9999
```

## 使用代理功能
```
curl --proxy http://s5.wuyun.org:9999 www.ip.cn
```
