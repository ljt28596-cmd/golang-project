package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
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

// 菜单
func (this *Client) Menu() bool {
	var flag int
	fmt.Println("1.public chat")
	fmt.Println("2.private chat")
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

// 更新用户名
func (this *Client) UpdateName() bool {
	fmt.Println(">>>>>>please input new username:")
	fmt.Scanln(&this.Name)

	sendMsg := "rename|" + this.Name + "\n"
	//发送给客户端处理
	_, err := this.conn.Write([]byte(sendMsg))

	if err != nil {
		fmt.Println("conn.write error:", err)
		return false
	}

	return true
}

// 公聊(可输入空格)
func (this *Client) Public() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("please input the content, input 'exit' to quit...")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("read error:", err)
			break
		}
		content := strings.TrimSpace(line) //去掉首尾空白字符>>>相当于python中的.strip()
		if content == "exit" {
			break
		}
		if _, err := this.conn.Write([]byte(content + "\n")); err != nil {
			fmt.Println("conn.write error:", err)
			break
		}
	}
}

// 公聊模式
func (this *Client) PublicChat() {
	//用户输入消息
	fmt.Println(">>>>>now the model is public chat<<<<<<")
	fmt.Println("please input the content,input 'exit' to quit...")
	var content string
	fmt.Scanln(&content)

	//消息不为空则发送
	for content != "exit" {
		if len(content) != 0 {
			sendMsg := content + "\n"
			_, err := this.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.write error:", err)
				break
			}
		}
		content = ""
		fmt.Println("please input the content,input 'exit' to quit...")
		fmt.Scanln(&content)

	}

}

func (this *Client) SelectUsers() {
	sendMsg := "who"
	_, err := this.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}

// 私聊模式
func (this *Client) PrivateChat() {
	var content, obj string

	fmt.Println(">>>>>>now the model is private chat")
	this.SelectUsers()
	fmt.Println("please input the receiver")
	fmt.Scanln(&obj)

	fmt.Println("please input the content,<input 'exit' to quit,input 'switch' to change receiver>...")
	fmt.Scanln(&content)
	for content != "exit" {
		if content == "switch" {
			this.SelectUsers()
			fmt.Println("please input the receiver")
			fmt.Scanln(&obj)
			fmt.Println("switch receiver successfully...")
		} else {
			if len(content) != 0 {
				sendMsg := "to|" + obj + "|" + content + "\n"
				//发送给客户端处理
				_, err := this.conn.Write([]byte(sendMsg))

				if err != nil {
					fmt.Println("conn.write error:", err)
				}
			}
		}
		content = ""
		fmt.Scanln(&content)
	}

}

// 接收server发来的信息
func (this *Client) DealResponse() {
	//永久监听client.conn(也就是server发来的)
	io.Copy(os.Stdout, this.conn)
}

func (this *Client) Run() {
	for this.flag != 0 {
		for this.Menu() != true {
		}

		//根据不同模式处理不同业务
		switch this.flag {
		case 1:
			//公聊模式
			this.PublicChat()

		case 2:
			//私聊模式
			this.PrivateChat()

		case 3:
			//更新用户名
			this.UpdateName()

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

	//开启一个goroutine去处理server发来的消息
	go client.DealResponse()

	fmt.Println(">>>>>>>>>连接服务器成功")

	client.Run()
}
