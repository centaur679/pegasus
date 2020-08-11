package main

import (
	"fmt"
	"log"
	"net"

	"github.com/centaur679/pegasus/server/utils"

	pb "github.com/centaur679/pegasus/server/api"
	proto "github.com/golang/protobuf/proto"
)

const (
	confPath = "server.yml"
)

func main() {
	// 解析配置文件
	var c utils.Conf
	conf := c.GetConf(confPath)
	defer fmt.Println("程序即将退出！")
	if !conf.Validate() {
		log.Panicln("配置校验错误")
	}
	fmt.Println(conf)

	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: conf.Port})
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("本地地址: <%s> \n", listener.LocalAddr().String())

	peers := []net.UDPAddr{}
	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			log.Println("接收消息出错：", err)
			continue
		}

		req := &pb.TestIPRequest{}
		proto.Unmarshal(data[:n], req)
		peers = append(peers, *remoteAddr)

		res := &pb.TestIPResponse{ReceiveIP: remoteAddr.IP.To4().String()}
		out, err := proto.Marshal(res)
		if err != nil {
			log.Println(err)
			continue
		}
		listener.WriteToUDP(out, remoteAddr)
	}

	// peers := make([]net.UDPAddr, 0, 2)
	// data := make([]byte, 1024)
	// for {
	// 	n, remoteAddr, err := listener.ReadFromUDP(data)
	// 	if err != nil {
	// 		fmt.Printf("error during read: %s", err)
	// 	}
	// 	log.Printf("<%s> %s\n", remoteAddr.String(), data[:n])
	// 	peers = append(peers, *remoteAddr)
	// 	if len(peers) == 2 {
	// 		log.Printf("进行UDP打洞,建立 %s <--> %s 的连接\n", peers[0].String(), peers[1].String())
	// 		listener.WriteToUDP([]byte(peers[1].String()), &peers[0])
	// 		listener.WriteToUDP([]byte(peers[0].String()), &peers[1])
	// 		time.Sleep(time.Second * 8)
	// 		log.Println("中转服务器退出,仍不影响peers间通信")
	// 		return
	// 	}
	// }
}
