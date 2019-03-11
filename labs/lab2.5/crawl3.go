// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 241.

// Crawl2 crawls web links starting with the command-line arguments.
//
// This version uses a buffered channel as a counting semaphore
// to limit the number of concurrent calls to links.Extract.
package main

import (
	"fmt"
	"log"
	"os"
	//"strconv"
	"flag"

	"gopl.io/ch5/links"
)

// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

func crawl(url string) []string {
	fmt.Println(url)
	tokens <- struct{}{} // acquire a token
	list, err := links.Extract(url)
	<-tokens // release the token

	if err != nil {
		log.Print(err)
	}
	return list
}

var seen = make(map[string]bool)

func crawler(depth int, url string, done chan bool) {

	if depth <= 0 {
		done <- true
		return
	}

	if visited, ok := seen[url]; visited && ok {
		done <- true
		return
	} else {
		seen[url] = true
	}

	links := crawl(url)
	linksDone := make(chan bool)

	for _, link := range links {
		go crawler(depth-1, link, linksDone)
		<-linksDone
	}
	done <- true

}

func main() {

	if len(os.Args) <= 2 {
		log.Fatal("We need 2 parameters, url and depth")
	}

	url := os.Args[2];
	depth := flag.Int("depth", 1, "");
	flag.Parse();
	//fmt.Println(*depth);
	//fmt.Println(url);

	done := make(chan bool)
	seen[url] = false

	go crawler(*depth, url, done)

	<-done
}

//!-
