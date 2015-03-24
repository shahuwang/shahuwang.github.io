之前看到generation的时候，就猜测很有可能是同一个表，重新创建之后保存为新的一代。然后我就去做测试，但是一直未能验证我的想法，每次运行代码之后都重新创建一个新的表，generation都是1。直到一次偶然，居然发现能够创建到generation 为2， 为3 等等的情况了。重新去看了下代码，发现自己被自己坑了。

依然使用上一篇文章里的测试代码，使用 CQL 写入 SSTable，代码如下：

    package com.shahuwang;

    import org.apache.cassandra.io.sstable.CQLSSTableWriter;
    import org.junit.Test;
    
    import java.io.IOException;
    
    /**
     * Created by rickey on 2015/3/11.
     */
    public class sstableTest {
        @Test
        public void cqlWriter(){
            String schema = "create table test.test(k int primary key, v1 text, v2 int);";
            String insert = "insert into test.test(k, v1, v2) values (?, ?, ?)";
            CQLSSTableWriter writer = CQLSSTableWriter.builder().inDirectory("E:\\share\\github\\cassandra\\data\\test").forTable(schema).using(insert).build();
            try{
                writer.addRow(0, "test1", 24);
                writer.addRow(1, "test2", null);
                writer.addRow(2, "test3", 42);
                writer.close();
            }catch (IOException e){
                System.out.println(e);
            }
    
    
        }
    }

这份代码，运行第一次，生成的文件是这样的：
![第一次运行生成文件](https://raw.githubusercontent.com/shahuwang/images/master/Cassandra/sstable.png)

而第二次运行，第三次运行，则生成的文件如下图：
![第一次运行生成文件](https://raw.githubusercontent.com/shahuwang/images/master/Cassandra/SSTable2.png)

可以很明显看到generation 有1， 2， 3 三种了。因为每次运行都执行 CQL 语句 

`create table test.test(k int primary key, v1 text, v2 int);`

每次创建一个column family，名为 test，由于之前已经存在了，所以新的 column family 就用 generation 进行区分。

那么这个 generation 生成的依据是什么呢？这也是我被坑的地方。要判断 generation 是否要增大，要增大到多少，首先需要确定你的 CQL 语句的 cfName（column family name）是什么，然后再确定你指定的那个文件夹对应的 cfName 是否与 CQL 的 cfName 一致，如果一致，说明你要覆盖掉以前创建的 column family，所以要增大 generation。 要增大到多大呢？此时就根据里面的文件里的generation 来判断了，譬如当前是2，那么就增大到3。

Cassandra 默认根据 CQL 的 key space 和 column family 创建存储数据的文件夹，譬如上面的这个 CQL 语句，Cassandra 会默认创建文件夹是 test/test, 第一个test是 key space name,  第二个 test 是column family name。我就是这里被坑的，由于我代码里指定的文件夹，与我 CQL 里面指定的 keyspace , column family 不一致， 导致我每次运行测试代码，generation 都是 1。

根据指定的文件夹提取 cfName 和 当前最大 generation 的代码主要是 Cassandra/io/sstable/Descriptor.java 里面的  Pair < Descriptor, String > fromFilename(File directory, String name, boolean skipComponent) 方法。设置几个断点在这里，就可以看明白了。

通过 CQL 获取到 cfName 的部分暂时不去考虑，这段时间主要聚焦于 SSTable。

而生成当前所需的 generation 的代码主要在 Cassandra/io/sstable/AbstractSSTableSimpleWriter.java 里面。

下面是我抽取出来的一段测试代码，可以通过这个测试代码看看 generation 的生成过程：

    import org.apache.cassandra.io.sstable.Component;
    import org.apache.cassandra.io.sstable.Descriptor;
    import org.apache.cassandra.io.sstable.SSTable;
    import org.apache.cassandra.utils.Pair;
    import org.junit.Test;
    
    import java.io.File;
    import java.io.FilenameFilter;
    
    import java.util.HashSet;
    import java.util.Set;
    import java.util.concurrent.atomic.AtomicInteger;
    
    /**
     * Created by rickey on 2015/3/11.
     */
    public class GenerationTest {
        protected static AtomicInteger generation = new AtomicInteger(0);
        @Test
        public void generation(){
            // 第二个参数 test 是 cfName
            int maxGen = getNextGeneration(new File("E:\\share\\github\\cassandra\\data\\test"), "test");
            System.out.println(maxGen);
        }
        private static int getNextGeneration(File directory, final String columnFamily)
        {
            final Set<Descriptor> existing = new HashSet<>();
            // 这里巧妙的利用了接口做类似lambda表达式的结果
            directory.list(new FilenameFilter()
            {
                public boolean accept(File dir, String name)
                {
                    Pair<Descriptor, Component> p = SSTable.tryComponentFromFilename(dir, name);
    
                    Descriptor desc = p == null ? null : p.left;
                    if (desc == null)
                        return false;
                    System.out.println(desc.cfname);
                    if (desc.cfname.equals(columnFamily))
                        existing.add(desc);
    
                    return false;
                }
            });
            int maxGen = generation.getAndIncrement();
            System.out.println(existing.size());
            for (Descriptor desc : existing)
            {
                while (desc.generation > maxGen)
                {
                    maxGen = generation.getAndIncrement();
                }
            }
            return maxGen;
        }
    }