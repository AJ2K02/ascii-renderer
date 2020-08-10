package main

import (
	"fmt"
	"math"
	"time"
)

type Vec3 interface {
	Dot(rhs Vec3)
	Add(rhs Vec3)
} 

type Vec3i struct {
	X int
	Y int
	Z int
}

func (v *Vec3i) Add(rhs *Vec3i) *Vec3i {
	v.X += rhs.X
	v.Y += rhs.Y
	v.Z += rhs.Z
	return v
}

func (v *Vec3i) Dot(rhs *Vec3i) int {
	return v.X*rhs.X + v.Y*rhs.Y + v.Z*rhs.Z
}

func (v *Vec3i) Norm() float64 {
	return math.Sqrt(float64(v.Dot(v)))
}

func (v *Vec3i) ScalarMul(a int) *Vec3i {
	v.X *= a
	v.Y *= a
	v.Z *= a
	return v
}

func (v *Vec3i) Sub(rhs *Vec3i) *Vec3i {
	return v.Add(rhs.ScalarMul(-1))
}

func generateCube(corner *Vec3i, sideLength int) []Vec3i {
	pts := make([]Vec3i, 0, (sideLength+1)*(sideLength+1)*(sideLength+1))
	for x := 0; x <= sideLength; x++ {
		for y := 0; y <= sideLength; y++ {
			for z := 0; z <= sideLength; z++ {
				pts = append(pts, Vec3i{corner.X+x, 
										corner.Y+y,
										corner.Z+z})
			}
		}
	}
	return pts
}

func render(pts []Vec3i, lightSource Vec3i, lightIntensity, xs, ys, zs int) [][]float64 {
	xs++
	ys++
	zs++
	camera := Vec3i{xs/2, ys/2, 0}
	//Initialize two 2D arrays for the screen and for the z-buffer.
	pixUnderlying  := make([]float64, xs*ys)
	zbufUnderlying := make([]float64, xs*ys)
	pix := make([][]float64, ys)
	zbuf := make([][]float64, ys)
	for i := range pix {
		pix[i]  = pixUnderlying[i*xs:i*xs+xs]
		zbuf[i] = zbufUnderlying[i*xs:i*xs+xs]
	}

	for _, pt := range pts {
		x, y, z := pt.X-camera.X, pt.Y-camera.Y, pt.Z-camera.Z
		xp := zs/z * x + camera.X
		yp := zs/z * y + camera.Y
		inv_z := 1.0/float64(z)
		if xp>0 && xp<xs && 
		   yp>0 && yp<ys && 
		   inv_z > zbuf[yp][xp] {
			zbuf[yp][xp] = inv_z

			//Calculate lighting
			//d : 		from light source to point
			//x, y, z : from camera       to point
			//should be surface normal...
			d := pt
			d.Sub(&lightSource)
			t := Vec3i{x, y, z}
			fmt.Println(d)
			light := d.Dot(&t)
			pix[yp][xp] = float64(light * lightIntensity) / (d.Norm()*t.Norm())
		}
	}
	return pix
}

func display(pix [][]float64) {
	fmt.Print("\033[2J")
	for _, row := range pix {
		for _, val := range row {
			var c byte
			if val < 0 {
				c = "."[0]
			} else if val > 11 {
				c = "@"[0]
			} else {
				c = ".,-~:;=!*#$@"[int(val)]
			}
			fmt.Print(string(c) + " ")
		}
		fmt.Println("\n")
	}
}

func main() {
	pts := generateCube(&Vec3i{10, 10, 8}, 3)
	pix := render(pts, Vec3i{0, -2, 0}, 12, 15, 15, 15)
	display(pix)

	time.Sleep(3*time.Second)

	pts = generateCube(&Vec3i{3, 10, 8}, 3)
	pix = render(pts, Vec3i{0, -2, 0}, 12, 15, 15, 15)
	display(pix)
}

/*

Plan :

1) Start by creating an iterable of all points to render (DONE)
2) Create a function to render a list of points (BUG??)
3) Display the pixels (DONE)

Questions along the way :
1) Is it ok to create a slice with make(type, 0, 84566565) and then append
2) When to pass by pointer/copy

TODO :
1) Min / max z distance for points
*/