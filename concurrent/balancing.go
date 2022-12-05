package concurrent

import (
	"math/rand"
	"sync"
)

// NewWorkBalancingExecutor returns an ExecutorService that is implemented using the work-balancing algorithm.
// @param capacity - The number of goroutines in the pool
// @param threshold - The number of items that a goroutine in the pool can
// grab from the executor in one time period. For example, if threshold = 10
// this means that a goroutine can grab 10 items from the executor all at
// once to place into their local queue before grabbing more items. It's
// not required that you use this parameter in your implementation.
// @param thresholdBalance - The threshold used to know when to perform
// balancing. Remember, if two local queues are to be balanced the
// difference in the sizes of the queues must be greater than or equal to
// thresholdBalance. You must use this parameter in your implementation.
func NewWorkBalancingExecutor(capacity, thresholdQueue, thresholdBalance int) ExecutorService {
	// Initiate the stealing executor
	balancingExecutor := executor{
		capacity:  capacity,
		numTasks:  0,
		wg:        &sync.WaitGroup{},
		threshold: thresholdBalance,
	}
	balancingExecutor.workers = make([]*worker, capacity)

	// Initiate workers
	for i := 0; i < capacity; i++ {
		balancingExecutor.workers[i] = &worker{
			localQueue: NewUnBoundedDEQueue(),
			id:         i,
		}
	}

	// Manage peers for each worker
	for i := 0; i < capacity; i++ {
		peers := make([]*worker, capacity)
		copy(peers, balancingExecutor.workers)
		balancingExecutor.workers[i].peers = append(peers[:i], peers[i+1:]...)
	}

	// Start all the workers in new threads
	for i := 0; i < capacity; i++ {
		go balancingExecutor.workers[i].runWorkBalancing(balancingExecutor.threshold)
	}

	return &balancingExecutor
}

func (w *worker) runWorkBalancing(threshold int) {
	for true {
		// Run local task
		task := w.localQueue.PopTop()
		if task != nil {
			executeFutureTask(task)
		}
		size := w.localQueue.Size()
		// Rebalance
		if rand.Intn(size+1) == size {
			// Pickup a victim randomly
			victim := w.peers[rand.Intn(len(w.peers))]
			// Compare the id of the workers to lock them in order
			if victim.id <= w.id {
				balancing(victim.localQueue, w.localQueue, threshold)
			} else {
				balancing(w.localQueue, victim.localQueue, threshold)
			}
		}
	}
}
