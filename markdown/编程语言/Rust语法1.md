# Rust 那些奇奇怪怪的语法之 Cow

我认真学习过的语言里，语法最复杂的估计就是 Python 和 Java 了，这导致我对复杂一些的语言，有恐惧感。另一方面呢，因为我对 C/C++很不熟悉，所以对系统级编程有学习的想法。所以， Rust 是最符合我这两个学习目标的语言了。

不过，Rust 的语法真的好复杂呀， 今晚就编译一段小程序，一直没有编译成功。弄到最后，各种胡乱尝试，终于尝试对了。

其实我学习 Rust 已经有一个多月了，但是，至今还是完全感觉不到自己入门了。一个是语法过于复杂，一个是自己学习断断续续的，一个是自己学习方法不对。

想想还是回到我的老本行，以前我学习一个东西，都是在学习的过程中写博客，通过写博客，令自己对知识的掌握更稳固。

今天要讲的是 [std::borrow::Cow][1], Cow 的作用是封装borrow来的数据，提供immutable的数据访问，以及在需要Mutation的时候，复制这份数据，提供惰性访问。

我照着文档，抄了这么一段小程序：

    use std::borrow::Cow;

    fn abs_all(input: &mut Cow<[i32]>){
        for i in 0..input.len(){
            let v = input[i];
            if v < 0{
                input.to_mut()[i] = -v;
            }
        }
    }
    
    fn main(){
        let v:Vec<i32> = vec![1,2,3, -10];
        let mut v1:Cow<[i32]> =Cow::Owned(v);
        abs_all(&mut v1);
        println!("{:?}", v1);
    }

就是这两句把我坑了半天，一直编译不成功，一直报类型错误：

    let mut v1:Cow<[i32]> =Cow::Owned(v);
    abs_all(&mut v1);
    
实际上， abs_all 的参数类型是 &mut Cow<[i32]>, v1 明明已经是 mut 的了， 为什么还要再一次声明为 &mut 呢？

我觉得会造成我的困扰，主要是关键词做得不好，如果弄成 mut& 就好懂多了，即可改写的reference。这个问题，我觉得 [Rust by Example][2] 解释得还是挺好的。首先，一个变量要能可变，那么第一步无论如何它都必须声明为 mut。 然后，& 实际上是 borrow，borrow 当然也要分可改动和不能改动，所以要使用 &mut 。

Cow 的定义，则再次让我陷入到 Rust 语法的迷雾中了：

    #[stable(feature = "rust1", since = "1.0.0")]
    pub enum Cow<'a, B: ?Sized + 'a> where B: ToOwned {
        /// Borrowed data.
        #[stable(feature = "rust1", since = "1.0.0")]
        Borrowed(&'a B),
    
        /// Owned data.
        #[stable(feature = "rust1", since = "1.0.0")]
        Owned(<B as ToOwned>::Owned)
    }

这个定义涉及到 lifetime， 泛型， 以及 trait，组合得还是相当复杂的。

首先， where B: ToOwned 表明类型 B 必须实现了 ToOwned 。

Lifetime 的作用是在编译期就把非法访问给杜绝了，对Borrow 过来的数据的操作，必须在 owner 的 scope 之内。

最让我头疼的是，`Cow<'a, B: ?Sized + 'a>` 这样的语法到底是什么个意思。

[Explicit lifetime][3] 我感觉可以解决这个问题的一部分了。

最基本的情况，在一个 enum 里声明一个带 Lifetime 的泛型，形式是这样的： `Cow<'a, B>` 。 `B: ?Sized` 说的是 B 可能实现了 Sized。所以现在我们可以理解 `Cow<'a, B:?Sized>` 了，因为 Cow 这个 enum 有 Borrow， 所以必须指定 Lifetime, 而且还用到了泛型，所以这两个都必须在尖括号里面声明。另外，这个泛型可能实现了 Sized。+ 'a 是因为你可能在 trait 的实现里包含了其他资源，譬如说另外一个对象的 reference，所以必须保证你的这个reference的Lifetime也是 'a, 也就是说，当前Cow的实现里的Lifetime 是 'a, B实现了的 trait 里面的资源的生命周期也必须是 'a 的。





  [1]: https://doc.rust-lang.org/std/borrow/enum.Cow.html
  [2]: http://rustbyexample.com/scope/borrow/mut.html
  [3]: http://stackoverflow.com/questions/27278401/explicit-lifetime-error-in-rust