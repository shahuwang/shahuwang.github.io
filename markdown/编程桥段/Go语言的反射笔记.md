如果去阅读Go语言写的一些库，特别是涉及大量数据变化的代码，会充斥着大量的reflect代码。因为Go并没有泛型，也不支持宏，所以需要类型转换的复杂技巧代码，都需要用到reflect。对Go的反射一知半解，每次用都得去查找相关文章，总是记不牢，理不清，所以这次打算写篇文章，把反射给搞懂了。

首先，需要理解的是Go的反射三大定律：


1. 通过interface获得反射对象

    每个interface值，实际上都带着一对(value, type)，value代表此interface的值，type表示原始类型。reflect.ValueOf  和 reflect.TypeOf 就是获取这对数据出来，并封装在 reflect.Value 和 reflect.Type 里。
    
    reflect.Kind 表示各种 type, 如 int32, Array, Map, Struct, Ptr, 说的是这个 interface 的值的是个什么东西。
    
    A Kind represents the specific kind of type that a Type represents. 这句话不知道该怎么翻译。
    
    Type和 Value都有一个 Kind 的方法，它会返回一个常量，表示底层数据的类型。这其实是最困惑人的地方，为什么两者都有这个方法，而且返回值还是一样的呢？感觉这是reflect设计不合理的地方
    
2. 通过反射对象能还原回interface{}对象

    reflect.Value 对象有一个 Interface 方法，与 reflect.ValueOf 作用相反，是获得 Value 对象的interface{}
    
    
3. 如果要修改反射对象，这个Value对象必须是可写的

    什么对象是可写的呢？指针类型的咯，就是 reflect.ValueOf(variable), 你传的必须是引用，而不是值复制。
    rv 不是可写的，因为 rv 是指针，它自己本身无法被修改，但是它指向的值能够被修改，通过Elem可以获得指针实际指向的值。
    
    
  	type Ref struct{}
	r := new(Ref)
	rv := reflect.ValueOf(r)
	fmt.Println(rv.CanSet() == false)
	fmt.Println(rv.Elem().CanSet() == true)


现在来了解下reflect.Type所具备的所有方法的具体含义。

`Type.Align()int` , 返回该类型在以多少字节来进行内存对齐

`Type.FieldAlign()int`, 返回该类型的值作为某个结构体的字段时，在内存对齐时是采取多少字节对齐的

`Type.Method(i int)Method`, 返回该值的第i个方法

`Type.MethodByName(name string)Method`, 返回该类型名为name的方法

`Type.NumMethod()int`, 返回该类型的方法数量

`Type.Name()string`, 返回该类型的名称

`Type.PkgPath()string`, 返回该类型所在的包名

`Type.Size()int`, 返回存储该类型的一个值所需要的字节数，但不包含实际数据的字节。譬如 `type Name struct{name : string}`, name字段存储的字符串长度对返回值并不会有什么改变。

`Type.Elem()Type`, 返回该类型（为指针）实际指向的数据

`Type.Field(i int)StructField` 返回一个结构体的第i个字段

`Type.In(i int)Type` 返回一个为函数Type的第i个参数

Go中的Select块，以及channel，都具有各自的反射对象以及使用方式

`reflect.Select` 用于以反射的方式执行 Select 操作 

发现 `reflect.Indirect` 的作用和 Elem 差不多。。。

写着写着，发现要写的东西太繁杂了，懒了，不想写了。

不过扫了几遍 reflect 包的文档，对它的体会深了不少，总得来说，reflect 提供了两方面的功能：取正常代码执行时的中间表达；直接构造中间表达来实现正常代码无法表达的内容。
