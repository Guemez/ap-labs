package main

import (
	"golang.org/x/tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {
	count := make(map[string]int)
	new_s := strings.Fields(s)
	for _, i := range new_s{
		if _, ok := count[i]; ok{
			count[i] = count[i] + 1
		}else{
			count[i] = 1
		}
	}
	return count
}

func main() {
	wc.Test(WordCount)
}