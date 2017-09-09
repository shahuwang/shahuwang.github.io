本文基本上算是对[Lua手册](https://www.lua.org/manual/5.1/manual.html#pdf-require)关于require机制的翻译。

现在工作上主要使用openresty，也是很惊叹这套基于Lua和Nginx的工具链如此之好用，又如此之强大。为了不再蹉跎岁月，打算把这套工具链精进研究一下。

昨晚突然想知道require机制到底是怎么个过程，所以查找了一下相关的资料，发现还是官方手册最完善。

##### require(modname)

require函数首先会去package.loaded这个table以modname做key查找对应的值，如果存在这个值，就返回。如果没有，则会去找一个loader来加载这个模块。

require会去package.loaders这个数组里面找loader，所以如果你想写自己的加载器，修改这个数组就行啦。

找到loader之后，require会把loader获取到的值，赋值给package.loaded[modname]。如果loader没有返回值，则赋值为true。

##### package.cpath

这是让require去查找C语言模块的。Lua在启动的时候，会根据系统的环境变量LUA\_CPATH设置package.cpath=LUA\_CPATH, 如果不存在这个环境变量，则使用默认的值，定义在luaconf.h里面的。

每个路径使用 `;` 分割开来，如`package.cpath= "./?.so;./?.dll;/usr/local/?/init.so"`, 那么 `require "foo"`, 会尝试打开 `./foo.so, ./foo.dll` 和 `/usr/local/foo/init.so` 。如果找到需要的C库，将会使用动态链接连接起来，然后查找里面以luaopen\_ 开头的函数，比如 `require "a.b.c" 会找到函数 luaopen\_a\_b\_c。

##### package.path

这是让require去查找Lua模块的。Lua在启动的时候，会根据系统变量LUA\_PATH设置`package.path=LUA_PATH`, 如果这个环境变量没有定义，则使用luaconf.h里面定义的路径。 如果你设置`LUA_PATH=";;"`, 这个值不会被使用，还会照旧使用默认路径。

每一个路径用 `;` 分割开来，如 `package.path= "./?.lua;./?.lc;/usr/local/?/init.lua"`, 那么 `require "foo"` 的时候，loader会尝试打开 `./foo.lua, ./foo.1c和/usr/local/foo/init.lua`。

##### package.loaded

这个表放置了所有加载过的模块，如果require(modname), 则先看package.loaded[modname]的返回值，如果是nil或者false，则表示这个模块没有加载。其他值表示模块已经加载，返回此值。

##### package.loaders

这是供require控制模块加载的table，每一个值都是一个查找函数。 Lua初始化这个表时提供四个函数，第一个查找函数会去package.preload表里面查找；第二个查找函数会去找一个用于加载Lua库的loader，这个loader会使用package.path路径去找寻要加载的模块；第三个查找函数会去找一个加载C模块的loader，这个loader会使用package.cpath作为查找路径；第四个查找函数会使用一个all-in-one loader, 主要用于C模块的查找，比如require a.b.c, 它会找一个名为a的C库，如果找到了，就看里面有没有一个luaopen\_a\_b\_c的open function。


对Lua的require机制了解了之后，就可以来了解Lua的热更新了。关于热更新的知识，主要来源于这几篇博客[https://asqbtcupid.github.io/](https://asqbtcupid.github.io/)

不知道其他人有没有疑问，我是一直有很大疑问的，每每说到Lua，Erlang之类的语言，都说它的特性支持热更新，但是我基本上没有找着篇文章，告诉我这些热更新操作到底该怎么做，具体步骤怎么样？这一点是相当奇怪的，只能说明，大部分的项目根本用不到热更新，所以大部分用这些语言的人，也基本上没有去钻研这个细节。

目前据我所了解到的Lua的热更新，主要有两个核心点：重新require，更新package.loaded里面的值；保留upvalue（啊，upvalue这个词也是很少解释说明，大概就是函数定义之前的变量）。比如如下这个模块, 如果我要热更新，把返回值修改成 `"hello world .. hotfix" .. count`, 就需要保证package.loaded里面的是新的代码，新的代码运行的时候能保留count的值（count就是函数test的upvalue)。

```lua
--- base.lua
local _M = {}
local count = 0
function _M.test()
    count = count + 1
    return "hello world..  " .. count
end
return _M
```

如何做到更换旧代码并且保留upvalue呢？代码如下：

```lua
--- hotfix.lua
local oldfunc = require "base"
package.loaded['base'] = nil
local newfunc = require "base"

for i = 1, math.huge do
    local name, value = debug.getupvalue(oldfunc.test, i)
    if not name then break end
    debug.setupvalue(newfunc.test, i, value)
end
print("hotfix complete!")
```

关于如何访问upvalue的方法，可以看 [Programming Lua](https://www.lua.org/pil/23.1.2.html)

代码写好了，另一个疑问就来了，如何通知正在运行着的程序，去加载这段热更新代码呢？这方面也没看到有多少文章有说的，也是很奇怪。网上看到的一个做法，倒是值得学习，他让Lua进程监听某个信号，一到监听到信号，就去指定的路径加载指定的文件，这个文件可以是热更新的代码，也可以是处理热更新的某些逻辑。

但是，此时另一个问题又出现了，假如我修改了 Module\_a里面的test函数，要进行热更新。而Module\_b里面用到了Module\_a的test函数，Module\_c则用到了Module\_b 。 a => b => c 这样的一个引用链，现在我替换了 Module\_a 的代码，然后在 Module\_b 进行了如上的热更新逻辑，后续如果有新的代码运行 `require("Module_b")` 那么确实会使用到了新的代码逻辑。但是对于 Module\_c 来说，它在启动之时就已经 `local b = require("Module_b")`, 不管Module\_b 是否进行了热更新逻辑，Module\_c里面一直都是使用着旧的代码。

难道我为了热更新一个函数，需要把整条调用链上的代码都进行热更新逻辑处理吗？看着也是相当不合理的啊。找到这篇[《Lua热更新原理(4) - 替换函数》](https://asqbtcupid.github.io/luahotupdate4-callback/)，发现是可以替换掉全部函数的，核心原理就是“遍历虚拟机， 找到旧函数所有的索引，并把这些索引指向新函数”，不过说实在的，我并没有怎么看懂，这一块的内容，等到下篇博客再来写了。

如下是我用来测试热更新逻辑的四个代码文件，base.lua 就是热更新时要更改的函数代码所在；hotfix.lua 就是热更新操作过程；app.lua 引用了 base.lua, 主要是用来验证替换函数这个逻辑的；main.lua 操作的入口，这里先运行旧代码，睡眠10秒钟，把base.lua代码替换掉，进行热更新逻辑操作，比对 base.test(), app.run()的输出，是不是新的代码的输出。

```lua
--- base.lua

local _M = {}
local count = 0
function _M.test()
    count = count + 1
    return "hello world.. after hotfix " .. count
end
return _M
``` 

```lua
--- hotfix.lua

local oldmod = require "base"
package.loaded['base'] = nil
local newmod = require "base"

for i = 1, math.huge do
    local name, value = debug.getupvalue(oldmod.test, i)
    if not name then break end
    debug.setupvalue(newmod.test, i, value)
end

local function replace(oldfunction, newfunction)
    local visited = {}
    local function f(t)
        if not t or visited[t] then return end
        visited[t] = true
        if type(t) == "function" then
            for i = 1, math.huge do
                local name, value = debug.getupvalue(t, i)
               
                if not name then break end
                f(value)
            end
        elseif type(t) == "table" then
            f(debug.getmetatable(t))
            for k,v in pairs(t) do
                f(k); f(v);
                if type(v) == "function" or type(k) == "function" then
                    if v == oldfunction then t[k] = newfunction end
                    if k == oldfunction then 
                        t[newfunction] = t[k]
                        t[k] = nil
                    end
                end
            end
        end
    end
    f(_G)
    local registryTable = debug.getregistry()
    for k, v in pairs(registryTable) do
        if v == oldfunction then
            registryTable[k] = newfunction
        end
    end
end

replace(oldmod.test, newmod.test)
print("hotfix complete!")
```

```lua
--- app.lua

local base = require "base"
local _M = {}
function _M.run() 
    local ret = base.test()
    print(ret)
end
return _M
```

```lua
--- main.lua

local app = require "app"
app.run()
local function sleep(s)
    local ntime = os.clock() + s
    repeat until os.clock() > ntime
end

sleep(10) --- 趁睡眠时，去更改base.lua里的test的返回值

local hotfix = require "hotfix"

local base = require "base"
print(base.test()) --- 热更新后的输出


app = require "app"
app.run() --- 也是热更新后的新输出

```