/*
 * Copyright (C) 2023 by Jason Figge
 */

package controller

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
	"us.figge/guilib/graphics"
	"us.figge/guilib/graphics/fonts"
)

type DirectionCd int

const (
	FPSX         = 205 //1085
	DoV          = 200
	D4           = 0.069813170079773
	D360         = 6.283185307179586
	WallHeight   = 40 //128
	PlayerHeight = 20 //64
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

type Entity struct {
	x    float64 // lateral
	y    float64 // depth
	l    float64 // Looking
	z    float64 // height
	a    float64 // Angle
	sin1 float64
	cos1 float64
	sin2 float64
	cos2 float64
}

type Fov struct {
	W  float64
	H  float64
	CW float64
	CH float64
	S  float64
}

type Controller struct {
	graphics.BaseHandler
	graphics.CoreMethods
	player   *Entity
	sections []Section
	fov      *Fov
}

func NewController(width, height uint32) *Controller {
	c := &Controller{
		fov: &Fov{
			W:  float64(width),
			H:  float64(height),
			CW: float64(width / 2),
			CH: float64(height / 2),
			S:  DoV,
		},
		player: &Entity{
			x:    70,
			y:    -110,
			z:    20,
			a:    0,
			l:    0,
			sin1: math.Sin(0),
			cos1: math.Cos(0),
			sin2: math.Sin(0),
			cos2: math.Cos(0),
		},
		sections: []Section{
			{
				walls: []Wall{
					{x: [2]float64{40, 40}, y: [2]float64{10, 290}, z: [2]float64{0, 40}, c: [3]uint32{0xFF0000FF, 0x00FF00FF, 0x000000FFFF}},
					//{x: [2]float64{350, 450}, y: [2]float64{150, 150}, z: [2]float64{0, 40}, c: [3]uint32{0x00FF00FF, 0x00FF00FF, 0x000000FFFF}},
					//{x: [2]float64{450, 450}, y: [2]float64{150, 250}, z: [2]float64{0, 40}, c: [3]uint32{0x0000FFFF, 0x00FF00FF, 0x000000FFFF}},
					//{x: [2]float64{450, 350}, y: [2]float64{250, 250}, z: [2]float64{0, 40}, c: [3]uint32{0xFFFF00FF, 0x00FF00FF, 0x000000FFFF}},
				},
			},
		},
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
	//offset := float64(c.fov.W)
	renderer.SetDrawColor(uint8(0), uint8(0), uint8(0xFF), uint8(0xFF))
	renderer.DrawLine(int32(c.fov.W), 0, int32(c.fov.W), int32(c.fov.H))
	//
	//renderer.SetDrawColor(uint8(0xFF), uint8(0), uint8(0), uint8(0xFF))
	//renderer.DrawPointF(float32(c.player.x+offset), float32(-c.player.y))
	//
	//for _, section := range c.sections {
	//	for _, w := range section.walls {
	//		renderer.SetDrawColor(uint8(w.c[0]>>24), uint8(w.c[0]>>16), uint8(w.c[0]>>8), uint8(w.c[0]))
	//		renderer.DrawLineF(float32(w.x[0]+offset), float32(w.y[0]), float32(w.x[1]+offset), float32(w.y[1]))
	//	}
	//}
}

func (c *Controller) rotate(x, y, ox, oy float64) sdl.FPoint {
	return sdl.FPoint{
		X: float32((x-ox)*c.player.cos2 - (y-oy)*c.player.sin2 + ox),
		Y: float32((y-oy)*c.player.cos2 + (x-ox)*c.player.sin2 + oy),
	}

}

func (c *Controller) translate(r *sdl.Renderer, x, y, z float64, e *Entity) (float32, float32) {
	// offset coordinates by entity
	dx := x - e.x
	dy := y - e.y
	dz := z / (e.z * 2)

	// rotate around entity
	rdx := dx*e.cos1 - dy*e.sin1
	rdy := dy*e.cos1 + dx*e.sin1

	scale := 3 - dz

	// Convert to screen position
	rdx = rdx*scale + c.fov.CW
	rdy = rdy*scale + c.fov.CH

	return float32(rdx), float32(rdy)
}

//	func (c *Controller) drawWall(renderer sdl.Renderer, x1, x2, b1, b2 float64) {
//		var x, y float64
//		dyb := b2 - b1
//		dx := x2 - x1
//		if dx == 0 {
//			dx = 1
//		}
//		xs := x1
//		for x = x1; x < x2; x++ {
//			y1 := dyb*(x-xs+.05)/dx + b1
//		}
//	}
func (c *Controller) draw3D(renderer *sdl.Renderer) {
	var wx, wy, wz [4]float64

	x1 := 40 - c.player.x
	y1 := 10 - c.player.y
	x2 := 40 - c.player.x
	y2 := 290 - c.player.y

	wx[0] = x1*c.player.cos1 + y1*c.player.sin1
	wx[1] = x2*c.player.cos1 + y2*c.player.sin1

	wy[0] = y1*c.player.cos1 - x1*c.player.sin1
	wy[1] = y2*c.player.cos1 - x2*c.player.sin1

	wz[0] = 0 - c.player.z + ((c.player.l * wy[0]) / 32)
	wz[1] = 0 - c.player.z + ((c.player.l * wy[1]) / 32)

	wx[0] = wx[0]*200/wy[0] + c.fov.CW
	wx[1] = wx[1]*200/wy[1] + c.fov.CW
	wy[0] = wz[0]*200/wy[0] + c.fov.CH
	wy[1] = wz[1]*200/wy[1] + c.fov.CH

	renderer.SetDrawColor(255, 255, 255, 255)
	if wx[0] > 0 && wx[0] < c.fov.W && wy[0] > 0 && wy[0] < c.fov.H {
		renderer.DrawPointF(float32(wx[0]), float32(wy[0]))
	}
	if wx[1] > 0 && wx[1] < c.fov.W && wy[1] > 0 && wy[1] < c.fov.H {
		renderer.DrawPointF(float32(wx[1]), float32(wy[1]))
	}
}

func (c *Controller) Xdraw3D(renderer *sdl.Renderer) {
	var wx, wy [4]float32
	for _, section := range c.sections {
		for _, w := range section.walls {
			wx[0], wy[0] = c.translate(renderer, w.x[0], w.y[0], w.z[0], c.player)
			wx[1], wy[1] = c.translate(renderer, w.x[1], w.y[1], w.z[0], c.player)
			wx[2], wy[2] = c.translate(renderer, w.x[0], w.y[0], w.z[1], c.player)
			wx[3], wy[3] = c.translate(renderer, w.x[1], w.y[1], w.z[1], c.player)

			renderer.SetDrawColor(uint8(w.c[0]>>24), uint8(w.c[0]>>16), uint8(w.c[0]>>8), uint8(w.c[0]))
			renderer.DrawLineF(wx[0], wy[0], wx[1], wy[1])
			renderer.DrawLineF(wx[2], wy[2], wx[3], wy[3])

			renderer.DrawLineF(wx[0], wy[0], wx[2], wy[2])
			renderer.DrawLineF(wx[1], wy[1], wx[3], wy[3])
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
	dx := 8 * c.player.sin2
	dy := 8 * c.player.cos2
	switch dir {
	case DirectionCdForward:
		c.player.x += dx
		c.player.y += dy
	case DirectionCdBackward:
		c.player.x -= dx
		c.player.y -= dy
	case DirectionCdStrafeLeft:
		c.player.x += dy
		c.player.y += dx
	case DirectionCdStrafeRight:
		c.player.x -= dy
		c.player.y -= dx
	case DirectionCdMoveUp:
		c.player.z += 4
		if c.player.z > WallHeight {
			c.player.z = WallHeight
		}
	case DirectionCdMoveDown:
		c.player.z -= 4
		if c.player.z > 0 {
			c.player.z = 0
		}
	case DirectionCdLookUp:
		c.player.l -= 1
		if c.player.l < -WallHeight {
			c.player.l = -WallHeight
		}
	case DirectionCdLookDown:
		c.player.l += 1
		if c.player.l > WallHeight {
			c.player.l = WallHeight
		}
	case DirectionCdAntiClockwise:
		c.player.a -= D4
		if c.player.a < 0 {
			c.player.a += D360
		}
		c.player.sin1 = math.Sin(-c.player.a)
		c.player.cos1 = math.Cos(-c.player.a)
		c.player.sin2 = math.Sin(c.player.a)
		c.player.cos2 = math.Cos(c.player.a)
	case DirectionCdClockwise:
		c.player.a += D4
		if c.player.a > D360 {
			c.player.a -= D360
		}
		c.player.sin1 = math.Sin(-c.player.a)
		c.player.cos1 = math.Cos(-c.player.a)
		c.player.sin2 = math.Sin(c.player.a)
		c.player.cos2 = math.Cos(c.player.a)
	}
}
