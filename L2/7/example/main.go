package main

import (
	f "fanin"
	"fmt"
)

func main() {
	var chan1 chan interface{} = make(chan interface{})
	var chan2 chan interface{} = make(chan interface{})
	var chan3 chan interface{} = make(chan interface{})
	go TestWrite(chan1)
	go TestWrite(chan2)
	go TestWrite(chan3)
	fin := f.FanIn(chan1, chan2, chan3)
	for v := range fin {
		fmt.Println(v)
	}
}
func TestWrite(c chan interface{}) {
	c <- 1
	c <- true
	c <- "22"
	close(c)
}
