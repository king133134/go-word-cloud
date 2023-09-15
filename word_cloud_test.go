package wordcloud

import (
	"fmt"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"math/rand"
	"testing"
)

func TestWordCloud_CalPixels(t *testing.T) {
	f, err := opentype.Parse(goitalic.TTF)
	if err != nil {
		log.Fatalf("Parse: %v", err)
	}
	wc := NewWordCloud(100, 100)
	wc.SetFont(f)
	words := []*Word{
		{text: "aaa", size: 13},
		{text: "bb", size: 13},
		{text: "ccc", size: 14},
		{text: "dddd", size: 14},
		{text: "e", size: 14},
		{text: "a", size: 12},
		{text: "bb", size: 13},
	}
	pixelsMap := wc.calPixels(words)
	fmt.Println(pixelsMap)
	wc.calDc.SavePNG("test/out.png")
}

func TestWordCloud_RenderEN(t *testing.T) {
	f, err := opentype.Parse(gomono.TTF)
	if err != nil {
		log.Fatalf("Parse: %v", err)
	}
	wc := NewWordCloud(400, 400)
	wc.SetFont(f)
	words := []*Word{
		{text: "a", size: 20},
		{text: "bb", size: 12},
		{text: "ccc", size: 30},
		{text: "dddd", size: 14},
		{text: "e", size: 11},
		{text: "a", size: 15},
		{text: "bb", size: 18},
		{text: "eeeeeeeeeeeeeefffss", size: 13},
	}
	wc.Render(words, "test/out.png")
}

func TestWordCloud_RenderCN(t *testing.T) {
	f, err := LoadFont("./fonts/SmileySans.ttf")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	const (
		width  = 800
		height = 450
		cm     = 1<<16 - 1
	)
	wc := NewWordCloud(width, height)
	wc.SetFont(f)
	wc.SetRotate(120, 5)
	strs := []string{"love", "movie", "animation", "music", "分子", "电子", "松子", "离子", "绝绝子", "孙子", "孟子",
		"孔子", "老子", "原子", "质子", "中子", "上帝粒子", "love", "movie", "animation", "music", "分子", "电子", "松子", "离子", "绝绝子", "孙子", "孟子",
		"孔子", "老子", "原子", "质子", "中子", "上帝粒子", "love", "movie", "animation", "music", "分子", "电子", "松子", "离子", "绝绝子", "孙子", "孟子",
		"孔子", "老子", "原子", "质子", "中子", "上帝粒子", "love", "movie", "animation", "music", "分子", "电子", "松子", "离子", "绝绝子", "孙子", "孟子",
		"孔子", "老子", "原子", "质子", "中子", "上帝粒子"}
	words := NewWords(strs)

	board, err := wc.RenderAll(words, "test/out.png", true)

	if err != nil {
		t.Fatalf("render all error: %v", err)
	}

	img := wc.dc.Image()
	cnt := 0
	for i := 0; i < 1000; {
		x, y := rand.Intn(width), rand.Intn(height)
		if !hasColor(img.At(x, y)) {
			continue
		}
		i++
		if !boardPointExists(board, y*width+x) {
			cnt++
		}
	}
	for _, w := range words {
		w.padding = 3
	}
	board, err = wc.RenderAll(words, "test/out.png", true)
	if err != nil {
		t.Fatalf("render all error: %v", err)
	}
	for i := 0; i < 1000; {
		x, y := rand.Intn(width), rand.Intn(height)
		if !hasColor(img.At(x, y)) {
			continue
		}
		i++
		if !boardPointExists(board, y*width+x) {
			cnt++
		}
	}
	if cnt > 100 {
		t.Fatalf("There are more than 10 percent incorrect points on the canvas.")
	}
}

func hasColor(color color.Color) bool {
	r, g, b, a := color.RGBA()
	return r != cm || g != cm || b != cm || a != cm
}

func boardPointExists(board []uint64, index int) bool {
	return 1<<(index&63)&board[index>>6] != 0
}

func TestWordCloud_CN1(t *testing.T) {
	f, err := LoadFont("./fonts/SmileySans.ttf")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	wc := NewWordCloud(800, 450)
	wc.SetFont(f)
	wc.SetRotate(120, 5)
	words := []*Word{{padding: 1, text: "松子", size: 52, rotate: 0}, {padding: 1, text: "love", size: 52, rotate: -60}, {padding: 1, text: "love", size: 52, rotate: -60}}
	_, err = wc.RenderAll(words, "test/out.png", true)
	if err != nil {
		t.Fatalf("render all error: %v", err)
	}
}

func TestWordCloud_isCollide1(t *testing.T) {
	f, err := opentype.Parse(goitalic.TTF)
	if err != nil {
		log.Fatalf("Parse: %v", err)
	}
	w, h := 40, 40
	wc := NewWordCloud(w, h)
	wc.SetFont(f)
	board := createBoard(w, h)
	px := &pixel{height: 8, width: 8, yOff1: 0, yOff2: 0, bits: []uint64{255}}
	wc.record(board, 8, 8, px)
	type test struct {
		name string
		x, y int
		want bool
	}

	tests := []test{
		{"test1", 8, 8, true},
		{"test2", 9, 8, true},
		{"test3", 15, 8, true},
		{"test4", 16, 8, false},
		{"test5", 18, 8, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := wc.isCollide(board, tt.x, tt.y, px); got != tt.want {
				t.Errorf("wc.isCollide() = %v, want %v", got, tt.want)
			}
		})
	}
}
