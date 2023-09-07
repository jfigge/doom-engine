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
	screenWidth  = 320
	screenHeight = 120
	Scale        = 4
)

func main() {
	graphics.Open("Text Renderer", screenWidth*Scale, screenHeight*Scale, controller.NewController(screenWidth/2, screenHeight))
	fmt.Println("Game over")
}
