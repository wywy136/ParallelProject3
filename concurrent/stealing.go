package concurrent

import (
	"math/rand"
	"sync"
)

// NewWorkStealingExecutor returns an ExecutorService that is implemented using the work-stealing algorithm.
// @param capacity - The number of goroutines in the pool
// @param threshold - The number of items that a goroutine in the pool can
// grab from the executor in one time period. For example, if threshold = 10
// this means that a goroutine can grab 10 items from the executor all at
// once to place into their local queue before grabbing more items. It's
// not required that you use this parameter in your implementation.
func NewWorkStealingExecutor(capacity, threshold int) ExecutorService {
	// Initiate the stealing executor
	stealingExecutor := executor{
		capacity: capacity,
		numTasks: 0,
		wg:       &sync.WaitGroup{},
	}
	stealingExecutor.workers = make([]*worker, capacity)

	// Initiate workers
	for i := 0; i < capacity; i++ {
		stealingExecutor.workers[i] = &worker{
			localQueue: NewUnBoundedDEQueue(),
			id:         i,
		}
	}

	// Manage peers for each worker
	for i := 0; i < capacity; i++ {
		peers := make([]*worker, capacity)
		copy(peers, stealingExecutor.workers)
		stealingExecutor.workers[i].peers = append(peers[:i], peers[i+1:]...)
	}

	// Start all the workers in new threads
	for i := 0; i < capacity; i++ {
		go stealingExecutor.workers[i].runWorkStealing()
	}

	return &stealingExecutor
}

func (w *worker) runWorkStealing() {
	for true {
		// Steal a work
		if w.localQueue.IsEmpty() {
			// Randomly pick up a victim from the worker's peers to steal from
			victim := w.peers[rand.Intn(len(w.peers))]
			// No more work, spin
			if victim.localQueue.IsEmpty() {
				continue
			} else {
				// Execute the stolen task
				task := victim.localQueue.PopBottom()
				if task != nil {
					executeFutureTask(task)
				}
			}
		} else { // Consume the task in the local queue
			task := w.localQueue.PopTop()
			if task != nil {
				executeFutureTask(task)
			}
		}
	}
}
