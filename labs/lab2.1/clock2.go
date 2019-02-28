// Clock2 is a concurrent TCP server that periodically writes the time.
package main

import (
	"io"
	"log"
	"net"
	"time"
	"os"
)

func handleConn(c net.Conn) {
	defer c.Close()
	timeZone := os.Getenv("TZ")

	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Fatal(err)
	}

	for {
		_, err := io.WriteString(c, (timeZone + ":  " + time.Now().In(loc).Format("15:04:05\n")))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	var port string
		if os.Args[1] == "-port" {
			port = os.Args[2]
	}
	listener, err := net.Listen("tcp", "localhost:" + port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn) // handle connections concurrently
	}
}
