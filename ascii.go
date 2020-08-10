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

type Vec3f struct {
	X float64
	Y float64
	Z float64
}

func (v *Vec3f) Add(rhs *Vec3f) *Vec3f {
	v.X += rhs.X
	v.Y += rhs.Y
	v.Z += rhs.Z
	return v
}

func (v *Vec3f) Dot(rhs *Vec3f) float64 {
	return v.X*rhs.X + v.Y*rhs.Y + v.Z*rhs.Z
}

func (v *Vec3f) Norm() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v *Vec3f) Normalize() *Vec3f {
	return v.ScalarMul(1.0/v.Norm())
}

func (v *Vec3f) ScalarMul(a float64) *Vec3f {
	v.X *= a
	v.Y *= a
	v.Z *= a
	return v
}

func (v *Vec3f) Sub(rhs *Vec3f) *Vec3f {
	return v.Add(rhs.ScalarMul(-1))
}

func generateCube(corner *Vec3f, sideLength int) []Vec3f {
	pts := make([]Vec3f, 0, (sideLength+1)*(sideLength+1)*(sideLength+1))
	for x := 0; x <= sideLength; x++ {
		for y := 0; y <= sideLength; y++ {
			for z := 0; z <= sideLength; z++ {
				pts = append(pts, Vec3f{corner.X+float64(x), 
										corner.Y+float64(y),
										corner.Z+float64(z)})
			}
		}
	}
	return pts
}

type PointIterator interface {
	Generate() <-chan [2]Vec3f
}

type CubeGenerator struct {
	Corner 		Vec3f
	SideLength  float64
	Step 		float64
}

/*
Creates a CubeGenerator instance and returns a pointer to it.

Parameters :
corner : coordinates of the corner with the lowest x, y, z values
sideLength : length of a side of the cube
*/
func NewCubeGenerator(corner Vec3f, sideLength, step float64) *CubeGenerator {
	return &CubeGenerator{corner, sideLength, step}
}

func (g *CubeGenerator) Generate() <-chan [2]Vec3f {
	nPts := int(g.SideLength*g.SideLength*6 - g.SideLength*12)
	c := make(chan [2]Vec3f, nPts)
	go func() {
		var r [2]Vec3f
		for z := 0.0; z < 2; z++ {
			r[1]   = Vec3f{Z: 2*z-1}
			r[0].Z = g.Corner.Z + z*(g.SideLength-1)
			for y := 0.0; y < g.SideLength; y += g.Step {
				for x := 0.0; x < g.SideLength; x += g.Step {
					r[0].X = g.Corner.X+x
					r[0].Y = g.Corner.Y+y 
					c <- r
				}
			}
		}

		for y := 0.0; y < 2; y++ {
			r[1]   = Vec3f{Y: 2*y-1}
			r[0].Y = g.Corner.Y + y*(g.SideLength-1)
			for z := 0.0; z < g.SideLength; z += g.Step {
				for x := 0.0; x < g.SideLength; x += g.Step {
					r[0].X = g.Corner.X+x
					r[0].Z = g.Corner.Z+z
					c <- r
				}
			}
		}

		for x := 0.0; x < 2; x++ {
			r[1]   = Vec3f{X: 2*x-1}
			r[0].X = g.Corner.X + x*(g.SideLength-1)
			for z := 0.0; z < g.SideLength; z += g.Step {
				for y := 0.0; y < g.SideLength; y += g.Step {
					r[0].Y = g.Corner.Y+y
					r[0].Z = g.Corner.Z+z
					c <- r
				}
			}
		}
		close(c)
	}()
	return c
}

func render(pts PointIterator, lightSource Vec3f, lightIntensity, xs, ys, zs int) [][]float64 {
	xs++
	ys++
	zs++
	camera := Vec3f{float64(xs/2), float64(ys/2), 0.0}
	//Initialize two 2D arrays for the screen and for the z-buffer.
	pixUnderlying  := make([]float64, xs*ys)
	zbufUnderlying := make([]float64, xs*ys)
	pix := make([][]float64, ys)
	zbuf := make([][]float64, ys)
	for i := range pix {
		pix[i]  = pixUnderlying[i*xs:i*xs+xs]
		zbuf[i] = zbufUnderlying[i*xs:i*xs+xs]
	}

	for r := range pts.Generate() {
		pt := r[0]
		x, y, z := pt.X-camera.X, pt.Y-camera.Y, pt.Z-camera.Z
		xp := int(float64(zs)/z * x + camera.X)
		yp := int(float64(zs)/z * y + camera.Y)
		inv_z := 1.0/z
		if xp>0 && xp<xs && 
		   yp>0 && yp<ys && 
		   inv_z > zbuf[yp][xp] {
			zbuf[yp][xp] = inv_z

			//Calculate lighting
			//d : from light source to point
			d := pt
			d.Sub(&lightSource).Normalize()
			d = Vec3f{Z:1}
			normal := r[1]
			light := d.Dot(&normal)
			pix[yp][xp] = light * float64(lightIntensity)
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
	//pts := generateCube(&Vec3f{10, 10, 8}, 3)
	g := NewCubeGenerator(Vec3f{10, 10, 8}, 3, 1)
	pix := render(g, Vec3f{0, -2, 0}, 12, 15, 15, 15)
	display(pix)

	time.Sleep(3*time.Second)

	//pts = generateCube(&Vec3f{3, 10, 8}, 3)
	g = NewCubeGenerator(Vec3f{3, 10, 8}, 3, 1)
	pix = render(g, Vec3f{0, -2, 0}, 12, 15, 15, 15)
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
2) Return error for PointIterator::Generate
*/