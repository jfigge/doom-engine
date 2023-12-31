/*
 * Copyright (C) 2023 by Jason Figge
 */

package controller

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
	"us.figge/guilib/graphics"
	"us.figge/guilib/graphics/fonts"
)

type DirectionCd int

const (
	FPSX         = 605 //1085
	DoV          = 200
	FoV          = 90
	D4           = 0.069813170079773
	D360         = 6.283185307179586
	WallHeight   = 40 //128
	cameraHeight = 20 //64
	backPlane    = float64(1000)
	screenPlane  = float64(0.1)
)

const (
	DirectionCdForward DirectionCd = iota
	DirectionCdBackward
	DirectionCdAntiClockwise
	DirectionCdClockwise
	DirectionCdMoveUp
	DirectionCdMoveDown
	DirectionCdLookUp
	DirectionCdLookDown
	DirectionCdStrafeLeft
	DirectionCdStrafeRight
)

type transformMatrix [4][4]float64

type Wall struct {
	x [2]float64
	y [2]float64
	z [2]float64
	c [3]uint32
}

type Section struct {
	x     float64
	y     float64
	z     uint8
	walls []Wall
}

type Camera struct {
	x     float64 // lateral
	y     float64 // depth
	l     float64 // Looking
	z     float64 // height
	a     float64 // Angle
	sin3D float64
	cos3D float64
	sin2D float64
	cos2D float64
}

type Fov struct {
	W           float64
	H           float64
	CW          float64
	CH          float64
	S           float64
	backPlane   float64
	screenPlane float64
}

type Controller struct {
	graphics.BaseHandler
	graphics.CoreMethods
	fov              *Fov
	camera           *Camera
	sections         []Section
	projectionMatrix transformMatrix
}

func NewController(width, height float64) *Controller {
	c := &Controller{
		fov: &Fov{
			W:  width,
			H:  height,
			CW: width / 2,
			CH: height / 2,
			S:  DoV,
		},
		camera: &Camera{
			x:     160,
			y:     160,
			z:     20,
			a:     0,
			l:     0,
			sin3D: math.Sin(0),
			cos3D: math.Cos(0),
			sin2D: math.Sin(0),
			cos2D: math.Cos(0),
		},
		sections: []Section{
			{
				walls: []Wall{
					//{x: [2]float64{-25, -25}, y: [2]float64{25, -75}, z: [2]float64{0, 40}, c: [3]uint32{0xFF0000FF, 0x00FF00FF, 0x000000FFFF}},
					//{x: [2]float64{-25, 25}, y: [2]float64{-75, -75}, z: [2]float64{0, 40}, c: [3]uint32{0x00FF00FF, 0x00FF00FF, 0x000000FFFF}},
					//{x: [2]float64{25, 25}, y: [2]float64{-75, 25}, z: [2]float64{0, 40}, c: [3]uint32{0x0000FFFF, 0x00FF00FF, 0x000000FFFF}},
					//{x: [2]float64{25, -25}, y: [2]float64{25, 25}, z: [2]float64{0, 40}, c: [3]uint32{0xFFFF00FF, 0x00FF00FF, 0x000000FFFF}},

					{x: [2]float64{175, 175}, y: [2]float64{125, 75}, z: [2]float64{10, 40}, c: [3]uint32{0xFF0000FF, 0x00FF00FF, 0x000000FFFF}},
					{x: [2]float64{175, 225}, y: [2]float64{75, 75}, z: [2]float64{10, 40}, c: [3]uint32{0x00FF00FF, 0x00FF00FF, 0x000000FFFF}},
					{x: [2]float64{225, 225}, y: [2]float64{75, 125}, z: [2]float64{10, 40}, c: [3]uint32{0x0000FFFF, 0x00FF00FF, 0x000000FFFF}},
					{x: [2]float64{225, 175}, y: [2]float64{125, 125}, z: [2]float64{10, 40}, c: [3]uint32{0xFFFF00FF, 0x00FF00FF, 0x000000FFFF}},
				},
			},
		},
	}

	a := width / height
	f := 1 / math.Tan(FoV*math.Pi/360)
	q := backPlane / (backPlane - screenPlane)
	c.projectionMatrix = transformMatrix{
		{a * f, 0, 0, 0},
		{0, f, 0, 0},
		{0, 0, q, 1},
		{0, 0, -screenPlane * q, 0},
	}

	return c
}

