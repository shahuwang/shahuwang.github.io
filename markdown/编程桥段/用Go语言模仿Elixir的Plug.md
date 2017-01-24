现在公司的项目，后端主要用两门语言，Python和Elixir，而且因为Elixir的一些基础库写得越来越完善了，所以现在大家都更偏向于用Elixir写新的项目或功能。目前所有的API接口，都是使用Plug来实现的，刚开始用的时候，对Elixir和Plug都不熟悉，网上文档非常少，当时觉得这破东西好难用。用多熟悉了之后，发现条理上还是很清晰的。

Plug大致就这么三个部分，底层使用cowboy来处理http相关的部分，提供Router进行路由匹配，以及各种plug函数，对整个请求流程做统一的处理。精华当然就是这个plug了，按照规定的格式，写一个你自己的plug模块，方便又好用。


```
defmodule Example.Plug.Router do
  use Plug.Router

  plug :match
  plug Plug.Parsers, parsers: [Plug.Parsers.Jiffy]
  plug :dispatch

  get "/", do: send_resp(conn, 200, "Welcome")
  match _, do: send_resp(conn, 404, "Oops!")
end
```

如上代码，就写好了一个简单的 API 应用，plug :match 是找到匹配的路径，plug Plug.Parsers 会将请求参数解析为json，plug dispatch 表示去执行匹配到的处理函数了。

plug的实现也很简单，如下 plug， 直接返回 hello world （后面的就不执行了）：


```
defmodule MyPlug do
  import Plug.Conn

  def init(options) do
    # initialize options

    options
  end

  def call(conn, _opts) do
    conn
    |> put_resp_content_type("text/plain")
    |> send_resp(200, "Hello world")
  end
end
```

公司的项目自己实现了相当多的plug，典型的如权限验证，参数解析，请求日志打印等等。

用多了之后，就萌生了用Go语言来实现一个Plug，于是乎就有了这个 [glug](http://note.youdao.com/) 。磨蹭了将近一个月，越写越没用动力，还有很多很多的功能没有写，目前就实现了最基本的路由功能，以及 glug 函数（参照plug）。

Go的http库非常的完善非常好用，不像erlang的cowboy，非常不好用，所以elixir写了好大一堆代码来封装cowboy。因此，其实实现核心功能，需要的代码并不是很多。

路由部分，我使用 Trie 树来进行路径查找，Get方法一棵树，Post方法一棵树。每个节点都包含一种segment，segment就是路径 /a/b/c 中的 a 呀，b 呀， c 呀。目前仅有两种 segment，NormalSegment和 VariantSegment，前者就是简单的字符串路径，后者就是 /a/:name/ 这个 :name 部分，用于捕获这个部分的值，值存放在 conn.PathParams 里面。后面如果需要，可以添加正则等之类的segment。

glug函数需要符合 `type GlugFunc func(*Connection) bool` , 比如目前实行的 logger 如下：


```
package glug

import (
	"log"
	"net/http"
	"time"
)

func Logger(conn *Connection) bool {
	start := time.Now()
	addr := conn.Request.Header.Get("X-Real-IP")
	if addr == "" {
		addr = conn.Request.Header.Get("X-Forwarded-For")
		if addr == "" {
			addr = conn.Request.RemoteAddr
		}
	}
	log.Printf("Started %s %s for %s", conn.Request.Method, conn.Request.URL.Path, addr)
	fun := func(resp *Resp) {
		status := resp.Status
		statusText := http.StatusText(int(status))
		duration := time.Since(start)
		log.Printf("Completed %v %s in %v\n", status, statusText, duration)
	}
	conn.Register(fun)
	return true
}
```

返回值为true，就会执行下去，否则，直接终止。

简单的使用例子如下：

```
package main
import(
    "github.com/shahuwang/glug"
    "net/http"
)

func main() {
    router := glug.NewRouter()
    router.Use(router.Match)
    router.Use(glug.Logger)
    router.Use(router.Dispatch)
    router.Get("/login", func(conn *glug.Connection){
        conn.Sendresp(200, conn.Request.Header, []byte("hello world"))
    })
    
    http.ListenAndServe(":8080", router)
}
```

大致就是这样的一个简单的东西，要做到能真正使用在生产环境中，还需要很多细节完善。

每次构思写个东西的时候，都会臆想写出个多么牛逼的东西，功能多么强大，然后很多人来用，迎娶白富美，走上人生巅峰。但是写着写着，拖延症就开始发作了，接着就是懒癌发作，然后就草草了事了，并没有动力去做出一个真正能用的东西出来。

无怪乎这些年我一直这么失败。
