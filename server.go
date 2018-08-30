package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"

	"./test"

	"github.com/golang/protobuf/proto"
)

const SERVER_ADDR = "127.0.0.1:8333"

func write() {
	fmt.Println(test.PhoneType_HOME, test.PhoneType_WORK)
	p1 := &test.Person{
		Id:   1,
		Name: "小张",
		Phones: []*test.Phone{
			{Type: 0, Number: "111111111"},
			{Type: 1, Number: "222222222"},
		},
	}
	p2 := &test.Person{
		Id:   2,
		Name: "小王",
		Phones: []*test.Phone{
			{Type: test.PhoneType_HOME, Number: "333333333"},
			{Type: test.PhoneType_WORK, Number: "444444444"},
		},
	}

	//创建地址簿
	book := &test.ContactBook{}
	book.Persons = append(book.Persons, p1)
	book.Persons = append(book.Persons, p2)

	//编码数据
	data, _ := proto.Marshal(book)
	//把数据写入文件
	ioutil.WriteFile("./test.txt", data, os.ModePerm)

}

func read() {
	//读取文件数据
	data, _ := ioutil.ReadFile("./test.txt")
	book := &test.ContactBook{}
	//解码数据
	proto.Unmarshal(data, book)
	for _, v := range book.Persons {
		fmt.Println(v.Id, v.Name)
		for _, vv := range v.Phones {
			fmt.Println(vv.Type, vv.Number)
		}
	}
}

func ServerHandler(listen *net.TCPListener) {
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Printf("Accept Client Connect Exception : %s", err.Error())
			continue
		}
		log.Printf("Client From : %s", conn.RemoteAddr().String())
		go func() {
			defer conn.Close()
			data := make([]byte, 4)
			n, _ := io.ReadFull(conn, data)
			fmt.Println(data, n)
			b_buf := bytes.NewBuffer(data)
			var x int32

			binary.Read(b_buf, binary.BigEndian, &x)

			fmt.Println(x)
			data3 := make([]byte, x)
			io.ReadFull(conn, data3)
			book := &test.ContactBook{}
			proto.Unmarshal(data3, book)
			kk := 0
			for _, v := range book.Persons {
				kk++
				fmt.Println(v.Id, v.Name)
				for _, vv := range v.Phones {
					fmt.Println(vv.Type, vv.Number)
				}
			}
			fmt.Println(kk)
			//	res := string(data[:])
			peerAddr := conn.RemoteAddr().String()

			conn.Write([]byte("ok," + peerAddr))
		}()
	}
}

func main() {
	serverAddr, _ := net.ResolveTCPAddr("tcp", SERVER_ADDR)

	listen, err := net.ListenTCP("tcp", serverAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("Local Address:<%s> \n", listen.Addr().String())
	ServerHandler(listen)
}
