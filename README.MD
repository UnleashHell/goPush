## GOLANG APNS2 HTTP服务

基于apns2(http2)的http推送api服务,相比第一代apns传统socket方式,更加稳定,不会因为某个错误的device token而导致连接关闭,同时每条推送都能得到实时的结果反馈.

### 使用和配置

1. 配置conf/app.ini

| 参数            | 说明                            |
| ------------- | ----------------------------- |
| runMode       | 运行模式 dev or pro dev模式将日志输出在终端 |
| port          | http服务监听端口号                   |
| pemPath       | pem文件路径                       |
| pemPass       | pem文件密码                       |
| package       | app的包名                        |
| iosWorkerNum  | 启用多少携程用于推送消息                  |
| iosChannelNum | 存储推送消息channel缓存长度             |

2. 启动程序

```
go build goPush.go
nohup ./goPush >goPush.log 2>&1 &
```

3. 日志 logs目录将产生推送成功或失败日志 分别对应info-{$date}.log、error-{$date}.log

### 发送消息

```golang
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func main(){
	c := http.Client{
		Timeout: time.Duration(3) * time.Second,
	}
	requestUrl := "http://127.0.0.1:8080/api/push"
	data := url.Values{}
	data.Add("token", "设备token")
	data.Add("alert", "信息内容")
	data.Add("badge", "数字脚标")
	data.Add("sound", "声音")
	response, err := c.PostForm(requestUrl, data)
	if err != nil {
		fmt.Println(err);
	}else{
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err);
		}else{
			fmt.Println(string(body))
		}
	}	
}
```