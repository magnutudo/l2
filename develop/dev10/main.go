package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	i, _ := strconv.Atoi(os.Args[3])
	duration := time.Duration(i) * time.Second
	conn, err := net.DialTimeout("tcp", os.Args[1]+os.Args[2], duration)
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")

	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 256)
	for {
		conn.SetReadDeadline(time.Now().Add(time.Second * 5))
		if err != nil {
			return
		}
		n, err := conn.Read(tmp)

		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}

		buf = append(buf, tmp[:n]...)

	}
	fmt.Println(string(buf))
}
