#编程桥段： 一个简单的TCP连接多路复用

这是我做的一个练习，目的是学习如何实现在一个TCP连接里，发送多个来源的数据，并保证不混淆。

因为是个简单的练习，因为 Go 语言还不甚熟悉，所以很多东西是没有去考虑的，而且我也不想弄得那么复杂，少少的上百行代码就好了，纯粹练习其核心思想。

这是我在学习 Gotunnel 的代码时，看到其实现，而模仿的，模仿对于我们这种笨人，也是有很大好处的。

代码分两部分，client 端和 server 端， client 端建立一条 TCP 连接到 server 端，然后用两个 telnet 连接 client 端， 各自发送消息给 client 端。 client 端将消息通过该 TCP 连接发送给 server 端，server 端只是简单的把消息原样返回, client 端将收到的 server 端消息返回给发送此消息的 telnet 客户端。两个 telnet 发送什么就收到什么，而没有收到别人，或者未收到回复，即表明多路复用成功。

运用的思想很简单， 每个和 client 端建立的连接，都分配一个 id， 然后发送的每份数据，数据头部包含了此数据属于哪个连接，以及此数据的长度。server 端只是原样返回信息，不需要过多说明。client 端收到回复之后，从收到的数据里，读取出数据头部信息，即连接的 id 和数据的长度。然后根据 id 从字典里面找到此连接，将数据发送回去。

这其实也是一种很简单的数据协议设计，当然， Gotunnel 采取这种做法的主要原因是建立加密信道。

代码如下：

    client 端
    
    package main
    
    import (
    	"bufio"
    	"encoding/binary"
    	"fmt"
    	"io"
    	"net"
    	"time"
    )
    
    var connMap map[uint32]*net.TCPConn
    
    func main() {
    	addr, err := net.ResolveTCPAddr("tcp", "localhost:8088")
    	if err != nil {
    		panic("addr invalid")
    	}
    	ln, err := net.ListenTCP("tcp", addr)
    	if err != nil {
    		panic(err)
    	}
    	tunnel, err := net.DialTimeout("tcp", "localhost:8080", time.Duration(30)*time.Second)
    	if err != nil {
    		fmt.Errorf(err.Error())
    		return
    	}
    	var i uint32 = 1
    	connMap = make(map[uint32]*net.TCPConn)
    	go receive(tunnel.(*net.TCPConn))
    	for {
    		fmt.Println("recieve from client")
    		conn, err := ln.AcceptTCP()
    		if err != nil {
    			fmt.Errorf(err.Error())
    			continue
    		}
    		go handle(conn, tunnel.(*net.TCPConn), i)
    		connMap[i] = conn
    		i = i + 1
    	}
    }
    
    func receive(tunnel *net.TCPConn) {
    	for {
    		fmt.Println("read from server")
    		reader := bufio.NewReader(tunnel)
    		var h header
    		binary.Read(reader, binary.LittleEndian, &h)
    		buf := make([]byte, h.Len)
    		io.ReadFull(reader, buf)
    		conn := connMap[h.Linkid]
    		fmt.Println(h.Linkid)
    		conn.Write(buf)
    	}
    }
    
    type header struct {
    	Linkid uint32
    	Len    uint32
    }
    
    func handle(conn *net.TCPConn, tunnel *net.TCPConn, i uint32) {
    	fmt.Println("handle connect")
    	writer := bufio.NewWriter(tunnel)
    	buf := make([]byte, 1024)
    	for {
    		n, err := conn.Read(buf)
    		if err != nil {
    			fmt.Errorf(err.Error())
    			return
    		}
    		binary.Write(writer, binary.LittleEndian, &header{i, uint32(n)})
    		writer.Write(buf[:n])
    		writer.Flush()
    		fmt.Println("write to server")
    	}
    }


server 端

    package main
    
    import (
    	"fmt"
    	"net"
    )
    
    func main() {
    	addr, err := net.ResolveTCPAddr("tcp", "localhost:8080")
    	if err != nil {
    		panic("add invalid")
    	}
    	ln, err := net.ListenTCP("tcp", addr)
    	if err != nil {
    		panic(err)
    	}
    	fmt.Println("server listening 8080")
    	for {
    		conn, _ := ln.AcceptTCP()
    		buf := make([]byte, 1024)
    		for {
    			fmt.Println("read from conn")
    			n, err := conn.Read(buf)
    			if err != nil {
    				fmt.Errorf(err.Error())
    				break
    			}
    			conn.Write(buf[:n])
    		}
    	}
    }

