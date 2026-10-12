package nursery

import (
	"sync/atomic"
)

// Represents a function running concurrently.
type task struct {
	done chan struct{}
	fn   func()
}

// An opaque handle to a task within a nursery.
type Task int

// Create a new task (don't run it yet).
func newTask(fn func()) task {
	var closed atomic.Bool
	closed.Store(false)
	return task{done: make(chan struct{}), fn: fn}
}

// Start running a task.
func (t *task) run() {
	defer close(t.done)
	t.fn()
}

// Check if a task has completed.
func (t task) isFinished() bool {
	select {
	case <-t.done:
		return true
	default:
		return false
	}
}

// Wait for a task to complete.
func (t task) await() {
	<-t.done
}

// A 'container' for a bunch of concurrently running tasks; tasks running as
// part of a nursery are guaranteed to be closed by the time the nursery is
// closed.
type Nursery struct {
	tasks  []task
	closed *atomic.Bool
}

// Create a nursery in which to run concurrent tasks.
func OpenNursery() Nursery {
	var closed atomic.Bool
	closed.Store(false)

	return Nursery{
		tasks:  make([]task, 0),
		closed: &closed,
	}
}

func (n *Nursery) isClosed() bool {
	return n.closed.Load()
}

// Start a new task in the nursery concurrently, and return a handle to it.
func (n *Nursery) Spawn(fn func()) Task {
	if n.isClosed() {
		panic("attempt to spawn a task in a closed nursery")
	}

	task := newTask(fn)
	n.tasks = append(n.tasks, task)
	id := len(n.tasks) - 1
	go task.run()
	return Task(id)
}

// Wait for a specific task in the nursery to complete.
func (n *Nursery) Await(id Task) {
	if n.isClosed() {
		panic("attempt to await a task in a closed nursery")
	}

	task := n.tasks[id]
	if !task.isFinished() {
		task.await()
	}
}

// Close the nursery, ending any running tasks within it.
func (n *Nursery) Close() {
	if n.isClosed() {
		panic("attempt to re-close a closed nursery")
	}

	for id := range n.tasks {
		n.Await(Task(id))
	}
	n.closed.CompareAndSwap(false, true)
}
