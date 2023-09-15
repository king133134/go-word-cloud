package wordcloud

import (
	"fmt"
	"golang.org/x/image/font/opentype"
	"image/color"
	"math/rand"
	"os"
	"sort"
	"strings"
)

var colors = [20]string{
	"#393b79",
	"#5254a3",
	"#6b6ecf",
	"#9c9ede",
	"#637939",
	"#8ca252",
	"#b5cf6b",
	"#cedb9c",
	"#8c6d31",
	"#bd9e39",
	"#e7ba52",
	"#e7cb94",
	"#843c39",
	"#ad494a",
	"#d6616b",
	"#e7969c",
	"#7b4173",
	"#a55194",
	"#ce6dbd",
	"#de9ed6",
}

const noRotated = 360

type Word struct {
	text    string
	size    int
	padding int
	color   color.Color
	rotate  int
}

func parseHexColor(x string) (r, g, b, a uint8) {
	x = strings.TrimPrefix(x, "#")
	a = 255
	format := ""
	if l := len(x); l == 3 {
		format = "%1x%1x%1x"
		_, _ = fmt.Sscanf(x, format, &r, &g, &b)
		r |= r << 4
		g |= g << 4
		b |= b << 4
		return
	} else if l == 6 {
		format = "%02x%02x%02x"
	} else if l == 8 {
		format = "%02x%02x%02x%02x"
	}
	_, _ = fmt.Sscanf(x, format, &r, &g, &b)
	return
}

func NewRandColor() color.Color {
	idx := rand.Intn(len(colors))
	r, g, b, a := parseHexColor(colors[idx])
	return color.RGBA{R: r, G: g, B: b, A: a}
}

func NewWords(words []string) []*Word {
	res := make([]*Word, len(words))
	for i, v := range words {
		res[i] = NewWord(v, 10+rand.Intn(46), Color(NewRandColor()), Padding(0))
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].size > res[j].size
	})
	return res
}

type WordOption func(w *Word)

func Color(color color.Color) WordOption {
	return func(w *Word) {
		w.color = color
	}
}

func HexColor(hex string) WordOption {
	return func(w *Word) {
		r, g, b, a := parseHexColor(hex)
		w.color = color.RGBA{R: r, G: g, B: b, A: a}
	}
}

func randRotate(rotateRange, orientations int) int {
	return rand.Intn(orientations)*rotateRange/(orientations-1) - rotateRange>>1
}

func Padding(padding int) WordOption {
	return func(w *Word) {
		w.padding = padding
	}
}

func Rotate(rotate int) WordOption {
	return func(w *Word) {
		w.rotate = rotate
	}
}

func LoadFont(filepath string) (f *opentype.Font, err error) {
	fontBytes, _ := os.ReadFile(filepath)
	return opentype.Parse(fontBytes)
}

func NewWord(text string, size int, opts ...WordOption) *Word {
	w := &Word{
		text:   text,
		size:   size,
		color:  nil,
		rotate: noRotated,
	}
	for _, o := range opts {
		o(w)
	}
	return w
}
