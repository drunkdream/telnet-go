package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
)

func telnet(host string, port int) int {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Println(fmt.Sprintf("Connect %s:%d failed: %s", host, port, err))
		return -1
	}
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(2)
	go handleWrite(conn, &wg)
	go handleRead(conn, &wg)
	wg.Wait()

	return 0
}

func handleRead(conn net.Conn, wg *sync.WaitGroup) int {
	reader := bufio.NewReader(conn)
	buff := make([]byte, 4096)
	for {
		var bytes int
		var err error
		bytes, err = reader.Read(buff)
		if err != nil {
			fmt.Println("Error to read from connection because of ", err)
			return -1
		}
		_, err = os.Stdout.Write(buff[:bytes])
		if err != nil {
			fmt.Println("Error to send message because of ", err)
			return -1
		}
	}
}

func handleWrite(conn net.Conn, wg *sync.WaitGroup) int {
	reader := bufio.NewReader(os.Stdin)
	buff := make([]byte, 4096)
	for {
		var bytes int
		var err error
		bytes, err = reader.Read(buff)
		if err != nil {
			fmt.Println("Error to read from stdin because of ", err)
			return -1
		}
		_, err = conn.Write(buff[:bytes])
		if err != nil {
			fmt.Println("Error to send message because of ", err)
			return -1
		}
	}

}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: telnet host port")
		os.Exit(-1)
	}
	host := os.Args[1]
	port, _ := strconv.Atoi(os.Args[2])
	fmt.Println(telnet(host, port))
}
