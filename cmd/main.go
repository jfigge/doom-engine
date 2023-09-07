/*
 * Copyright (C) 2023 by Jason Figge
 */

package main

import (
	"fmt"

	"doom-engine/internal/controller"

	"us.figge/guilib/graphics"
)

const (
	screenWidth  = 720
	screenHeight = 360
	Scale        = 2
)

func main() {
	graphics.Open("Text Renderer", screenWidth*Scale, screenHeight*Scale, controller.NewController(screenWidth/2, screenHeight))
	fmt.Println("Game over")
}