func (c *Controller) Init(canvas *graphics.Canvas) {
	fonts.LoadFonts(canvas.Renderer())
	graphics.ErrorTrap(canvas.Renderer().SetDrawBlendMode(sdl.BLENDMODE_BLEND))
	canvas.Renderer().SetLogicalSize(int32(c.fov.W*2), int32(c.fov.H))
	c.AddDestroyer(fonts.FreeFonts)
}

func (c *Controller) OnUpdate() {
	c.processKeys()
}

func (c *Controller) OnDraw(renderer *sdl.Renderer) {
	graphics.ErrorTrap(c.Clear(renderer, uint32(0x232323)))
	c.draw2D(renderer)
	c.draw3D(renderer)
	graphics.ErrorTrap(c.WriteFrameRate(renderer, FPSX, 0))
}

func (c *Controller) draw2D(renderer *sdl.Renderer) {
	offset := float64(c.fov.W)
	renderer.SetDrawColor(uint8(0), uint8(0), uint8(0xFF), uint8(0xFF))
	renderer.DrawLine(int32(c.fov.W), 0, int32(c.fov.W), int32(c.fov.H))

	renderer.SetDrawColor(uint8(0xFF), uint8(0xFF), uint8(0xFF), uint8(0xFF))
	renderer.DrawLinesF([]sdl.FPoint{
		c.rotate(c.camera.x+offset, c.camera.y, c.camera.x+offset, c.camera.y),
		c.rotate(c.camera.x-5+offset, c.camera.y+20, c.camera.x+offset, c.camera.y),
		c.rotate(c.camera.x+5+offset, c.camera.y+20, c.camera.x+offset, c.camera.y),
		c.rotate(c.camera.x+offset, c.camera.y, c.camera.x+offset, c.camera.y),
	})
	renderer.SetDrawColor(uint8(0xFF), uint8(0), uint8(0), uint8(0xFF))
	renderer.DrawPointF(float32(c.camera.x+offset), float32(c.camera.y))

	for _, section := range c.sections {
		for _, w := range section.walls {
			renderer.SetDrawColor(uint8(w.c[0]>>24), uint8(w.c[0]>>16), uint8(w.c[0]>>8), uint8(w.c[0]))
			renderer.DrawLineF(float32(w.x[0]+offset), float32(w.y[0]), float32(w.x[1]+offset), float32(w.y[1]))
		}
	}
}

func (c *Controller) rotate(x, y, ox, oy float64) sdl.FPoint {
	return sdl.FPoint{
		X: float32((x-ox)*c.camera.cos2D - (y-oy)*c.camera.sin2D + ox),
		Y: float32((y-oy)*c.camera.cos2D + (x-ox)*c.camera.sin2D + oy),
	}

}

func (c *Controller) project(x float64, y float64, z float64, w float64) (float64, float64, float64) {
	mx := x*c.projectionMatrix[0][0] + y*c.projectionMatrix[1][0] + z*c.projectionMatrix[2][0] + w*c.projectionMatrix[3][0]
	my := x*c.projectionMatrix[0][1] + y*c.projectionMatrix[1][1] + z*c.projectionMatrix[2][1] + w*c.projectionMatrix[3][1]
	mz := x*c.projectionMatrix[0][2] + y*c.projectionMatrix[1][2] + z*c.projectionMatrix[2][2] + w*c.projectionMatrix[3][2]
	mw := x*c.projectionMatrix[0][3] + y*c.projectionMatrix[1][3] + z*c.projectionMatrix[2][3] + w*c.projectionMatrix[3][3]
	fmt.Println(mw)
	if mw != 0 {
		//mx = mx / z
		//my = my / z
	}
	fmt.Println(mw)
	return mx, my, mz
}

