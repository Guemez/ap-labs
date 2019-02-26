package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func handler(conn net.Conn, city string, channel chan int) {
	var err error
	for true {
		time.Sleep((1 * time.Second))
		fmt.Printf(city + ": ")
		_, err = io.CopyN(os.Stdout, conn, 9)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}
	channel <- 2
}

func main() {
	wt := make(chan int)

	for _, args := range os.Args[1:]{
		params := strings.Split(args, "=")
		city := params[0]
		port := params[1]
		conn, err := net.Dial("tcp", port)
		if err != nil {
			log.Fatal(err)
		}
		go handler(conn, city, wt)
	}

	_ = <-wt
	close(wt)
}
