现在许多 NoSQL 受谷歌的 BigTable 的设计思想影响，都采用 commitlog-->memtable-->sstable 的结果，所有进来的数据都优先写入到commitlog里，然后写入 memtable，一定程度之后sync到硬盘的sstable。只有commitlog写成功了，一次写操作才会是成功的。memtable和sstable没有写成功，都可以根据commitlog复盘出来。

今天要写的，是Cassandra的commitlog部分的代码。代码的开始从 cassandra/db/Keyspace.java 开始，在 apply 方法里，使用到了 CommitLog.instance.add(mutation); 这里便是 commitlog 使用的开始。

cassandra/db/Mutation.java 这个类，是对一行数据的封装，每次执行 CQL，插进来的每一行数据都被封装成了Mutation实例。

Mutation这个类花了我很多时间去看，还是没多看懂，毕竟涉及到其他的类太多了。不过知道其代表着一行数据就好了。

commitlog要做的，主要是三件事：添加新的数据，将数据sync到硬盘上，将硬盘里的数据复盘到内存里。

在看commitlog如何实现这三个功能前，需要了解一下其他几个类的作用和意义。

第一个要了解的是 
cassandra/db/commitlog/CommitLogSegment.java  这个类，commitlog实例，对应着多个log存储文件，CommitLogSegment 这个类就是对着这一个log文件。

CommitLogSegment  是一个抽象类，但其实际上已经实现了大部分的方法了，具体实现，主要是数据如何写和存储上（譬如压缩与否）。

CommitLogSegment  实现了一个内部类 Allocation，这个内部类的作用其实是记录这次add的数据放置的相关信息。

CommitLogSegment 的第一个作用是初始化log文件。其构造函数即初始化了一个log文件。

第二个作用是 allocate，传进来mutation和size，返回 Allocation的实例，即分配好这个数据要写的位置，然后由上层将数据写入这个位置（不是写入到log硬盘文件里）。

第三个作用是sync到硬盘文件上。

CommitLogSegment 里面的数据分两种， dirty和clean 。已经flush到硬盘的log是clean的，刚写进来，还没有写到硬盘上的，是dirty的。

CommitLogSegment 的两个具体实现是 CompressedSegment 和 MemoryMappedSegment, 前者是在写入到log文件时进行压缩，后者是进行文件的内存映射。

CommitLogSegmentManager 则负责管理众多 CommitLogSegment, 其在初始化的时候就启动一个线程，用几个队列维护着新增的segment，使用着的segment，以及等待要写入到硬盘的segment。同时，还对segment进行写入、回收或者删除。

CommitLogArchiver 是根据配置文件，对commitlog进行备份或者恢复，具体可看 
[http://docs.datastax.com/en/cassandra/2.0/cassandra/configuration/configLogArchive_t.html](http://docs.datastax.com/en/cassandra/2.0/cassandra/configuration/configLogArchive_t.html)

CommitLogDescriptor 主要存放Commitlog的一些元数据，以及log文件的写入与读取。

CommitLogReplayer 很明显，就是将log文件复盘到内存的具体实现过程。

AbstractCommitLogService 则是控制 commitlog 写操作的过程，如信号机制。具体的两种实现，一个是 BatchCommitLogService, 一个是 PeriodicCommitLogService ,两者的差异可以看这篇wiki介绍：[ArchitectureCommitLog](http://wiki.apache.org/cassandra/ArchitectureCommitLog)

这篇文章只是对 Commitlog 的代码做个大概的介绍，具体的实现细节，我自己也还没有看怎么清晰。不过目前先知道每个类主要是做什么的，这样到了抠细节的时候，就容易多了