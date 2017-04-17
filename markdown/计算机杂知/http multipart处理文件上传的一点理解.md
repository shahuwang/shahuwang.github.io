工作中用 Elixir 的 Plug 框架做了一个 API 接口，用于处理语音翻译为文字，Plug处理上传的文件用到了这个结构[https://hexdocs.pm/plug/Plug.Upload.html](http://note.youdao.com/)，
以前一直使用 HTML 或者 js 的表单功能做文件上传，几年前粗略得读了一遍 《HTTP权威指南》，也基本上忘得差不多了。
所以当我要写文档，用Python来说明怎么样使用这个接口的时候，才发现我对http协议中的multipart/form-data基本不了解，加上requests库文档不完善，居然花了我快两天才知道怎么样做。

我的接口需要用户提供两个字段，rate（采样率）， audio(语音文件)，Plug对上传的文件封装为 %Plug.Upload{} 结构，此结构包含三个数据，path表示文件路径，filename表示文件名，content_type表示文件的类型，我的接口要求这个content_type为audio/amr。

不知道读到这里，你是否产生疑惑了？正常情况下，难道不是一个http请求，只有一个content-type吗， 为什么上传的文件还会单独有一个content-type？

这就是我疑惑的地方，去看了multipart/form-data的解释，才理解了些。正如其名字所示，一个http请求的content-type设置为multipart/form-data, 表示这个http请求，混杂了多个部分，每个部分都是不一样的数据类型，每个部分都可以设置自己的content-type以及其他一些相关的header参数（对，把每个part看出独立的http请求就好）。

对于html表单来说，每个field就是其中的一个part，每一个part，都包含两个header： `Content-Disposition: form-data; name="field1"`, Content-Disposition都是form-data, name就是这个字段的名称。普通的字段，其content-type就是默认的text/plain。

一个multipart/form-data请求，会自动生成一串字符串做boundary，即每个part的分割标志。

我这里用Go语言写一个简单的http服务器，把http请求报文以字符串返回回去。然后使用chrome的一个插件postman来做form表单的提交，并看看返回回来的http报文是怎么样的。

Go的代码如下：

```
package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func Hello(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("%+v\n", req.Header)
	dump, _ := httputil.DumpRequest(req, true)
	w.Write(dump)
}

func main() {
	http.HandleFunc("/hello", Hello)
	http.ListenAndServe(":8989", nil)
}
```

postman的设置如图, 我设置了两个普通字段，一个普通文本文件，一个图片上传

![image](https://raw.githubusercontent.com/shahuwang/images/master/%E8%AE%A1%E7%AE%97%E6%9C%BA%E6%9D%82%E7%9F%A5/post.png)

返回回来的http报文结构如下：


```
POST /hello HTTP/1.1
Host: 127.0.0.1:8989
Accept: */*
Accept-Encoding: gzip, deflate, br
Accept-Language: zh-CN,zh;q=0.8,en;q=0.6,zh-TW;q=0.4,la;q=0.2
Cache-Control: no-cache
Connection: keep-alive
Content-Type: multipart/form-data; boundary=----WebKitFormBoundaryMfzICyp6us00g3IO
Origin: chrome-extension://aicmkgpgakddgnaphhhpliifpcfhicfo
Postman-Token: e6dd5c45-ed9e-8cdb-1fc6-efdb4abe3caa
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36

------WebKitFormBoundaryMfzICyp6us00g3IO
Content-Disposition: form-data; name="field1"

I am first-----------------------------
------WebKitFormBoundaryMfzICyp6us00g3IO
Content-Disposition: form-data; name="field2"

I am second====================
------WebKitFormBoundaryMfzICyp6us00g3IO
Content-Disposition: form-data; name="txtfield"; filename="test.txt"
Content-Type: text/plain

txt file post test
------WebKitFormBoundaryMfzICyp6us00g3IO
Content-Disposition: form-data; name="imagefield"; filename="初代目.jpg"
Content-Type: image/jpeg

????
(注：剩下的都是图片的乱码)
------WebKitFormBoundaryMfzICyp6us00g3IO--
```

可以看到，此http请求的主体部分，content-type是`Content-Type: multipart/form-data; boundary=----WebKitFormBoundaryMfzICyp6us00g3IO`，自动附带上了boundary字符串。

然后每个字段都使用boundary包围着，里面的header都根据各自的内容不同而有所不同。以图片为例，其content-type为image/jpeg,也有了filename。如果文件的类型是非标准的，默认的content-type是application/octet-stream 。之前做过一个文件下载的api，因为下载的文件有可能是文本，有可能是图片，还有可能是其他，无法得知具体格式，所以我直接把response的content-type设置为application/octet-stream,所有文件都能下载。

不知道为啥，各门语言的http库的post方法，要做multipart/form-data上传文件，好像都不是很方便，甚至是文档不是很详细。

那么回到最开始的问题上来，用python的requests库如何实现上面这个postman发送的效果呢？


```
import requests
url = "http://127.0.0.1:8989/hello"
data = {"field1": "I am first", "field2": "I am second"}
file1 = ("raceon.go", open("raceon.go", "rb"), "application/octet-stream", {"Content-Type": "application/octet-stream"})
file2 = ("初代目.jpg", open("初代目.jpg", "rb"), "image/jpeg", {"Content-Type": "image/jpeg"})
resp = requests.post(url, data=data, files={"field3": file1, "field4": file2}, headers={})
print resp.content[0:500]
```

上面的代码就能实现postman同样的效果了，注意的是，每个file，我都设置了两次content-type, 不知道为啥在公司的电脑上，只有{"Content-Type": "image/jpeg"} 这样设置才是有效的，而在我自己的电脑上，则是只有第一个 "image/jpeg" 才是有效的，应该是个bug吧。


