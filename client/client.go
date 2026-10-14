package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	//连接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn
	return client

}

func (this *Client) Menu() bool {
	var flag int
	fmt.Println("1.public chat")
	fmt.Println("1.private chat")
	fmt.Println("3.update username")
	fmt.Println("0.exit")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		this.flag = flag
		return true
	} else {
		fmt.Println(">>>>>>please input validate figure<<<<<<<")
		return false
	}
}

func (this *Client) Run() {
	for this.flag != 0 {
		for this.Menu() != true {
		}

		//根据不同模式处理不同业务
		switch this.flag {
		case 1:
			//公聊模式
			fmt.Println("public chat...")
			break
		case 2:
			//私聊模式
			fmt.Println("private chat...")
			break
		case 3:
			//更新用户名
			fmt.Println("update username...")
			break
		}
	}

}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认是8888)")
}

func main() {
	//命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>连接服务器失败...")
		return
	}
	fmt.Println(">>>>>>>>>连接服务器成功")

	client.Run()
}
