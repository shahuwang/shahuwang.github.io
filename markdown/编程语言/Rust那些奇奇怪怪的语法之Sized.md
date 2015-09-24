Rust那些奇奇怪怪的语法之：?Sized

Sized 这个 trait呢，是比较特殊的，是自动实现了的。一个类型是否是Sized，要看它在编译期内的size是否已知且固定不变。

比如 u8 的大小是 1 byte，但是 [T] 和 trait则是size未知的。一个 slice的[T]是未知的，是因为在编译期，不知道你到底会有多少个T存在。一个trait的size是未知的，是因为不知道实现这个trait的结构是什么。

不过，把unsized的类型放到指针或者Box里面，它们就变成了sized了，通过指针找到源头，然后顺着源头找到其他的数据。

所有的类型参数，如`fn foo<T>(){}`中的T，默认都是实现了Sized的了（自动实现），这就限制了我们传参数的数据了，于是出现了?Sized,这是一个非常特殊的用法，其他的trait bound，都是用来缩小数据的类型范围的，这个是用来扩大类型范围的，也就是说`fn foo<T:?Sized>(){}`现在可以接受unsized的数据类型了。

本文主要是看这篇文章理解的[：The Sized Trait][1]


  [1]: http://huonw.github.io/blog/2015/01/the-sized-trait/