func (c *Controller) translate(x, y, z float64, cam *Camera) (float32, float32) {
	// Offset to 0,0,0 based on camera location
	dx := x - cam.x
	dy := y - cam.y
	dz := z - cam.z

	// Rotate
	rdx := dx*cam.cos3D - dy*cam.sin3D
	rdy := dy*cam.cos3D + dx*cam.sin3D

	// project into 3D space
	mx, _, mz := c.project(rdx, rdy, dz, 1)

	return float32(mx + c.fov.CW), float32(mz + c.fov.CH)
}

func mark(r *sdl.Renderer, x, y float64) {
	r.SetDrawColor(255, 255, 255, 255)
	r.DrawLineF(float32(x-2), float32(y-2), float32(x+2), float32(y+2))
	r.DrawLineF(float32(x+2), float32(y-2), float32(x-2), float32(y+2))
}
func (c *Controller) draw3D(renderer *sdl.Renderer) {
	mark(renderer, c.fov.CW, c.fov.CH)
	var wx, wy [4]float32
	//wv := make([]sdl.Vertex, 6)
	for _, section := range c.sections {
		for _, w := range section.walls {
			wx[0], wy[0] = c.translate(w.x[0], w.z[0], w.y[0], c.camera)
			wx[1], wy[1] = c.translate(w.x[1], w.z[1], w.y[0], c.camera)
			wx[2], wy[2] = c.translate(w.x[0], w.z[0], w.y[1], c.camera)
			wx[3], wy[3] = c.translate(w.x[1], w.z[1], w.y[1], c.camera)

			renderer.SetDrawColor(uint8(w.c[0]>>24), uint8(w.c[0]>>16), uint8(w.c[0]>>8), uint8(w.c[0]))
			renderer.DrawLinesF([]sdl.FPoint{
				{wx[0], wy[0]},
				{wx[1], wy[1]},
				{wx[2], wy[2]},
				{wx[0], wy[0]},
				{wx[3], wy[3]},
				{wx[3], wy[2]},
			})
			//wv[0] = sdl.Vertex{
			//	Position: sdl.FPoint{X: wx[0], Y: wy[0]},
			//	Color:    sdl.Color{R: uint8(w.c[0] >> 24), G: uint8(w.c[0] >> 16), B: uint8(w.c[0] >> 8), A: uint8(w.c[0])},
			//	TexCoord: sdl.FPoint{},
			//}
			//wv[1] = sdl.Vertex{
			//	Position: sdl.FPoint{X: wx[1], Y: wy[1]},
			//	Color:    sdl.Color{R: uint8(w.c[0] >> 24), G: uint8(w.c[0] >> 16), B: uint8(w.c[0] >> 8), A: uint8(w.c[0])},
			//	TexCoord: sdl.FPoint{},
			//}
			//wv[2] = sdl.Vertex{
			//	Position: sdl.FPoint{X: wx[2], Y: wy[2]},
			//	Color:    sdl.Color{R: uint8(w.c[0] >> 24), G: uint8(w.c[0] >> 16), B: uint8(w.c[0] >> 8), A: uint8(w.c[0])},
			//	TexCoord: sdl.FPoint{},
			//}
			//wv[3] = sdl.Vertex{
			//	Position: sdl.FPoint{X: wx[3], Y: wy[3]},
			//	Color:    sdl.Color{R: uint8(w.c[0] >> 24), G: uint8(w.c[0] >> 16), B: uint8(w.c[0] >> 8), A: uint8(w.c[0])},
			//	TexCoord: sdl.FPoint{},
			//}
			//wv[4] = sdl.Vertex{
			//	Position: sdl.FPoint{X: wx[1], Y: wy[1]},
			//	Color:    sdl.Color{R: uint8(w.c[0] >> 24), G: uint8(w.c[0] >> 16), B: uint8(w.c[0] >> 8), A: uint8(w.c[0])},
			//	TexCoord: sdl.FPoint{},
			//}
			//wv[5] = sdl.Vertex{
			//	Position: sdl.FPoint{X: wx[2], Y: wy[2]},
			//	Color:    sdl.Color{R: uint8(w.c[0] >> 24), G: uint8(w.c[0] >> 16), B: uint8(w.c[0] >> 8), A: uint8(w.c[0])},
			//	TexCoord: sdl.FPoint{},
			//}
			//renderer.RenderGeometry(nil, wv, nil)
		}
	}
}

