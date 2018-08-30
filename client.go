package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

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
		Name: "老傻逼",
		Phones: []*test.Phone{
			{Type: test.PhoneType_HOME, Number: "333333333"},
			{Type: test.PhoneType_WORK, Number: "444444444"},
		},
	}

	//创建地址簿
	book := &test.ContactBook{}
	book.Persons = append(book.Persons, p1)
	book.Persons = append(book.Persons, p2)
	for iii := 0; iii < 10; iii++ {
		book.Persons = append(book.Persons, p2)
	}

	//编码数据
	data, _ := proto.Marshal(book)
	//把数据写入文件

	serverAddr, _ := net.ResolveTCPAddr("tcp", SERVER_ADDR)
	log.Printf("SERVER IP : %s", serverAddr.String())

	conn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		log.Printf("SERVER IP : %s", serverAddr.String())
		log.Printf("Connect Error : %s", err.Error())
		return
	}
	fmt.Println("Connect Success!")
	defer conn.Close()
	size := len(data)
	fmt.Println(size)

	b_buf := bytes.NewBuffer([]byte{})

	err = binary.Write(b_buf, binary.BigEndian, int32(size))
	b_buf.Write(data)
	//	fmt.Println(b_buf.Bytes())

	conn.Write(b_buf.Bytes())
	buf := make([]byte, 1024)
	c, err := conn.Read(buf)
	if err != nil {
		log.Printf("Read Server Data Exception: %s", err.Error())
		return
	}
	//	ioutil.WriteFile("./test.txt", data, os.ModePerm)
	fmt.Println(c, string(buf))

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

func main() {
	for i := 0; i < 10000; i++ {
		go func() {
			write()
		}()
	}
	time.Sleep(time.Second * 199)
	//	read()
}
