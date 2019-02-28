package main

import (
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func handler(conn net.Conn, city string, channel chan int) {
	for true {
		_, err := io.CopyN(os.Stdout, conn, 9)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}
	channel <- 2
}

func main() {
	done := make(chan int)

	for _, args := range os.Args[1:]{
		params := strings.Split(args, "=")
		city := params[0]
		port := params[1]
		conn, err := net.Dial("tcp", port)
		if err != nil {
			log.Fatal(err)
		}
		go handler(conn, city, done)
	}

	_ = <-done
	close(done)
}
