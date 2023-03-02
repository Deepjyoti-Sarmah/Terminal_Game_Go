package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"
)

const (
	NOTHING = 0
	WALL    = 1
	PLAYER  = 69

	MAX_SAMPLES = 100
)

type input struct {
	pressedKey byte
}

func (i *input) update() {
	b := make([]byte, 1)
	os.Stdin.Read(b)
	i.pressedKey = b[0]
}

type position struct {
	x, y int
}

type player struct {
	pos   position
	level *level
	input *input

	reverse bool
}

func (p *player) update() {
	if p.reverse {
		p.pos.x -= 1
		if p.pos.x == 2 {
			p.pos.x += 1
			p.reverse = false
		}
		return
	}

	p.pos.x += 1
	if p.pos.x == p.level.width-1 {
		p.pos.x -= 1
		p.reverse = true
	}
}

type stats struct {
	start  time.Time
	frames int
	fps    float64
}

func newStats() *stats {
	return &stats{
		fps:   69,
		start: time.Now(),
	}
}

func (s *stats) update() {
	s.frames++
	if s.frames == MAX_SAMPLES {
		s.fps = float64(s.frames) / time.Since(s.start).Seconds()
		s.frames = 0
		s.start = time.Now()
	}
}

type level struct {
	width, height int
	data          [][]int
}

func newLevel(width, height int) *level {
	data := make([][]int, height)

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			data[h] = make([]int, width)
		}
	}
	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			if h == 0 {
				data[h][w] = WALL
			}
			if w == 0 {
				data[h][w] = WALL
			}
			if w == width-1 {
				data[h][w] = WALL
			}
			if h == height-1 {
				data[h][w] = WALL
			}
		}
	}
	return &level{
		width:  width,
		height: height,
		data:   data,
	}
}

func (l *level) set(pos position, v int) {
	l.data[pos.y][pos.x] = v
}

type game struct {
	isRunning bool
	level     *level
	stats     *stats
	player    *player
	input     *input
	drawBuf   *bytes.Buffer
}

func newGame(width, height int) *game {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	var (
		lvl  = newLevel(width, height)
		inpu = &input{}
	)
	return &game{
		level:   lvl,
		drawBuf: new(bytes.Buffer),
		stats:   newStats(),
		input:   inpu,
		player: &player{
			input: inpu,
			level: lvl,
			pos:   position{x: 2, y: 5},
		},
	}
}

func (g *game) start() {
	g.isRunning = true
	g.loop()
}

func (g *game) loop() {
	for g.isRunning {
		g.input.update()
		g.update()
		g.render()
		g.stats.update()
		time.Sleep(time.Millisecond * 16) //limits FPS
	}
}

func (g *game) update() {
	g.level.set(g.player.pos, NOTHING)
	g.player.update()
	g.level.set(g.player.pos, PLAYER)
}

func (g *game) renderLevel() {
	for h := 0; h < g.level.height; h++ {
		for w := 0; w < g.level.width; w++ {
			if g.level.data[h][w] == NOTHING {
				g.drawBuf.WriteString(" ")
			}
			if g.level.data[h][w] == WALL {
				g.drawBuf.WriteString("▢")
			}
			if g.level.data[h][w] == PLAYER {
				g.drawBuf.WriteString("☺")
			}
		}
		g.drawBuf.WriteString("\n")
	}
}

func (g *game) render() {
	g.drawBuf.Reset()
	fmt.Fprint(os.Stdout, "\033[23\033[1;1H")
	g.renderLevel()
	g.renderStats()
	fmt.Fprint(os.Stdout, g.drawBuf.String())
}

func (g *game) renderStats() {
	g.drawBuf.WriteString("-- STATS\n")
	g.drawBuf.WriteString(fmt.Sprintf("FPS: %.2f\n", g.stats.fps))
}

func main() {
	width := 80
	height := 18

	g := newGame(width, height)
	g.start()
}
