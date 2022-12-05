package vector

import (
	"fmt"
	"math"
)

type Vector struct {
	x float64
	y float64
	z float64
}

func NewVector(x, y, z float64) *Vector {
	return &Vector{x: x, y: y, z: z}
}

func (v *Vector) Print() {
	fmt.Printf("(%.2f, %.2f, %.2f)\n", v.x, v.y, v.z)
}

func (v *Vector) Getx() float64 {
	return v.x
}

func (v *Vector) Gety() float64 {
	return v.y
}
func (v *Vector) Getz() float64 {
	return v.z
}

func (v *Vector) DotProduct(u *Vector) float64 {
	return v.x*u.x + v.y*u.y + v.z*u.z
}

func (v *Vector) Norm() float64 {
	return math.Sqrt(v.DotProduct(v))
}

func (v *Vector) Scale(k float64) *Vector {
	return NewVector(k*v.x, k*v.y, k*v.z)
}

func (v *Vector) LinearComb(u *Vector, kv, ku float64) *Vector {
	return NewVector(
		kv*v.x+ku*u.x,
		kv*v.y+ku*u.y,
		kv*v.z+ku*u.z,
	)
}
