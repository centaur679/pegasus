package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/centaur679/pegasus/client/utils"

	pb "github.com/centaur679/pegasus/client/api"
	proto "github.com/golang/protobuf/proto"
)

var tag string

const (
	handMsg  = "我是打洞消息"
	confPath = "client.yml"
)

func main() {
	// 解析配置文件
	var c utils.Conf
	conf := c.GetConf(confPath)
	defer fmt.Println("程序即将退出！1")
	defer fmt.Println("程序即将退出！2")
	if !conf.Validate() {
		log.Panicln("配置校验错误")
	}
	fmt.Println(conf)

	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: conf.Port}
	dstAddr := &net.UDPAddr{IP: net.ParseIP(conf.Server.IP), Port: conf.Server.Port}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		log.Panicln(err)
	}

	req := &pb.TestIPRequest{Name: conf.Name}
	out, err := proto.Marshal(req)
	if err != nil {
		log.Panicln(err)
	}
	if _, err = conn.Write(out); err != nil {
		log.Panicln(err)
	}

	data := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		log.Panicln(err)
	}

	if remoteAddr.IP.To4().String() != dstAddr.IP.To4().String() {
		log.Panicln("源IP不符合")
	}

	conn.Close()
	res := &pb.TestIPResponse{}
	proto.Unmarshal(data[:n], res)
	fmt.Println("res.ReceiveIP: ", res.ReceiveIP)
	hasNat := testNAT(res.ReceiveIP)
	if hasNat {
		// 打洞
		fmt.Println("有NAT")
	} else {
		// 直接通信
		fmt.Println("没有NAT")
	}

	// bidirectionHole(srcAddr, &anotherPeer)
}

// 服务器返回 sendip:receiveip,通过比较两个IP 地址来决定是否打洞
func testNAT(addr string) bool {
	ips, err := utils.GetAvailableIPAddress()
	if err != nil {
		log.Panicln(err)
	}
	return ips[0] != addr
}

func bidirectionHole(srcAddr *net.UDPAddr, anotherAddr *net.UDPAddr) {
	conn, err := net.DialUDP("udp", srcAddr, anotherAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	// 向另一个peer发送一条udp消息(对方peer的nat设备会丢弃该消息,非法来源),用意是在自身的nat设备打开一条可进入的通道,这样对方peer就可以发过来udp消息
	if _, err = conn.Write([]byte(handMsg)); err != nil {
		log.Println("send handshake:", err)
	}
	go func() {
		for {
			time.Sleep(10 * time.Second)
			if _, err = conn.Write([]byte("from [" + tag + "]")); err != nil {
				log.Println("send msg fail", err)
			}
		}
	}()
	for {
		data := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Printf("error during read: %s\n", err)
		} else {
			log.Printf("收到数据:%s\n", data[:n])
		}
	}
}
