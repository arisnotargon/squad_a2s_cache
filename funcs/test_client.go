package funcs

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"net"
	"sync"

	"github.com/davecgh/go-spew/spew"
)

var (
	REQUEST_DATA_PREFIX *[]byte
	challengeStore      sync.Map
	once                sync.Once
)

func init() {
	once.Do(func() {
		REQUEST_DATA_PREFIX = new([]byte)
		*REQUEST_DATA_PREFIX = make([]byte, 0)
		*REQUEST_DATA_PREFIX = append(*REQUEST_DATA_PREFIX, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x54}...)
		*REQUEST_DATA_PREFIX = append(*REQUEST_DATA_PREFIX, []byte("Source Engine Query")...)
		*REQUEST_DATA_PREFIX = append(*REQUEST_DATA_PREFIX, byte(0x00))
	})
}

func SendTestPackage() {
	spew.Dump(*REQUEST_DATA_PREFIX)
	// 启动测试服务器
	// serverOk := make(chan struct{})
	// go testServer(serverOk)
	// <-serverOk
	testClient()
}

func StartTestServerInfo() {
	spew.Dump(*REQUEST_DATA_PREFIX)
	// 启动测试服务器
	serverOk := make(chan struct{})
	go testServer(serverOk)
	<-serverOk
	fmt.Println("测试服务器已启动")

	select {}
}

func handleReq(listen *net.UDPConn, data []byte, addr *net.UDPAddr) error {
	responseData := make([]byte, 0)
	if bytes.Equal(data, *REQUEST_DATA_PREFIX) {
		challenge := make([]byte, 4)
		_, err := rand.Read(challenge)
		if err != nil {
			fmt.Printf("rand Read failed, err: %+v", err)
			return err
		}
		responseData = append(responseData, *REQUEST_DATA_PREFIX...)
		responseData = append(responseData, challenge...)
		challengeStore.Store(fmt.Sprintf("%s-%d", addr.IP, addr.Port), challenge)
		_, err = listen.WriteToUDP(responseData, addr)
		if err != nil {
			fmt.Println("Write to udp failed, err: ", err)
			return err
		}
	} else {
		// check challenge
		challenge, has := challengeStore.LoadAndDelete(fmt.Sprintf("%s-%d", addr.IP, addr.Port))
		if !has {
			return fmt.Errorf("no challenge key : %s", fmt.Sprintf("%s-%d", addr.IP, addr.Port))
		}
		dataLen := len(data)

		if dataLen < 4 {
			return fmt.Errorf("data len error : %b", data)
		}

		if !bytes.Equal(challenge.([]byte), data[dataLen-4:]) {
			return fmt.Errorf("challenge not matched,challenge  : %v, data: %b", challenge, data)
		}
		responseData = append(responseData, []byte("Squad test server info")...)
		_, err := listen.WriteToUDP(responseData, addr)
		if err != nil {
			fmt.Println("Write server info to udp failed, err: ", err)
			return err
		}
	}
	return nil
}

func testServer(ok chan struct{}) {
	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 10002,
	})
	if err != nil {
		fmt.Println("Listen failed, err: ", err)
		return
	}
	defer listen.Close()
	ok <- struct{}{}
	for {
		var data [1024]byte
		n, addr, err := listen.ReadFromUDP(data[:])
		spew.Dump("n, addr, err", n, data[:n])
		// return

		if err != nil {
			fmt.Println("read udp failed, err: ", err)
			continue
		}
		go func() {
			err = handleReq(listen, data[:n], addr)
			if err != nil {
				fmt.Println("handleReq udp failed, err: ", err)
			}
		}()
	}
}

func testClient() {
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 3, 13),
		Port: 10002,
	})
	if err != nil {
		fmt.Printf("testClient连接UDP服务器失败 err: %+v", err)
		return
	}
	defer socket.Close()
	_, err = socket.Write(*REQUEST_DATA_PREFIX)
	if err != nil {
		fmt.Printf("首次发送数据失败,err: %+v", err)
		return
	}
	data := make([]byte, 1024)
	n, remoteAddr, err := socket.ReadFromUDP(data) // 接收数据
	if err != nil {
		fmt.Println("首次接收数据失败, err: ", err)
		return
	}
	// fmt.Printf("testClient ==> recv:%v addr:%v count:%v\n", string(data[:n]), remoteAddr, n)
	spew.Dump("testClient ==> recv:", data[:n], remoteAddr, n)

	if n < 4 {
		fmt.Println("首次接收数据失败, err: n < 4")
		return
	}

	_, err = socket.Write(append(*REQUEST_DATA_PREFIX, data[n-4:n]...))
	if err != nil {
		fmt.Printf("二次发送数据失败,err: %+v", err)
		return
	}
	n, remoteAddr, err = socket.ReadFromUDP(data) // 接收数据
	if err != nil {
		fmt.Println("二次接收数据失败, err: ", err)
		return
	}
	spew.Dump("testClient ==> recv server info:", data[:n], remoteAddr)
}
