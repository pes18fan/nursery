package main

import (
	"fmt"
	"time"

	"github.com/pes18fan/nursery"
)

func do_some_work() {
	// Create a nursery, and end it when the function is done.
	n := nursery.OpenNursery()
	defer n.Close()

	// Spawn a couple of tasks.
	n.Spawn(func() {
		fmt.Println("doing some stuff")
	})

	id := n.Spawn(func() {
		fmt.Println("doing some other stuff")
		time.Sleep(time.Second)
	})

	// Stop here until the second one is done.
	n.Await(id)

	// Spawn another!
	n.Spawn(func() {
		fmt.Println("doing some other other stuff after the other stuff")
	})
}

func main() {
	fmt.Println("about to start concurrent work")

	// Do concurrent tasks within here, exiting when all are done
	do_some_work()

	// We will only get here once all the concurrent tasks in the nursery are
	// done!
	fmt.Println("all work done!")
}
