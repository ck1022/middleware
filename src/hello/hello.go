package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 100; i++ {
		printi(i)
	}
	time.Sleep(time.Second)
	fmt.Println(time.Second)
}
func printi(i int) {
	fmt.Println(i)
	time.Sleep(10000)
}
