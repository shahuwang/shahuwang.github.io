�������� Elixir �� Plug �������һ�� API �ӿڣ����ڴ�����������Ϊ���֣�Plug�����ϴ����ļ��õ�������ṹ[https://hexdocs.pm/plug/Plug.Upload.html](http://note.youdao.com/)��
��ǰһֱʹ�� HTML ���� js �ı��������ļ��ϴ�������ǰ���Եö���һ�� ��HTTPȨ��ָ�ϡ���Ҳ���������ò���ˡ�
���Ե���Ҫд�ĵ�����Python��˵����ô��ʹ������ӿڵ�ʱ�򣬲ŷ����Ҷ�httpЭ���е�multipart/form-data�������˽⣬����requests���ĵ������ƣ���Ȼ�����ҿ������֪����ô������

�ҵĽӿ���Ҫ�û��ṩ�����ֶΣ�rate�������ʣ��� audio(�����ļ�)��Plug���ϴ����ļ���װΪ %Plug.Upload{} �ṹ���˽ṹ�����������ݣ�path��ʾ�ļ�·����filename��ʾ�ļ�����content_type��ʾ�ļ������ͣ��ҵĽӿ�Ҫ�����content_typeΪaudio/amr��

��֪������������Ƿ�����ɻ��ˣ���������£��ѵ�����һ��http����ֻ��һ��content-type�� Ϊʲô�ϴ����ļ����ᵥ����һ��content-type��

��������ɻ�ĵط���ȥ����multipart/form-data�Ľ��ͣ��������Щ��������������ʾ��һ��http�����content-type����Ϊmultipart/form-data, ��ʾ���http���󣬻����˶�����֣�ÿ�����ֶ��ǲ�һ�����������ͣ�ÿ�����ֶ����������Լ���content-type�Լ�����һЩ��ص�header�������ԣ���ÿ��part����������http����ͺã���

����html����˵��ÿ��field�������е�һ��part��ÿһ��part������������header�� `Content-Disposition: form-data; name="field1"`, Content-Disposition����form-data, name��������ֶε����ơ���ͨ���ֶΣ���content-type����Ĭ�ϵ�text/plain��

һ��multipart/form-data���󣬻��Զ�����һ���ַ�����boundary����ÿ��part�ķָ��־��

��������Go����дһ���򵥵�http����������http���������ַ������ػ�ȥ��Ȼ��ʹ��chrome��һ�����postman����form�����ύ�����������ػ�����http��������ô���ġ�

Go�Ĵ������£�

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

postman��������ͼ, ��������������ͨ�ֶΣ�һ����ͨ�ı��ļ���һ��ͼƬ�ϴ�

![image](https://raw.githubusercontent.com/shahuwang/images/master/%E8%AE%A1%E7%AE%97%E6%9C%BA%E6%9D%82%E7%9F%A5/post.png)

���ػ�����http���Ľṹ���£�


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
Content-Disposition: form-data; name="imagefield"; filename="����Ŀ.jpg"
Content-Type: image/jpeg

????
(ע��ʣ�µĶ���ͼƬ������)
------WebKitFormBoundaryMfzICyp6us00g3IO--
```

���Կ�������http��������岿�֣�content-type��`Content-Type: multipart/form-data; boundary=----WebKitFormBoundaryMfzICyp6us00g3IO`���Զ���������boundary�ַ�����

Ȼ��ÿ���ֶζ�ʹ��boundary��Χ�ţ������header�����ݸ��Ե����ݲ�ͬ��������ͬ����ͼƬΪ������content-typeΪimage/jpeg,Ҳ����filename������ļ��������ǷǱ�׼�ģ�Ĭ�ϵ�content-type��application/octet-stream ��֮ǰ����һ���ļ����ص�api����Ϊ���ص��ļ��п������ı����п�����ͼƬ�����п������������޷���֪�����ʽ��������ֱ�Ӱ�response��content-type����Ϊapplication/octet-stream,�����ļ��������ء�

��֪��Ϊɶ���������Ե�http���post������Ҫ��multipart/form-data�ϴ��ļ������񶼲��Ǻܷ��㣬�������ĵ����Ǻ���ϸ��

��ô�ص��ʼ��������������python��requests�����ʵ���������postman���͵�Ч���أ�


```
import requests
url = "http://127.0.0.1:8989/hello"
data = {"field1": "I am first", "field2": "I am second"}
file1 = ("raceon.go", open("raceon.go", "rb"), "application/octet-stream", {"Content-Type": "application/octet-stream"})
file2 = ("����Ŀ.jpg", open("����Ŀ.jpg", "rb"), "image/jpeg", {"Content-Type": "image/jpeg"})
resp = requests.post(url, data=data, files={"field3": file1, "field4": file2}, headers={})
print resp.content[0:500]
```

����Ĵ������ʵ��postmanͬ����Ч���ˣ�ע����ǣ�ÿ��file���Ҷ�����������content-type, ��֪��Ϊɶ�ڹ�˾�ĵ����ϣ�ֻ��{"Content-Type": "image/jpeg"} �������ò�����Ч�ģ��������Լ��ĵ����ϣ�����ֻ�е�һ�� "image/jpeg" ������Ч�ģ�Ӧ���Ǹ�bug�ɡ�


