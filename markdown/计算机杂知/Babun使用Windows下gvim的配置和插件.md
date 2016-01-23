# Babun使用Windows下gvim的配置和插件

windows自带的CMD屎一样，但是平时在家的时候又不想使用Linux系统。开着虚拟机，总是有些问题导致使用体验不是很爽。所以很是希望windows下也能有Linux下那些常用的工具，这样就不用那么多折腾了。

最近发现 Babun 正是我需要的，集成了 zsh和cygwin，以及其他大量的Linux常用工具。不过，由于cygwin是不带GUI的，在babun里面启动 gvim，打开的是windows下安装的gvim。但是又由于cygwin的一些环境变量原因，打开的gvim没能使用windows下的配置和安装好的插件。总不能babun里面又配制一份吧，浪费硬盘也不方便。

昨晚折腾了一晚上，终于把这个问题解决了。

*注: 我的用户名是rickey，babun安装的时候默认路径是C:/Users/rickey/.babun, babun的默认用户名也是rickey*

+ 第一步，在babun安装目录下的cygwin/home/rickey，删除掉 .vim 文件夹
+ 第二步，在开始菜单找到cmd，右键选择管理员权限打开
+ 第三步，在cmd里用mklink将C:/Users/rickey/.vim 链接到 babun的home目录，即 cygwin/home/rickey，命令如下：mklink /D C:\Users\rickey\.babun\cygwin\home\rickey\.vim  C:\Users\rickey\.vim
+ 第四步：打开babun，将 .vimrc 链接到home目录，ln  C:/Users/rickey/.vimrc /home/rickey/.vimrc ，注意软链接是不行的
+ 第五步：在babun里，打开 nano ~/.zshrc, 在末尾填上：export VIM='C:/Program Files/Vim/vim74'

现在打开的gvim，和在windows下打开gvim就是一模一样的了。不过babun自带的vim此时会变得不好用了，只是我用gvim比较多，就不管vim了。

不过，光是这样配置还是有问题的，因为在babun启动的时候使用的是babun的路径模式，但是在使用gvim的时候，使用的是windows的路径，一保存，就会出现保存不了，因为路径找不到。

所以首先，在/home/rickey/bin 目录下创建一个gvim.sh, 内容如下：

```
#!/usr/bin/env bash
gvim `cygpath -W $a` &
```

保存完毕后，对这个文件执行命令 chmod a+x gvim.sh, 然后，再打开 ~/.zshrc 文件，在最后面添加上 alias gvim=gvim.sh,大功告成。
另外，还有一个要注意的点，就是在cygwin下使用git， git st的时候会出现乱码，使用命令：git config --global core.quotepath false 解决问题。



