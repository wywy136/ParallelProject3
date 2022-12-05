package main

import (
	"fmt"
	"proj3/concurrent"
	"sync/atomic"
)

type ZipcodeInfo struct {
	cases  int // Total number of cases for the month
	tests  int // Total number of tests for the month
	deaths int // Total number of deaths for the month
}
type SharedContext struct {
	records     map[string]bool
	zipcode     int
	month       int
	year        int
	recordsFlag int32
	totalCases  int
	totalTests  int
	totalDeaths int
}
type FileTask struct {
	ctx  *SharedContext
	file int
}

func readData(filePath string, records map[string]ZipcodeInfo, kZipcode, kMonth, kYear int) {
	// some code
}

func NewFileTask(ctx *SharedContext, file int) concurrent.Runnable {
	return &FileTask{ctx, file}
}
func (task *FileTask) Run() {

	gloablCtx := task.ctx

	// records := make(map[string]ZipcodeInfo)

	// file := fmt.Sprintf("data/covid_%v.csv", task.file)
	// readData(file, records, task.ctx.zipcode, gloablCtx.month, gloablCtx.year)

	for !atomic.CompareAndSwapInt32(&gloablCtx.recordsFlag, 0, 1) {
	}

	// for key, value := range records {
	// 	if _, prs := gloablCtx.records[key]; !prs {
	// 		gloablCtx.totalCases += value.cases
	// 		gloablCtx.totalTests += value.tests
	// 		gloablCtx.totalDeaths += value.deaths
	// 		gloablCtx.records[key] = true
	// 	}
	// }

	gloablCtx.totalCases += 3
	gloablCtx.totalTests += 2
	gloablCtx.totalDeaths += 1

	atomic.StoreInt32(&gloablCtx.recordsFlag, 0)

}

func main() {

	context := SharedContext{make(map[string]bool), 606040, 5, 2020, 0, 0, 0, 0}

	threads := 3

	// executor := concurrent.NewWorkStealingExecutor(threads, 10)
	executor := concurrent.NewWorkBalancingExecutor(threads, 10, 10)
	var futures []concurrent.Future

	for i := 1; i <= 1000; i++ {
		task := NewFileTask(&context, i)
		futures = append(futures, executor.Submit(task))
	}
	executor.Shutdown()

	fmt.Printf("%v,%v,%v\n", context.totalCases, context.totalTests, context.totalDeaths)

}
