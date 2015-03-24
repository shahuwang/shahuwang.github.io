weed-fs是用Go语言编写的一个分布式文件系统，主要用于存储大批量的小文件，因此其在存储结构的处理上，妙用了许多方法来达到这个目标，使得大批量的小文件存取更加高效。

正常的文件系统一般提供了文件夹定位，文件定位，数据定位这三步，所以在磁盘访问上会多了几步。而Weed-fs则直接提供文件ID到数据定位这一步。当然，为了兼容传统的文件系统接口，如FUSE，Weed-fs提供了一个filer模块，模仿文件夹功能。

下图是weed-fs的代码目录结构，代码量并不是很大，其中我觉得比较值得阅读的是filer，storage，topology和weed这四个。 filer实现了文件目录的功能， storage是实际存储的实现，topology则是分布式拓扑结构的实现， weed则是终端和web的接入口。

<!--more-->

![directory layout](http://a.hiphotos.bdimg.com/album/s%3D550%3Bq%3D90%3Bc%3Dxiangce%2C100%2C100/sign=c02bc808acc379317968862cdbffc678/f636afc379310a550bc2d4fab44543a9832610dd.jpg?referer=f4d351381e950a7b2c227af450b4&x=.jpg)

filer 既然是模仿传统上我们所认知的文件目录功能，那么就必然需要提供文件夹和文件的增删改查这四大功能。也因此，filer分两块实现，文件夹的增删改查，文件的增删改查，然后两者还要实现关联。

weed-fs认为文件夹的数量远远小于文件的数量，但是文件夹的访问次数又远远高于文件的次数，所以存储的处理上，文件夹是存放在内存中，而文件则是存放在leveldb中（go实现）。

文件夹的实现并不是很复杂，就是一大堆的字典组成一个目录树，而且，每一个文件夹的增删改操作都会记录到一个log文件里面去。每次启动的时候，都会先去加载这个log文件，根据上面记录的增删改操作在内存中构建出原有的目录树来。

文件增删改查的实现上，则要理解一点，即file id，在weed-fs中，每个存进来的文件都会给你一个file id，如果不启用 filer 功能的话，你要找一个文件是要通过 file id来查找的，而不是根据文件路径名来找到的。

创建文件的时候，要传入两个参数，分别是 filepath和file id，filepath是类似于这样的路径：/home/test/a.txt 。然后将这个filepath 分成两部分，即/home/test/ 部分和 a.txt部分。用前者去文件夹里面查找，如果存在就返回文件夹id，不存在，则创建一个并返回文件夹id。然后用文件夹id和文件名生成一个唯一id，与该文件对应的file id组成一个键值对存储在leveldb里面。所以，如果启用了filer功能，那么操作文件的时候，先根据你提供的filepath生成一个key，根据这个key在leveldb里面找到文件的file id，然后再根据file id找到文件数据。

现在来看 weed-fs的实际存储上的设计，下图是storage主要的结构：

![storage](http://f.hiphotos.bdimg.com/album/s%3D550%3Bq%3D90%3Bc%3Dxiangce%2C100%2C100/sign=02cc350f0955b31998f982707392f31b/78310a55b319ebc41a4a0f978126cffc1e171689.jpg?referer=b09a464501087bf424fb62d95288&x=.jpg)

在启动的一个weed-fs服务里，会有一个store结构，下有多个DiskLocation，即同一个服务，可以设置多个存储文件夹。每个DiskLocation下面，又会有多个volume，每个volume下面，如果开启了collection功能，则每个volume下会有多个collection。每个collection或者volume下面，主要有两种文件，一个是存放实际数据的dat文件，一个是存放文件索引的idx文件。如果开启了collection，则文件的命名是collectionname_1.dat , collectionname_2.dat, collectionname_1.idx, collectoinname_2.idx 等一系列这样的名字。后面的这个数字就是volume id。

启动一个 volume server，在指定的文件夹上预先生成7个volume。

其他的还有 .cpd 和 .cpx 即将dat和idx压缩后的数据，另有cdb这种东西，其中的cdb部分代码是将idx和cdb互换的，cdb是一种[constant databases](http://cr.yp.to/cdb.html) ，当一个volume被转换为readOnly形式的时候，一般会先将idx转换为cdb，具体为什么要采用cdb，不是很清楚，也是是快加上数据不能更改吧。

每个上传上来的文件，都被封装为 Needle 这样的一个数据结构：

![Needle](http://a.hiphotos.bdimg.com/album/s%3D550%3Bq%3D90%3Bc%3Dxiangce%2C100%2C100/sign=0f122db584d6277fed12323d18036e0d/a08b87d6277f9e2f8f1ba4061c30e924b899f396.jpg?referer=ce3085c9935298225c240cf33d8c&x=.jpg)

每个dat文件在读取到内存的时候，都被封装成Volume这个数据结构了：

    type Volume struct{
        Id VolumeId
        dir string
        Collection string
        dataFile *os.File
        nm NeedleMapper
        readOnly bool
        SuperBlock
        accessLock sync.Mutex
    }

每个idx文件，即文件索引文件，在使用的时候都会被加载到内存里面，构造出 NeedleMap 这个数据结构：

    type NeedleMap struct{
        indexFile *os.File
        m CompactMap
        mapMetric
    }

每个在索引文件里的文件索引记录，都被封装为NeedleValue这个数据结构：

    type NeedleValue struct{
        Key Key
        Offset uint32
        Size uint32
    }

weed-fs存储文件的过程是这样的，你先申请一个file key，然后拿着这个file key和你的文件发送过来存储。一个file key 是这样的：3,01637037d6，逗号前面的这个3，就是NeedleValue的这个key的值了。索引文件用一个长度为100000的数组来保存NeedleValue，按照Key值从小到大排序。不知道为什么，如果一个索引文件记录的文件数超过十万，weed-fs则用一个字典来记录超过的部分。也许是超过十万之后，二分查找的效率不如字典了吧。

同样，文件存储到dat里面也是线性Append过去的，以8 bytes为最小的一个数据块，所以NeedleValue里面的Offset表示这个文件在dat文件里存储的起始位置。Size 则表示文件的大小，根据这两个就能到dat文件里面找到我们要的那个文件了。

索引的插入和写入到索引文件里的代码主要在NeedleMap的Put方法里面，代码主要是这几行：

    bytes := make([]byte, 16)
    util.Uint64toBytes(bytes[0:8], key)
    util.Uint32toBytes(bytes[8:12], offset)
    util.Uint32toBytes(bytes[12:16], size)
    nm.indexFile.Write(bytes)

根据上面的这个规则，同样也可以轻松得从索引文件构造出索引来。

dat文件的写入过程主要在 Needld的 Append方法里面，里面提供了两个版本的写入，这里主要看第一个版本，因为比较简单明了，容易说明文件存储过程。写入过程从volume的write方法里开始，先通过

    offset, err:= v.dataFile.Seek(0,2)

找到dat文件最末尾的位置，这个位置也就是将要插入文件的起始位置。

然后调用Append方法，将数据写入，再根据 key，offset，file.size构建一个NeedleValue，插入到索引里面，以便后续查找。

在Append方法的里，每个文件写入到dat文件里，都包含三个数据，header，即文件的元数据：

    header := make([]byte, 16)
    util.Uint32toBytes(header[0:4], n.Cookie)
    util.Uint64toBytes(header[4:12], n.Id)
    util.Uint32toBytes(header[12:16], n.Size)
    w.Write(header)

这里的n.Cookie和n.Id的值，是通过 ParseKeyHash这个方法处理得来的。之前说过，上传一个文件需要一个fileid（简写为fid），即类似于 3,01637037d6这里的，则 n.Id, n.Cookie, _ := ParseKeyHash(01637037d6) 获得的。

关于一个 fid是如何生成的，可以去看 file_id.go这个代码。

把文件的元数据作为头部写入到dat之后，再将文件本身数据写入到dat里面。

然后接着，写入一个唯一校验码和补全零值，保证每个8 bytes的数据块都被填满了。否则，offset这个计数的功能就没有意义了，因为dat的最后位置不是 8 bytes的倍数，无法计算位置了。

看到文件系统相关资料的，应该都知道，很多文件系统的开始部分都是一个Super Block，这个Super Block记录了这个文件系统的一些元数据。weed-fs也为每个volume（即dat）初始化了一个长为8 bytes的Super Block，用来存储这个Volume的版本，生存时限（TTL），复制集模式。

下一篇文章将主要介绍 topology 这部分的代码设计。