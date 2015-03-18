之前看 Go 语言的项目，很快懂了不少。最近看 Cassandra 的代码，Java 各种继承绕得我头都晕了，实在是好难懂。

目前对 Cassandra 的认知还停留在比较初始的阶段，所以好多东西都看不明白。看代码的时候，找着注释的说明，写了如下一个测试程序：

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
            String schema = "create table test2.myTable2(k int primary key, v1 text, v2 int);";
            String insert = "insert into test2.myTable2(k, v1, v2) values (?, ?, ?)";
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


代码运行后，在指定的文件夹下，生成了如下图所示的几个文件：
![sstable files](https://raw.githubusercontent.com/shahuwang/images/master/Cassandra/sstable.png)

看了这些文件之后，我主要想知道各个文件所存储的内容，文件名里各个部分的意义，以及这些文件名生成的过程。

各个文件的意义，在cassandra/io/sstable/Component.java 里面有讲解：

+  Data.db 是 SSTable 的基础数据，其他所有的文件都可以根据这个文件重新生成出来。
+ Index.db 是 row key的索引，具有指向其实际数据位置的指针
+ Filter.db 是序列化后的布隆过滤器，处理的是 row keys.
+ CompressionInfo.db 存储未压缩数据长度，chunk 偏移量等数据
+ Statistics.db SSTable 上内容的统计元数据
+ Digest.sha1 存储data file 的 sha1 和
+ TOC.txt , 存储 SSTable 的所有 component 的列表

文件名的第一部分： la， 代表 SSTable的Version，在 cassandra/io/sstable/format/big/BigFormat.java 可以看到，la 是3.0.0 版本的。

数字 1 代表当前 SSTable 的generation，1是第一代的意思。

Big 是 SSTable 的格式名。

这个名字生成的开始是 cassandra/io/sstable/Descriptor.java 的 filenameFor(Component component)。

generation 的作用是什么？ Big 这种格式有什么新内涵？ 下一篇再继续说吧