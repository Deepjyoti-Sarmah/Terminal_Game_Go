package main

import (
	"bytes"
	"fmt"
)

const (
	NOTHING = 0
	WALL    = 1
	PLAYER  = 69
)

type game struct {
}

func (g *game) update() {

}

func (g *game) render() {

}

func main() {
	width := 80
	height := 18
	level := make([][]byte, height)

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			level[h] = make([]byte, width)
		}
	}
	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			if h == 0 {
				level[h][w] = WALL
			}
			if w == 0 {
				level[h][w] = WALL
			}
			if w == width-1 {
				level[h][w] = WALL
			}
			if h == height-1 {
				level[h][w] = WALL
			}
		}
	}
	buf := new(bytes.Buffer)
	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			if level[h][w] == NOTHING {
				buf.WriteString(" ")
			}
			if level[h][w] == WALL {
				buf.WriteString("H")
			}
		}
		buf.WriteString("\n")
	}

	fmt.Println(buf.String())
}
