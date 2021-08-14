package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

func telnet(host string, port int) int {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 10*time.Second)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("[FAIL] Connect %s:%d failed: %s", host, port, err))
		return -1
	}
	defer conn.Close()
	fmt.Fprintln(os.Stderr, "[OKAY]")

	var wg sync.WaitGroup
	wg.Add(2)
	go handleWrite(conn, &wg)
	go handleRead(conn, &wg)
	wg.Wait()

	return 0
}

func handleRead(conn net.Conn, wg *sync.WaitGroup) int {
	defer wg.Add(-2)
	reader := bufio.NewReader(conn)
	buff := make([]byte, 4096)
	for {
		var bytes int
		var err error
		bytes, err = reader.Read(buff)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error to read from upstream because of", err)
			return -1
		}

		_, err = os.Stdout.Write(buff[:bytes])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error to write to stdout because of", err)
			return -1
		}
	}
}

func handleWrite(conn net.Conn, wg *sync.WaitGroup) int {
	defer wg.Add(-2)
	reader := bufio.NewReader(os.Stdin)
	buff := make([]byte, 4096)
	for {
		var bytes int
		var err error
		bytes, err = reader.Read(buff)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error to read from stdin because of", err)
			return -1
		}

		_, err = conn.Write(buff[:bytes])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error to write to upstream because of", err)
			return -1
		}
	}

}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: telnet host port")
		os.Exit(-1)
	}
	host := os.Args[1]
	port, _ := strconv.Atoi(os.Args[2])
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for sig := range c {
			// sig is a ^C, handle it
			fmt.Fprintln(os.Stderr, "Signal", sig)
			os.Exit(0)
		}
	}()
	os.Exit(telnet(host, port))
}
