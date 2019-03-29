// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	//"os"
	"time"
)

//!+broadcaster
type client chan<- string // an outgoing message channel


var users = make(map[string]string);

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)


func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg
			}
			
			

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

//!-broadcaster


//!+handleConn
func handleConn(conn net.Conn, user string) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)
	who := conn.RemoteAddr().String()
	ch <- "You are " + user
	messages <- user + " has arrived"
	entering <- ch


	users[user] = who

	

	input := bufio.NewScanner(conn)
	fmt.Println(input.Text());
	for input.Scan() {
		var s []string
		s = strings.Fields(input.Text())
		if(len(s) != 0){
		switch s[0] {
		case "/users":
			var u string
			for k, _ :=range users{
				u = u + k + ", "
			}
			ch <- u
		case "/user":
			if(len(s) == 2){
			var u string
			u = "not user found"
			for k, v :=range users{
				if(strings.Trim(k, " ") == strings.Trim(s[1], " ")){
					u = k + " " + v
				}
			}
			ch <- u
			}else{
				ch <- "not enough arguments /user <username>"
			}
			
		case "/time":
			ch <- time.Now().Format("15:04:05\n")
		case "/msg":
			var u string

			ch <- u
		default:
			messages <- user + ": " + input.Text()
		}
		}
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- ch
	messages <- user + " has left"
	delete(users, user)
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

//!+main
func main() {

	
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		message,_:= bufio.NewReader(conn).ReadString(' ');
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn, string(message));
	}
}

//!-main
