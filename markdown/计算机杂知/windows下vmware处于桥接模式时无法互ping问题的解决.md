由于实在无法离开谷歌（虽然 [AOL search](http://search.aol.com/aol/search?enabled_terms=&s_it=comsearch&q=china&s_chn=prt_aol20) 也是调用谷歌的API，效果等同，但是很多网站还是无法打开)， 所以就买了云梯的vpn。虽然它提供了所谓的智能模式，也就是自动识别国内国外网站，过滤掉国内网站不通过vpn，但是实在做得不好，导致微博打开的速度受到非常大的影响。

作为一个程序员，自然是不能被这样的问题影响的。自己写一个过滤器这样的事情我是做不来，所以我用vmware开了一个Ubuntu虚拟机，在虚拟机里面连接vpn。然后Ubuntu开启squid2的代理转发服务。用chrome的代理插件，将代理地址指向Ubuntu。这样的话，就不会影响我上一些不需要翻墙网站的速度了。

但是问题也来了，之前我的vmware都是使用NAT连接模式的，这种模式挺方便的，只是vpn经常断掉，最近还干脆连不上了。后来试了桥接模式，发现这种模式下vpn极其稳定。本来很兴奋的，却又发现，windows和Ubuntu互相ping不通了。搜索了好久，说是防火墙的问题，把防火墙一关，果然就行了。但防火墙不能不开吧，我把所有的规则都设置为通过了，但只要开着防火墙，windows和ubuntu就互相不能ping通。

找了好久也没有找到解决的办法，后来偶然看到说在桥接模式下，要网段设为一致才可以，且桥接模式下宿主机是处于 192.168.0.xx 的网段的。所以我在控制面板里查看我的网络连接，如下图：

![网络连接](https://raw.githubusercontent.com/shahuwang/images/master/%E8%AE%A1%E7%AE%97%E6%9C%BA%E6%9D%82%E7%9F%A5/1.png)

我目前连着的是WLAN，VMnet1 和 VMnet8 是 VMware 设置出来的。然后，我再点击 WLAN 里面看，点击属性--》共享，如下图：
![连接状态](https://raw.githubusercontent.com/shahuwang/images/master/%E8%AE%A1%E7%AE%97%E6%9C%BA%E6%9D%82%E7%9F%A5/2.png)

不太明白意思，貌似是不是 VMnet8 是通过 WLAN 这个连接来进行？然后回到网络连接页面，点击 VMnet8，点击属性，如下图：

![vmnet8](https://raw.githubusercontent.com/shahuwang/images/master/%E8%AE%A1%E7%AE%97%E6%9C%BA%E6%9D%82%E7%9F%A5/3.png)

把 VMware Bridge Protocol 这个勾选上，然后拉到下面，找到 Internet 协议版本 4（TCP/IPv4) ，选择它（不要去掉勾选），然后点击属性，设置如下图即可：

![ip设置](https://raw.githubusercontent.com/shahuwang/images/master/%E8%AE%A1%E7%AE%97%E6%9C%BA%E6%9D%82%E7%9F%A5/4.png)

另外，这里的 IP 地址一栏， 你可以填 192.168.0. xx， 不一定是 192.168.0.6 。

现在在 Windows 下可以 ping通虚拟机的ubuntu了，但是虚拟机的Ubuntu却无法ping通windows。现在需要修改 Ubuntu 的 ip 设置：
sudo vim /etc/network/interfaces

默认只有两行，然后将内容修改如下：

    auto lo
    iface lo inet static
    address 192.168.0.6
    gateway 192.168.0.1
    netmask 255.255.255.0

再关闭掉防火墙 sudo ufw disable

重启虚拟机，应该就可以ping通windows了。


