

上次看完 weed-fs 在存储处理上的代码之后，原想着能很快把分布式处理这一块的代码看完的，没想到到现在还是处于半懵懂状态，只能边写边整理思路了。

由于我经验尚浅，加上学习上走了不少的弯路，所以对很多东西真的是完全不懂。像 weed-fs 这样子的东西，它们是如何处理分布式的，如何把数据分不到不到结点上的，如何保证数据一致性的等等这些问题我都不知道答案。因此这次刚好借着 weed-fs，以点带面，逐步把这一块的相关知识给补全了。

weed-fs 分布式处理的代码，主要在 topology 和 weed/weed_server 这两个包里面。 weed_server 主要提供外部以及结点间访问接口， 而 topology 则是表达结点特征以及处理结点间关系的。

<!--more-->

注意以下行文提及 server 时， 分两种，volume_server 指的是 VolumeServer 这个结构的完整实现。

weed_server 的代码布局如图：

![weed_server](http://h.hiphotos.bdimg.com/album/s%3D550%3Bq%3D90%3Bc%3Dxiangce%2C100%2C100/sign=08abef6141a7d933bba8e4769d70a02e/0e2442a7d933c895d65b4d0ed21373f082020099.jpg?referer=ddacaaa82d738bd49d3687016078&x=.jpg)

filer_server 上篇文章已经提及，此处不表。

关于分布式处理的，就主要是三种 server： master_server, volume_server, raft_server。 这里 master_server 和 volume_server 都比较好理解， 所有的请求都先放给 master_server, 然后 master_server 再决定发送给那个 volume_server。

关于 raft_server, 则不引申过多，首先要记住的是我们这里使用了一个 GO 语言实现的 Raft 协议，即 goraft, 为了使用 goraft，使得它能为 weed-fs 这种结构服务，那就需要告诉它一些东西，比如配置文件放哪了，你数据结点的拓扑结构是怎么样的， 发过来的命令要怎么处理等等。 所以 raft_server 干的就是初始化好一个 raft 客户端 ， 然后 Set 到 MasterServer 上， 让 master_server 一致性判断等功能。

所以，很明显，master_server 是这部分代码的核心，因此理解代码也先从 master_server 开始。 但是要理解 server 部分的代码，又要先了解 topology 部分的代码。 而了解 topology 部分的代码，又需要知道 server 部分的代码，只有这样才能够了解 topology 部分的代码是用来干什么的。 正是因为这种交叉关联，所以读懂代码真的是好难呀。

首先，要明确了解的一点：weed-fs 是可以同时启动多个 master_server 的，然后由 Raft 选出一个 Leader。 你发送请求到指定的 master_server 上去，然后每个 master_server 都会把请求转发到 Leader 这个 master_server 上来，由它发送给各个 volume 结点。

先来看看 MasterServer 的结构：

    type MasterServer struct {
        port                    int
        metaFolder              string
        volumeSizeLimitMB       uint
        pulseSeconds            int
        defaultReplicaPlacement string
        garbageThreshold        string
        whiteList               []string

        Topo   *topology.Topology
        vg     *topology.VolumeGrowth
        vgLock sync.Mutex

        bounedLeaderChan chan int
    }

port 是指这个 MasterServer 开放给外部的端口；

metaFolder 是启动 master_server 时指定的文件夹，里面放着这个 master_server 的元数据。命令 weed master -mdir="." 就会在当前文件夹上启动一个 master_server 。里面有几个初始文件，看看就大概明白。

pulseSeconds 是指心跳指令发送时间间隔， 主要是用于 Leader 监控各个结点，把死结点从网络中移除出去。

defaultReplicaPlacement 是 001, 200, 110 等等这样的字符串，表示一种复制策略， 具体含义可以去看 weed-fs 的文档，上面的说明已经很清晰了。

whiteList 是信任白名单，为了安全而设置的，暂时不管。

bounedLeaderChan 起到类似于锁的功能，于我们要理解的部分关系不大，暂时不管。

Topo 和 vg 这两个就难点了，特别是 vg， 到底 VolumeGrowth 的作用是什么，现在我还不是很清楚。

现在则要回过头来看看 weed-fs 的 Replica （复制集）配置脚本：

    <Configuration>
        <Topology>
            <DataCenter name="dc1">
                <Rack name="rack1">
                    <Ip>192.168.1.1</Ip>
                </Rack>
            </DataCenter>
            <DataCenter name="dc2">
                <Rack name="rack1">
                    <Ip>192.168.1.2</Ip>
                </Rack>
                <Rack name="rack2">
                    <Ip>192.168.1.3</Ip>
                    <Ip>192.168.1.4</Ip>
                </Rack>
            </DataCenter>
        </Topology>
    </Configuration>

这里是把整个搭建起来的 weed-fs 分布式网络的结点布局称之为 Topology (拓扑），完整的布局下有 DataCenter （数据中心）， 每个数据中心下会有多个 Rack （支架，可理解为服务器机器集中放置的机柜）， 每个 Rack 下有多个 DataNode (结点)。

所以，这里对 Topo 的理解就明白多了。 vg 暂时先不管，后面再提及。

现在来看看 VolumeServer 的结构：

    type VolumeServer struct {
        masterNode   string
        pulseSeconds int
        dataCenter   string
        rack         string
        whiteList    []string
        store        *storage.Store
        FixJpgOrientation bool
    }

从上可以看出， 每个 volume_server 都和一个 master_server 联系，同时也记录自己属于哪个 rack，哪个 dataCenter。

FixJpgOrientation 是存储图片时是否要对图片进行处理的选项，这里不管。

这里值得注意的是 store， storage.Store 是什么？ 上一篇在存储的关系图里面已经提及了 Store 了。 实际上，我们可以把 storage.Store 看做是每个 volume_server 里对硬盘上数据的管理者。

启动一个 volume_server 的命令： weed volume -max=100,90 -mserver="localhost:9333" -dir="./data1,./data2" 。

这里实际上是将 data1 和 data2 这两个文件夹都作为 volume 数据，即上篇所说的 dat 数据存放处， 用逗号分隔。 -max 表示这个文件夹里可以生成的 dat 文件数量。 mserver 指定了 master_server

这里可以看到， 每个 volume_server 其实都关联着一个 master_server, 同时标记自己属于哪个 rack, 哪个 dataCenter。

现在则又牵扯到另一个问题，即 rack 和 dataCenter 到底是怎么样的存在，以及他们与 volume_server , master_server ，raft_server 的存在之间又有什么区别？

总的来说是这样的： Topology, DataCenter, Rack, DataNode 是模拟物理结点的拓扑结构，而 volume_server 则是 DataNode 的对外接口，每个 volume_server 都与一个指定的 master_server 联系，但是如何处理则交给 master_server 里中的 Leader 。

那么，问题来了：

*   启动一个 master server 和启动第二个 master server 并加入到网络中到底发生了什么？

*   启动一个 volume 并加入到网络中又发生了什么？

*   一个数据操作请求发送过来，会经过哪些个步骤？

*   一个 master server 结点死了，会发生什么事情？

*   一个 volume server 结点死了，又会发生什么事情？

先来看第一个问题，启动一个 master server 的时候发生了什么，这部分代码可以在 weed/master.go 和 weed/weed_server/master_server.go 看到。

首先，在 master.go 里可以看到， 这里是先 NewMasterServer，然后 启动一个 NewRaftServer 并 set 到master server 上， 自动选择当前master server 为 Leader。然后去加载 配置 xml，勾勒出当前服务器的拓扑结构。

当前的 master server 被初始化为 Leader 的代码就只是这几句：

    _, err := s.raftServer.Do(&raft.DefaultJoinCommand{
        Name: s.raftServer.Name(),
        ConnectionString: "http://" + s.httpAddr,
    })

这个 raftServer 是 goraft 实现的那个。上面基本没有文档，所以暂时还不太清楚具体缘由。

那启动了一个 master server， 再启动一个，与第一个连接，有发生什么事情呢？

按照如下命令启动两个 master server：

    weed master -defaultReplication=010
    weed master -port=9334 -peers="localhost:9333"

这里其实我还有两个疑问的， 上一篇文章我说 volume 上会保存 Replication 模式数据，而这里启动的 master server 也设置了一个，这两者什么关系呢？ 启动第二个 master server 不加 defaultReplication 参数，会有什么不同吗？暂时先把问题搁置一下，后面看看能否解答。

看了代码，主要的过程是这样的：

先调用 raft_server.go 的 Join 函数，然后到 peers master server的 raft_server_handlers.go 的 joinHandler 函数，然后执行一句：

    if _, err:= s.raftServer.Do(command);

这里应该是试图让当前启动的 master server 成为 Leader， 然后 raft 协议内部就进行选举了。

说完了 master server， 就到了 volume server 如何加入到网络中的问题了。这部分的代码相当庞大和复杂，许多地方我还没有看明白。 大致过程就是把一个 volume 封装为一个 DataNode， 然后加入到 master server 里的 Topo 结构中去， 这个结构应该是用一个字典来保存映射关系的。

由于不是很清晰，所以关于 volume server 的内容暂时不多说了。

写入数据的过程也大致类似。个人感觉这部分的代码比较复杂和奇怪，没怎么看懂，所以略过了。

那么 volume 的ReplicaPlacement 和 master server 自己的 ReplicaPlacement 有什么区别呢？

每次提交文件，都可以指定一种 replicaPlacement，如果没有设定，则默认使用 master server 设定的。而 volume 的 replicaPlacement 则主要在 volume_layout.go 这里使用，一个 volume 的增删改都要使用这个进行判断。这里则牵涉到 master server 如何根据 replicaPlacment 去找到对应 volume 的过程，因为同一个 volume id，两个server都会有，volume_layout 里就是根据 replica 数量要求，判断当前 volume id全部加起来的 server 数量是否小于这个数量。

心跳指令的发送则在启动 server 的时候用一个 goroutine 一直 for 循环执行。

写到这就发现这篇文章越写越乱，其实我还有很多个地方没有搞清楚，但是由于看这份代码实在花了我太多的时间了，所以想先暂停一下。下一步我打算先了解一下 raft 协议的实现， 然后再用 Lua 实现一个简单版的 weed—fs ，之所以选择 Lua，是因为我们公司是超重度使用 Lua 的游戏公司，我希望一两年后可以去做游戏服务器开发。

要是用 Lua 重写的话，估计还要继续再看 weed-fs 的代码很多次，所以这里暂时不写了。等以后有了更深的体会，我会回来继续写的。