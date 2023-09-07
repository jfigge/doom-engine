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
	screenHeight = 160
)

func main() {
	graphics.Open("Text Renderer", screenWidth*4, screenHeight*4, controller.NewController())
	fmt.Println("Game over")
}
