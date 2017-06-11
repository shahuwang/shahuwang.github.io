学艺不精，用http，特别是写前端代码的时候，总觉得表单这一块是个谜，搞不懂。之前写过一篇文件[《http multipart处理文件上传的一点理解》](http://shahuwang.com/%E8%AE%A1%E7%AE%97%E6%9C%BA%E6%9D%82%E7%9F%A5/http%20multipart%E5%A4%84%E7%90%86%E6%96%87%E4%BB%B6%E4%B8%8A%E4%BC%A0%E7%9A%84%E4%B8%80%E7%82%B9%E7%90%86%E8%A7%A3.html)，这篇讲的是http请求如何将multipart/form封装起来发给服务器的，本文要写的则是服务器如何解析Form和multipart/form的。

首先来看看普通的表单，也就是`Content-Type: application/x-www-form-urlencoded`这种类型的表单，其http报文结构是怎么样的，依然采用[《http multipart处理文件上传的一点理解》](http://shahuwang.com/%E8%AE%A1%E7%AE%97%E6%9C%BA%E6%9D%82%E7%9F%A5/http%20multipart%E5%A4%84%E7%90%86%E6%96%87%E4%BB%B6%E4%B8%8A%E4%BC%A0%E7%9A%84%E4%B8%80%E7%82%B9%E7%90%86%E8%A7%A3.html)里用到的那段Go代码。


```
POST /hello HTTP/1.1
Host: 127.0.0.1:8989
Accept: */*
Accept-Encoding: gzip, deflate, br
Accept-Language: zh-CN,zh;q=0.8,en;q=0.6,zh-TW;q=0.4,la;q=0.2
Cache-Control: no-cache
Connection: keep-alive
Content-Length: 36
Content-Type: application/x-www-form-urlencoded
Origin: chrome-extension://aicmkgpgakddgnaphhhpliifpcfhicfo
Postman-Token: b98e76e9-3552-b53d-c9a6-ae7d425e99ef
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36

first=I+am+first&second+=I+am+second
```
从报文结构可以看出，表单部分的内容被进行了URL encode 放到 body 部分里面去了。

所以对于服务器端来说，解析`Content-Type: application/x-www-form-urlencoded`就是把body部分URL decode成一个map即可。

之前我一直疑惑multipart/form里的文件，到了服务器端之后，是怎么个处理过程。仔细读了Go的http库代码，才了解到它是这么一个处理过程：

1. 确定当前请求为multipart/form
2. 逐字段解析，如果只是普通值，则存到map里
3. 如果是文件类型，则读取该字段内容，如果内容长度小于某个值，直接将内容封装到一个FileHeader结构里，等于就是把文件放到内存里了。如果内容长度超过某个值，则将内容写入到tmp文件里面去，同时记录路径，以备后面使用者自己去读取。