func (c *Controller) processKeys() {
	codes := sdl.GetKeyboardState()
	shift := codes[sdl.SCANCODE_LSHIFT] == 1 || codes[sdl.SCANCODE_RSHIFT] == 1
	if codes[sdl.SCANCODE_W] == 1 {
		switch {
		case shift:
			c.move(DirectionCdLookUp)
		case codes[sdl.SCANCODE_M] == 1:
			c.move(DirectionCdMoveUp)
		default:
			c.move(DirectionCdForward)
		}
	} else if codes[sdl.SCANCODE_S] == 1 {
		switch {
		case shift:
			c.move(DirectionCdLookDown)
		case codes[sdl.SCANCODE_M] == 1:
			c.move(DirectionCdMoveDown)
		default:
			c.move(DirectionCdBackward)
		}
	}
	if codes[sdl.SCANCODE_COMMA] == 1 {
		c.move(DirectionCdAntiClockwise)
	} else if codes[sdl.SCANCODE_PERIOD] == 1 {
		c.move(DirectionCdClockwise)
	}
	if codes[sdl.SCANCODE_A] == 1 {
		c.move(DirectionCdStrafeLeft)
	} else if codes[sdl.SCANCODE_D] == 1 {
		c.move(DirectionCdStrafeRight)
	}
}

func (c *Controller) move(dir DirectionCd) {
	dx := 8 * c.camera.sin2D
	dy := 8 * c.camera.cos2D
	switch dir {
	case DirectionCdForward:
		c.camera.x += dx
		c.camera.y -= dy
	case DirectionCdBackward:
		c.camera.x -= dx
		c.camera.y += dy
	case DirectionCdStrafeLeft:
		c.camera.x -= dy
		c.camera.y -= dx
	case DirectionCdStrafeRight:
		c.camera.x += dy
		c.camera.y += dx
	case DirectionCdMoveUp:
		c.camera.z += 4
		if c.camera.z > c.fov.H {
			c.camera.z = c.fov.H
		}
		fmt.Println(c.camera.z)
	case DirectionCdMoveDown:
		c.camera.z -= 4
		if c.camera.z < -c.fov.H {
			c.camera.z = c.fov.H
		}
	case DirectionCdLookUp:
		c.camera.l += 1
		if c.camera.l > cameraHeight {
			c.camera.l = cameraHeight
		}
	case DirectionCdLookDown:
		c.camera.l -= 1
		if c.camera.l < -cameraHeight {
			c.camera.l = -cameraHeight
		}
	case DirectionCdAntiClockwise:
		c.camera.a -= D4
		if c.camera.a < 0 {
			c.camera.a += D360
		}
		c.camera.sin3D = math.Sin(-c.camera.a)
		c.camera.cos3D = math.Cos(-c.camera.a)
		c.camera.sin2D = math.Sin(c.camera.a)
		c.camera.cos2D = math.Cos(c.camera.a)
	case DirectionCdClockwise:
		c.camera.a += D4
		if c.camera.a > D360 {
			c.camera.a -= D360
		}
		c.camera.sin3D = math.Sin(-c.camera.a)
		c.camera.cos3D = math.Cos(-c.camera.a)
		c.camera.sin2D = math.Sin(c.camera.a)
		c.camera.cos2D = math.Cos(c.camera.a)
	}
}
