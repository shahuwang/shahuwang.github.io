﻿### Rust的那些语法：trait和 Any

Rust有很多的Trait，日常编程中会需要用到这些trait，如果对他们不了解，就很难写好代码。去看Rust写得一些库，那真的是trait和泛型满天飞。

每一种编程语言，都会大量使用到某些个结构或者模式。在用这些语言写代码的时候，这些模式或者结构会比其他语言更常用到。Rust就真的是很喜欢用trait和泛型。

第一个想讲的是Any ，所有的'static生存期的类型都实现了Any，具备在运行时对数据的类型进行动态作用。实现了Any的类型，如果是 &Any, 则有 is 和 downcast_ref 两个方法来判断数据是不是某个类型的，并可以拿到这个类型。

看下面代码：

    fn log<T: Any + Debug>(value: &T) {
    let value_any = value as &Any;
       // try to convert our value to a String.  If successful, we want to
        // output the String's length as well as its value.  If not, it's a
        // different type: just print it out unadorned.
        match value_any.downcast_ref::<String>() {
            Some(as_string) => {
                println!("String ({}): {}", as_string.len(), as_string);
            }
            None => {
                println!("{:?}", value);
            }
        }
    }

说实在的，value_any.downcast_ref::() 这里真的是不忍吐槽Rust的泛型的语法形式。定义的时候是这样的：`fn downcast_ref<T>(&self> -> Option<&T> where T:Any` ， 然后T如果为String类型，就变成上面的那种调用方式了。 :: 本来是模块使用的符合，这里泛型也用这个，意思完全不明确。

由于文章是写给我自己看为主的，所以思维会很跳跃。

看了下 Any 实现的代码，用了很多trait的语法，发现我是相当不熟悉，所以我觉得这里要补一下trait的语法先。

在trait里面，都有一个大写的Self，用来指向实现这个trait的类型。

trait objects是一个比较难以理解的概念，在Rust里面， 形如 &SomeTrait 或者 Box 都叫做trait object，示例如下：

     trait Printable {
          fn stringify(&self) -> String;
        }
    
     impl Printable for i32 {
          fn stringify(&self) -> String { self.to_string() }
        }
    
        fn print(a: Box<Printable>) {
           println!("{}", a.stringify());
        }
    
        fn main() {
           print(Box::new(10) as Box<Printable>);
        }

就像Java里面，大部分时候，我们都要使用interface的实现，而不是interface本身，但是有时候用interface本身也挺好的。

trait objects的主要目的是为了 late binding，Box::new(10) as Box<Printable> 这句可以最好的说明这个需求了。我的数据可能有多种用途，我不想在声明的时候就把所有需要实现的trait都声明了，因为可能两个trait会有冲突呢。通过trait objects的方式，可以更方便。

就像上面说的，用Java的interface去理解trait objects就对了。

另外，在trait objects 里面，是可以通过Self取得实现trait的真正类型。

类型参数（type parameters),即trait里面的方法需要用到的泛型类型参数，需要用如下的语法：

    trait Seq<T> {
       fn len(&self) -> u32;
       fn elt_at(&self, n: u32) -> T;
       fn iter<F>(&self, F) where F: Fn(T);
    }

这里的就是类型参数，是给里面的方法使用的，和函数的泛型类型参数类似。

还有一个比较复杂的，应该是 Associated Types 了，说白了，这是对泛型的类型参数泛滥的一种优化，使得不需要声明那么多的类型参数，或者根本就不需要去声明了。链接中已经解释得比较清楚了，通过将：trait Graph<N, E> {} 修改为 trait Graph{type N; type E;},调用方式则从 fn distance<N, E, G:Graph<N, E>> 简化为 fn distance<G:Graph>。

这是我觉得Rust做得很不好的地方，很多语法并不正交，原则上第一种方法调用就不应该存在。

之前我一直没有意识到，原来也可以 impl trait for trait的，比如Any实现的源码里面：impl fmt::Debug for Any

转念想想，trait object和struct在类型意义上也没有多大不同，所有 impl trait for trait，也是挺合理的



