看大部分的C/C++代码，自然第一步是从Makefile整起。不过该死的是，Makefile这本来就相当奇葩的东西，网上教程居然非常的少。前段时间花了一个星期去看Redis的源代码，Redis的Makefile比较复杂，也弄了比较多的技巧，当时为了把它的Makefile弄明白，真的是花了很多的时间。leveldb的makefile则简单多了，所以本系列的第一篇文章，Makefile必知必会，就从Redis的Makefile说起。

make这个工具，其实有很多的实现版本，各有各的优缺点，网上我好像就找到了[GNU Make](http://www.gnu.org/software/make/manual/make.html)的文档。现在还有一个新的工具，叫做CMake，有兴趣的可以去看一下，貌似其让编译文件的编写变得简单多了。

<!--more-->

GNU Make的文档里给出的第一个例子，能把Makefile的基本语法规则说明白了，例子如下：

     edit : main.o kbd.o command.o display.o \
            insert.o search.o files.o utils.o
             cc -o edit main.o kbd.o command.o display.o \
                        insert.o search.o files.o utils.o

     main.o : main.c defs.h
             cc -c main.c
     kbd.o : kbd.c defs.h command.h
             cc -c kbd.c
     command.o : command.c defs.h command.h
             cc -c command.c
     display.o : display.c defs.h buffer.h
             cc -c display.c
     insert.o : insert.c defs.h buffer.h
             cc -c insert.c
     search.o : search.c defs.h buffer.h
             cc -c search.c
     files.o : files.c defs.h buffer.h command.h
             cc -c files.c
     utils.o : utils.c defs.h
             cc -c utils.c
     clean :
             rm edit main.o kbd.o command.o display.o \
                insert.o search.o files.o utils.o

上面这个例子里，位于冒号左侧的，我们称之为target，target可以表示这条makefile语句要生成的文件名，如上面的edit和main.o，都表示我们要生存一个edit的可执行文件，和一个main.o的编译文件。位于冒号右侧的第一行表示生成这个文件需要的文件（注意，Makefile类似于C语言里面的宏，只有单行，只是为了人类的方便，用\表示换行）。第二行，如cc -c main.c 这一行，千万注意要以Tab开头，然后才开始书写。这一行表示生成这个文件所要执行的命令。我的Vim是配置了用四个空格代替Tab的，书写makefile就遇到麻烦了。还好makefile提供了.RECIPEPREFIX关键字，可以用它来替换默认的行为，如：.RECIPEPREFIX = > ，就表示用 > 符号来代替默认的Tab。

但是，clean这一行却有所不同，它冒号的右侧是没有第一行的。这种可以称之为伪目标，即不生成文件，只执行命令。 这个clean就是清理生成文件的作用。

对着上面这个makefile，执行make的结果是什么呢？其实是生成了一个edit文件。make如果没有指定目标的话，默认执行第一个目标，也就是edit。

上面这个例子，潜在的知识还有很多，暂时不表，回到Redis的Makefile来看看现实与理想的差距。

Redis的主makefile只有这么几行：

    default:all
    .DEFAULT:
        cd src && $(MAKE) $@

    install:
        cd src && $(MAKE) $@

    .PHONY:install

这里首先定义了一个名为default的target，但是其冒号右侧的第一行却是all。这里还有两个.DEFAULT 和 .PHONY这样的关键词。首先理解一下.DEFAULT和.PHONY，.DEFAULT表示找不到匹配规则时，就执行该命令。.PHONY表示强制伪目标，比如如果目录下游install这个文件，make可能就不执行了，因为它看到install没有变化过，就不需要执行命令编译了。而.PHONY则表示，这是一个伪目标，不管存在install文件与否，都强制执行。这些built-in的关键字，可以看这个链接[Special Built-in Target Names](http://www.gnu.org/software/make/manual/make.html#Special-Targets)。

$@符号的作用，表示让这里命令执行的操作的名字，比如我输入make hello,上面这个文件找不到匹配，执行默认的.DEFAULT，所以$@的值就是hello了。Makefile里面还有很多这样特殊的变量，即 [Automatic Variables](http://www.gnu.org/software/make/manual/make.html#Automatic-Variables) 。知道这些奇奇怪怪的符号叫自动变量，其他的问题有谷歌就行了。当时我不知道这种东西叫自动变量，真的搜索了很久才查到完整的信息。

这个$(MAKE)是什么意思呢？很显然是一个make命令，为什么不直接用make呢？具体可参考[How the MAKE Variable Works](http://www.gnu.org/software/make/manual/make.html#MAKE-Variable) 。原文有一句：Recursive make commands should always use the variable MAKE, not the explicit command name ‘make’； 什么事Recursive？其实像上面这个makefile，进入到子目录下执行子目录下的makefile，就叫做 [Recursive Use](http://www.gnu.org/software/make/manual/make.html#Recursion)。

makefile默认的规则是：第一个target就是make命令默认的target。即我们执行make，不加任何参数时所对应的就是第一个target。所以上面这个makefile的默认target就是default。

执行make的时候，观察到输出是：cd src && make all。 为什么这里会进入到.DEFAULT里面呢？其实还是要去看.DEFAULT的详细定义：“ .DEFAULT: The rule for this target is used to process a target when there is no other entry for it, and no implicit rule for building it” 。对于上面这个makefile，主要是 no implicit rule for building it 这一句，default:all 这一句是没有building rule的，所以还是会跑到.DEFAULT里面去，同时把all带入给$@。

现在进入到src目录下的makefile看看。开头第一句就是：release_hdr := $(shell sh -c './mkreleasehdr.sh'。很显然，这里执行了一个shell命令，并把值赋给了release_hdr。 看到 := 这个符号，如果你写过 Go 语言，就会很熟悉的。makefile变量定义有两种方式，另一种是 = ,其实还有一种 ?= ,可看：[The Two Flavors of Variables](http://www.gnu.org/software/make/manual/make.html#Flavors) 。 关于 shell 命令，可看[The shell Function](http://www.gnu.org/software/make/manual/make.html#Shell-Function)。

关于变量的定义，有这样两句：

    OPTIMIZATION?=-O2
    OPT=$(OPTIMIZATION)

上面已经给出了变量赋值的三种形式，取变量值则用 $(var)的方式，具体可看[Variables Make Makefiles Simpler](http://www.gnu.org/software/make/manual/make.html#Variables-Simplify)

Redis的文档里，有这么一种执行make的方法：make MALLOC=jemalloc。 这种方式叫做[overriding variables](http://www.gnu.org/software/make/manual/make.html#index-overriding-variables-with-arguments-728), 这样的方式，不管makefile里面如花定义MALLOC，最后都会被重写为jemalloc。

在Redis的makefile里面，有LDFLAGS，CC，CXX等几个变量没有定义，却可以使用。这样的变量称之为[Variables Used by Implicit Rules](http://www.gnu.org/software/make/manual/make.html#Implicit-Variables)。

关于如何在makefile里面写 if-else，可参考 [Conditional Parts of Makefiles](http://www.gnu.org/software/make/manual/make.html#Conditionals)。

Redis的makefile里有一句：

    %.o:%.c .make-prerequisites
        $(REDIS_CC) -c $<

这里出现的%.o , %.c叫做模式匹配，%表示匹配任意字符串的意思，可以参考：[Introduction to Pattern Rules](http://www.gnu.org/software/make/manual/make.html#Pattern-Intro)。$< 符号的意思，上面已经提到过[Automatic Variables](http://www.gnu.org/software/make/manual/make.html#Automatic-Variables)，在这里面找就行了。实际上我写这篇文章的初衷就是提供makefile的基本规则，基本关键字等的知识和出处，以及对应的英文名称，这样你遇到了其他类似的，就能很快找到其对应的含义了。之前我不太了解这些，浪费了许多的时间。

还有一个知识点，是引入其他makefile的，[Including Other Makefiles](http://www.gnu.org/software/make/manual/make.html#Include)

最后一个知识点是@符号，可参考：[Recipe Echoing](http://www.gnu.org/software/make/manual/make.html#Echoing)。 默认make会把每一行输出， 加@开头表示不输出。

至此，只要顺着上面提到过的术语和链接，翻一翻GNU make的文档，就完全能把Redis的makefile看明白了，leveldb的Makefile就更不在话下了。