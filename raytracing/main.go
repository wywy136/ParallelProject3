package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"proj3/concurrent"
	"proj3/utils"
	"proj3/vector"
	"strconv"
	"time"
)

const usage = "go run main.go number_of_rays size_of_grid strategy number_of_threads number_of_subtasks\n" +
	"number_of_rays: the number of simulated rays for the rendering. Recommended: 100000000\n" +
	"size_of_grid: number of pixels along one dimension. Recommended: 1000\n" +
	"strategy: s - sequential; ws - work stealing; wb - work balancing\n" +
	"number_of_threads: number of threads for parallel strategy\n" +
	"number_of_subtasks: number of subtasks for the strategy to manage. Recommended: 1000"

type SharedContext struct {
	numRays   int
	numGrid1d int
	grids     [][]float64
}

type RayTracingTask struct {
	ctx             *SharedContext
	randomGenerator *rand.Rand
}

func NewRayTracingTask(ctx *SharedContext, seed int64) concurrent.Runnable {
	// Generate a new random generator for every task
	// Using the current system time as seed
	return &RayTracingTask{
		ctx:             ctx,
		randomGenerator: rand.New(rand.NewSource(seed)),
	}
}

func (t *RayTracingTask) Run() {
	// Light source
	L := vector.NewVector(4.0, 4.0, -1.0)
	// Position of the ball
	C := vector.NewVector(0.0, 12.0, 0.0)
	// For every light of this task
	for i := 0; i < t.ctx.numRays; i++ {
		V, W := utils.GetRandomVectors(t.randomGenerator, C)
		VC := V.DotProduct(C)
		term := VC - math.Sqrt(VC*VC+utils.R*utils.R-C.DotProduct(C))
		I := V.Scale(term)
		IMC := I.LinearComb(C, 1, -1)
		N := IMC.Scale(1 / IMC.Norm())
		LMI := L.LinearComb(I, 1, -1)
		S := LMI.Scale(1 / LMI.Norm())

		var b float64
		b = math.Max(0.0, S.DotProduct(N))

		var iIndex int
		var jIndex int
		jIndex = int(W.Getz()*float64(t.ctx.numGrid1d/2)/utils.WMAX + float64(t.ctx.numGrid1d)/2)
		iIndex = int(float64(t.ctx.numGrid1d)/2 - W.Getx()*float64(t.ctx.numGrid1d/2)/utils.WMAX)

		utils.AtomicAddFloat64(&t.ctx.grids[iIndex][jIndex], b)
	}
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Invalid command-line arguments!")
		fmt.Println(usage)
		return
	}

	nRays, _ := strconv.Atoi(os.Args[1])
	nGrid1d, _ := strconv.Atoi(os.Args[2])

	if nGrid1d < 1000 {
		fmt.Println("To ensure good visualization, n_grid should be at least 1000. Recommend: 1000")
		return
	}

	if nRays < nGrid1d*nGrid1d {
		fmt.Println("To ensure good visualization, n_rays should be at least n_grid*n_grid. Recommend: 100*n_grid*n_grid")
		return
	}

	grid := make([][]float64, nGrid1d)
	for i := 0; i < nGrid1d; i++ {
		grid[i] = make([]float64, nGrid1d)
	}

	strategy := os.Args[3]
	ctx := SharedContext{
		numRays:   nRays,
		numGrid1d: nGrid1d,
		grids:     grid,
	}

	// Sequential
	start := time.Now()
	if strategy == "s" {
		task := NewRayTracingTask(&ctx, time.Now().Unix())
		task.Run()
	} else {
		// Argument checking for parallel
		if len(os.Args) < 6 {
			fmt.Println("Invalid command-line arguments!")
			fmt.Println(usage)
			return
		}
		nThreads, _ := strconv.Atoi(os.Args[4])
		if nThreads <= 1 {
			fmt.Println("Number of threads should be greater than 1 for work stealing / balancing strategy.")
			return
		}

		nTasks, _ := strconv.Atoi(os.Args[5])
		if nTasks < nThreads {
			fmt.Println("Number of sub-tasks should be at least the number of threads.")
			return
		}

		nRaysPerTask := int(nRays / nTasks)
		ctx.numRays = nRaysPerTask

		var futures []concurrent.Future
		randSeed := time.Now().Unix()

		if strategy == "ws" {
			executor := concurrent.NewWorkStealingExecutor(nThreads, 10)
			for i := 0; i < nTasks; i++ {
				futures = append(futures, executor.Submit(NewRayTracingTask(&ctx, randSeed)))
				randSeed++
			}
			executor.Shutdown()
		} else if strategy == "wb" {
			executor := concurrent.NewWorkBalancingExecutor(nThreads, 10, 2)
			for i := 0; i < nTasks; i++ {
				futures = append(futures, executor.Submit(NewRayTracingTask(&ctx, randSeed)))
				randSeed++
			}
			executor.Shutdown()
		} else {
			fmt.Println("Unknown strategy.")
			fmt.Println(usage)
			return
		}
	}
	// fmt.Print(ctx.grids)
	elapsed := time.Since(start).Seconds()
	fmt.Println("Elapsed time: ", elapsed)

	utils.SaveGrid(ctx.grids, nGrid1d, "output")

}
