# http 转 socket5 的方式（shadowsocket）

有一些穿越长城的工具如是只提供socket5的上网方式的，但是很多的访问时需要使用到http和https的协议的。总不能买好几个不同类型的账号来穿越吧。

ubuntu 下有一款工具挺好用的，Privoxy， sudo apt-get install privoxy 即可安装。

然后打开文件 /etc/privoxy/config， 在文件的最后面添加这么一句： `forward-socks5 / 127.0.0.1:1080 .` 

一定要看仔细了， 是 forward-socks5, 不是 socket5 喔， 特别是最后是有一个圆点句号的. 127.0.0.1:1080 是你的socket5 穿越工具的代理地址。

然后执行命令 sudo service privoxy restart 就可以用了。

然后现在所有的http请求都可以指向 127.0.0.1:8118了

我用的是 linux mint， 在控制中心找到网络代理，设置手动配置，将http代理设置为 127.0.0.1， 端口8118就行了.





