
在 Windows 下有 QQ , 截图非常方便，　Ctrl+Alt+a, 就可以自由截图了。但是在 Liunx 下，截图就得使用一些软件了。我用的是 Linux Mint(Mate), 一个 Ubuntu 的衍生版， 附件里自带了一个截图工具。但是找不到快捷键启动它，每次都只能鼠标去点开，点开后还要选择一下自由截图（默认是截全屏），用起来非常不方便。更麻烦的是，这个软件的界面上没有提供任何帮助文档，以及快捷键设定。

今天特意去搜索了一下，发现这个软件叫做 mate-screenshot, 在命令行里面 mate-screenshot --help-all 就可以看到所有的帮助文档了。

在命令行里面执行 mate-screenshot -a ，这样鼠标就变成选取区域自由截图了。

在 Linux Mint 的首选项里，是有设定键盘快捷键的，添加一个快捷键，弹出来的对话框，填入名称和命令。保存后，在设置框里就可以看到这个要设置快捷键的命令了。然后如何设置快捷键启动这个命令呢？这里搞混了我，在界面上，点击 “禁用” 所在的位置，就可以按下你的快捷键，使之对应到这个命令上了。

然而，用这个设置方式，使 Ctrl+Alt+a 启动 mate-screenshot -a ,却是不行的，这是 Mate 的一个 Bug。还好皇天不负有心人，找到了解决方法。

+  在 /bin 目录下创建一个文件: screenshot
+  在 screenshot 写入如下内容：

    \#!/bin/bash
    sleep 0.1
    exec /usr/bin/mate-screenshot -a $@
+ 执行 sudo chmod a+x /bin/screenshot
+ 如上设置快捷键，不过命令变成： screenshot

现在，快捷键 Ctrl+Alt+a 就可以像 QQ 那样启动自由选取区域截图了




