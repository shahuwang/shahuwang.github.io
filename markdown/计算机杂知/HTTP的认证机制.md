读 Go 的 http 库代码，看到了 BasicAuth 部分，之前对 http 的基本认证都没什么了解，也没用过，所以就好奇，找来《HTTP权威指南》来看，发现这 http 的认证，还挺多知识点的，值得写篇文章记录下。

先来说下 http 的基本认证，这是一种相当简单的认证方式，把用户名密码 base64 一下，放到 header 里发送给服务器端即可。
因为这种方式，基本算是明文传输密码，所以是相当不安全。不安全归不安全，用得还是挺多的。

基本认证的交流过程如下：

client 发送一个请求到 server， server 检查其头部的 Authorization, 如果没有值或者解析出的用户名和密码错误，则返回
401 Unauthorized, 同时在 response 的 header 里，添加头部 WWW-Authenticate, 参考 [MDN WWW-Authenticate](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/WWW-Authenticate),
其格式为 `WWW-Authenticate: <type> realm=<realm>` , 如 `WWW-Authenticate: Basic`, `WWW-Authenticate: Basic realm="Access to the staging site"` 。 

这里面 realm 是指安全域，填入你想填的内容，让客户端知道它该填什么用户名和密码。

client 收到 401 回复，明白要添加什么用户名和密码，则将用户名和密码用 `:` 拼接，然后再 base64 编码，放入头部的
Authorization 里。即 `request.Header.put('Authorization', base64encode(username+":"+password))`， 再次发起请求。

server 收到请求，解析出用户名和密码，授权通过，返回 client 所需要的资源。

用 Go 写一个支持基本认证的 server， 代码参考[此文](http://www.dotcoo.com/golang-http-auth), 如下：


```
package main

import (
    "fmt"
    "io"
    "net/http"
    "log"
    "encoding/base64"
    "strings"
)

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
    auth := req.Header.Get("Authorization")
    if auth == "" {
        w.Header().Set("WWW-Authenticate", `Basic realm="Dotcoo User Login"`)
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    fmt.Println(auth)

    auths := strings.SplitN(auth, " ", 2)
    if len(auths) != 2 {
        fmt.Println("error")
        return
    }

    authMethod := auths[0]
    authB64 := auths[1]

    switch authMethod {
    case "Basic":
        authstr, err := base64.StdEncoding.DecodeString(authB64)
        if err != nil {
            fmt.Println(err)
            io.WriteString(w, "Unauthorized!\n")
            return
        }
        fmt.Println(string(authstr))

        userPwd := strings.SplitN(string(authstr), ":", 2)
        if len(userPwd) != 2 {
            fmt.Println("error")
            return
        }

        username := userPwd[0]
        password := userPwd[1]

        fmt.Println("Username:", username)
        fmt.Println("Password:", password)
        fmt.Println()

    default:
        fmt.Println("error")
        return
    }


    io.WriteString(w, "hello, world!\n")
}

func main() {
    http.HandleFunc("/hello", HelloServer)
    err := http.ListenAndServe(":12345", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
```

有时候为了开发方便，不好在所有的服务都添加基本认证，那么就可以添加一层代理服务器，由代理服务器进行基本认证。

代理服务器的认证，称之为代理认证，与上面的认证过程和方法都是一样的，唯一不同的是，代理认证返回的是 407 Unauthorized, 使用的 header 也变成了 Proxy-Authenticate , Proxy-Authorization 和 Proxy-Authentication-Info。

基本认证由于是明文传输，所以要保证安全，最好与 HTTPS 结合使用。

另一种更加安全点的认证方式，就是摘要验证。其基本思想就是不传输密码，而是传输密码的摘要，比如对密码进行md5
取摘要。服务器端根据用户名，找到用户密码，按照相当的方式对密码进行摘要，然后判断此摘要和客户端传输上来的摘要是否一致，一致则通过。

为了防止重放攻击（拿到你的用户名和密码摘要，其实和拿到用户名和密码是一样的效果），实际使用中，需要将密码和随机字符串拼接后再进行摘要计算。服务器端返回给用户一个随机字符串，客户端将之与密码拼接计算摘要，服务器端亦如是，然后比对摘要，判断是否通过验证。

摘要验证在《HTTP权威指南》里有比较详细的叙述，这里有篇文章，算是把书中的内容抄上来了：[前端学HTTP之摘要认证](http://www.cnblogs.com/xiaohuochai/p/6189065.html)
