// pi.go estimates pi using goroutines and  infinite Series Approach
//
// Usage: pi interval threads
//
//	  interval = the number of iterations to perform
//	  threads = the number of threads (i.e., goroutines to spawn)
//
//	Author:, Lamont Samuels
//	Date: 04/09/21
//	Cite: https://towardsdatascience.com/how-to-make-pi-part-1-d0b41a03111f
package main

import (
	"fmt"
	"math/rand"
	"os"
	"proj3/concurrent"
	"strconv"
)

type goContext struct {
	piFlag        int32
	intervals     int
	numOfSummands float64
	rGen          *rand.Rand
}

func calculateIntervals(intervals, start, end int) float64 {
	//some code
	return 1.0
}

type IntervalTask struct {
	ctx   *goContext
	rank  int
	start int
	end   int
}

func NewIntervalTask(ctx *goContext, rank, start, end int) concurrent.Callable {
	return &IntervalTask{ctx, rank, start, end}
}

// threadWork is the function that is called t
func (task *IntervalTask) Call() interface{} {

	//Define local interval
	localSums := calculateIntervals(task.ctx.intervals, task.start, task.end)

	return localSums

}
func main() {

	//Retrieve the command-line arguments and perform conversion if needed
	threadCount, _ := strconv.Atoi(os.Args[2])
	intervals, _ := strconv.Atoi(os.Args[1])

	// executor := concurrent.NewWorkStealingExecutor(threadCount, 10)
	executor := concurrent.NewWorkBalancingExecutor(threadCount, 10, 10)

	workAmount := intervals / threadCount
	var total, work, start, end int
	var futures []concurrent.Future
	var context goContext
	context.intervals = intervals
	var sum float64

	for i := 0; i < threadCount; i++ {
		start = total
		if i == threadCount-1 {
			work = intervals - total
		} else {
			work = workAmount
		}
		total += work
		end = total
		futures = append(futures, executor.Submit(NewIntervalTask(&context, i, start, end)))

	}
	for _, future := range futures {
		//Get back the local sums from the futures
		value := future.Get()
		localSum := value.(float64)
		sum += localSum
	}
	executor.Shutdown()

	//Print out the estimate
	piEstimate := 4.0 * sum
	fmt.Printf("%.10f\n", piEstimate)
}
