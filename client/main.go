package main

import (
	"fmt"
	"log"
	"net"
	"pegasus"
	"strings"
	"time"
)

var tag string

const (
	HAND_SHAKE_MSG = "我是打洞消息"
	CONF_PATH      = "conf.yml"
)

func main() {
	// 解析配置文件
	var c pegasus.Conf
	conf := c.GetConf(CONF_PATH)
	defer fmt.Println("程序即将退出！")
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

	if _, err = conn.Write([]byte("hello, I'm new peer:" + tag)); err != nil {
		log.Panicln(err)
	}

	data := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(data)
	if err != nil {
		log.Panicln(err)
	}

	conn.Close()
	hasNat := testNAT(string(data[:n]))
	if hasNat {
		// 打洞
	} else {
		// 直接通信
	}

	// bidirectionHole(srcAddr, &anotherPeer)
}

// 服务器返回 sendip:receiveip,通过比较两个IP 地址来决定是否打洞
func testNAT(addr string) bool {
	t := strings.Split(addr, ":")
	return t[0] != t[1]
}

func bidirectionHole(srcAddr *net.UDPAddr, anotherAddr *net.UDPAddr) {
	conn, err := net.DialUDP("udp", srcAddr, anotherAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	// 向另一个peer发送一条udp消息(对方peer的nat设备会丢弃该消息,非法来源),用意是在自身的nat设备打开一条可进入的通道,这样对方peer就可以发过来udp消息
	if _, err = conn.Write([]byte(HAND_SHAKE_MSG)); err != nil {
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
