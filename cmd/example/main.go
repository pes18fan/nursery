package main

import (
	"fmt"
	"time"

	"github.com/pes18fan/nursery"
)

func main() {
	n := nursery.OpenNursery()
	defer n.Close()

	n.Spawn(func() {
		fmt.Println("doing some stuff")
		panic("oh noes")
	})

	id := n.Spawn(func() {
		fmt.Println("doing some other stuff")
		time.Sleep(time.Second)
	})

	n.Await(id)

	n.Spawn(func() {
		fmt.Println("doing some other other stuff after the other stuff")
	})
}
