package utils

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"proj3/vector"
	"sync/atomic"
	"unsafe"
)

const PI = 3.141592
const R = 6
const WY = 10
const WMAX = 10

func RandomFloat64Range(r *rand.Rand, min, max float64) float64 {
	return min + r.Float64()*(max-min)
}

func GenerateRandomVector(r *rand.Rand) *vector.Vector {
	phi := RandomFloat64Range(r, 0, 2*PI)
	cos := RandomFloat64Range(r, -1.0, 1.0)
	sin := math.Sqrt(1 - cos*cos)
	return vector.NewVector(
		sin*math.Cos(phi),
		sin*math.Sin(phi),
		cos,
	)
}

func GetRandomVectors(r *rand.Rand, c *vector.Vector) (*vector.Vector, *vector.Vector) {
	for true {
		// Randomly generate view ray
		v := GenerateRandomVector(r)
		// Get intersection of the view ray with the window
		w := v.Scale(WY / v.Gety())
		// Check whether the ray is in the window as well as intersect with the ball
		vc := v.DotProduct(c)
		c2 := c.DotProduct(c)
		if math.Abs(w.Getx()) < WMAX && math.Abs(w.Getz()) < WMAX && (vc*vc+R*R-c2) > 0 {
			return v, w
		}
	}
	return nil, nil
}

func AtomicAddFloat64(val *float64, delta float64) (new float64) {
	for {
		old := *val
		new = old + delta
		if atomic.CompareAndSwapUint64(
			(*uint64)(unsafe.Pointer(val)),
			math.Float64bits(old),
			math.Float64bits(new),
		) {
			break
		}
	}
	return
}

func SaveGrid(grid [][]float64, size1d int, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		panic("Unable to open file!")
	}
	for i := 0; i < size1d; i++ {
		for j := 0; j < size1d; j++ {
			_, _ = fmt.Fprintf(file, "%.4f", grid[i][j])
			// Not the last element in a line, use \t to seperate fields
			if j != size1d-1 {
				_, _ = fmt.Fprint(file, "\t")
			}
		}
		_, _ = fmt.Fprint(file, "\n")
	}
	_ = file.Close()
}
