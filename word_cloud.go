package wordcloud

import (
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"image/color"
	"math"
	"math/rand"
	"strings"
)

const dpi = 72
const cm = 1<<16 - 1

type pixel struct {
	width, height int
	ow, oh        int
	yOff1, yOff2  int
	xOff1, xOff2  int
	rotate        int
	padding       int
	bits          []uint64
}

type WordCloud struct {
	calDc, dc                *gg.Context
	width, height            int
	angleRange, orientations int
	font                     *sfnt.Font
}

func NewWordCloud(width, height int) *WordCloud {
	s := 1 << 12
	cloud := &WordCloud{height: height, width: width, angleRange: 0, orientations: 0}
	cloud.calDc = gg.NewContext(s, s)
	cloud.dc = gg.NewContext(width, height)

	return cloud
}

func (_this *WordCloud) SetFont(font *sfnt.Font) {
	_this.font = font
}

func (_this *WordCloud) SetFontFile(file string) (err error) {
	_this.font, err = LoadFont(file)
	return
}

func (_this *WordCloud) SetRotate(angleRange, orientations int) {
	_this.angleRange = min(max(angleRange, 30), 120)
	_this.orientations = max(orientations, 2)
}

func MeasureString(dc *gg.Context, text string) (int, int) {
	w, h := dc.MeasureString(text)
	h += h / 2
	return int(w), int(h)
}

// faceOpts
// size    float64      // Size is the font size in points
//
//	dpi     float64      // DPI is the dots per inch resolution
//	hinting font.Hinting // Hinting selects how to quantize a vector font's glyph nodes
func faceOpts(size, dpi float64, hinting font.Hinting) *opentype.FaceOptions {
	return &opentype.FaceOptions{Size: size, DPI: dpi, Hinting: hinting}
}

func (_this *WordCloud) wordRotate(w *Word) int {
	if w.rotate != noRotated {
		return w.rotate
	}
	if _this.angleRange == 0 {
		return 0
	}
	return randRotate(_this.angleRange, _this.orientations)
}

func rotatedDraw(dc *gg.Context, x, y, radians float64, f func(dc *gg.Context)) {
	dc.Translate(x, y)
	dc.Rotate(radians)
	f(dc)
	dc.Rotate(-radians)
	dc.Translate(-x, -y)
}

func transformPoint(dc *gg.Context, x, y int) (int, int) {
	_x, _y := dc.TransformPoint(float64(x), float64(y))
	return int(_x), int(_y)
}

func paddingTransformDraw(dc *gg.Context, x, y, width, height, padding int) {
	yMin := 0
	img := dc.Image()
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			_x, _y := transformPoint(dc, x+j, y-i)
			if _x < 0 || _y < 0 {
				continue
			}
			if isColored(img.At(_x, _y)) {
				yMin = min(yMin, -i)
			}
		}
	}
	dc.SetLineWidth(float64(padding))
	w, h := width+padding, y-yMin+1+padding
	dc.DrawRectangle(float64(x-padding>>1), float64(yMin-padding>>1), float64(w), float64(h))
	dc.Stroke()
}

func paddingDraw(dc *gg.Context, x, y, width, height, padding int) {
	yMin := y
	img := dc.Image()
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			_x, _y := x+j, y-i
			if _x < 0 || _y < 0 {
				continue
			}
			if isColored(img.At(_x, _y)) {
				yMin = min(yMin, _y)
			}
		}
	}
	dc.SetLineWidth(float64(padding))
	w, h := width+padding, y-yMin+1+padding
	dc.DrawRectangle(float64(x-padding>>1), float64(yMin-padding>>1), float64(w), float64(h))
	dc.Stroke()
}

func isColored(color color.Color) bool {
	r, g, b, a := color.RGBA()
	return r != cm || g != cm || b != cm || a != cm
}

func setFontFace(dc *gg.Context, f *sfnt.Font, size float64) {
	if f == nil {
		panic("font is empty.")
	}
	face, err := opentype.NewFace(f, faceOpts(size, dpi, font.HintingNone))
	if err != nil {
		panic(err)
	}
	dc.SetFontFace(face)
}

func (_this *WordCloud) calPixels(words []*Word) []*pixel {
	cw, ch := _this.calDc.Width(), _this.calDc.Height()
	dc := _this.calDc
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	top, left, maxH := 0, 0, 0
	var ret []*pixel
	for i := 0; i < len(words); i++ {
		w := words[i]
		setFontFace(dc, _this.font, float64(w.size))
		ow, oh := MeasureString(dc, w.text)

		// set padding
		if w.padding > 0 {
			ow, oh = ow+w.padding<<1, oh+w.padding<<1
		}

		width, height := ow, oh
		rotate := _this.wordRotate(w)
		sin, cos := 0.0, 1.0
		if rotate != 0 {
			radians := gg.Radians(float64(rotate))
			sin, cos = math.Sin(radians), math.Cos(radians)
			_sw := math.Abs(sin * float64(ow))
			_cw := math.Abs(cos * float64(ow))
			_sh := math.Abs(sin * float64(oh))
			_ch := math.Abs(cos * float64(oh))
			_w, _h := int(_cw+_sh+1), int(_ch+_sw+1)
			if left+_w >= cw {
				top += maxH
				left = 0
			}
			if _h+top >= ch {
				break
			}
			_x, _y := math.Max(float64(left)+math.Ceil(float64(w.padding)*(math.Abs(sin)+math.Abs(cos))), math.Ceil(float64(left)-sin*float64(oh)*.9)), math.Ceil(float64(oh*9/10)*cos)+float64(top)
			if sin < 0 {
				_y = math.Ceil((_ch+_sw)*9/10) + float64(top)
			}
			rotatedDraw(dc, _x, _y, radians, func(dc *gg.Context) {
				dc.DrawString(w.text, 0, 0)
				if w.padding > 0 {
					paddingTransformDraw(dc, 0, 0, ow-w.padding<<1, oh, w.padding)
				}
			})
			width, height = _w, _h
		} else {
			if left+width >= cw {
				top += maxH
				left = 0
			}
			if height+top >= ch {
				break
			}
			_x, _y := float64(left+w.padding), float64(height-height/10+top)
			dc.DrawString(w.text, _x, _y)
			if w.padding > 0 {
				paddingDraw(dc, int(_x), int(_y), ow-w.padding<<1, oh, w.padding)
			}
		}

		img := dc.Image()
		bits := make([]uint64, ((width-1)>>6+1)*height)
		idx := 0

		yOff1, yOff2 := 0, -1
		xOff1, xOff2 := _this.width-1, 0
		rowBits := (width-1)>>6 + 1
		for y := 0; y < height; y++ {
			arr, ai := make([]uint64, rowBits), 0
			for x := 0; x < width; {
				var bit uint64
				for k := 0; k < 64 && x < width; k++ {
					if isColored(img.At(x+left, y+top)) {
						bit |= 1 << k
						xOff1 = min(xOff1, x)
						xOff2 = max(xOff2, x)
					}
					x++
				}
				arr[ai] = bit
				ai++
				if bit > 0 {
					yOff2 = y
				}
			}
			if yOff2 != -1 {
				copy(bits[idx:], arr)
				idx += len(arr)
			} else {
				yOff1++
			}
		}

		ret = append(ret, &pixel{width: width, height: height, ow: ow, oh: oh, yOff1: yOff1, yOff2: yOff2, xOff1: xOff1, xOff2: xOff2, rotate: rotate, padding: w.padding, bits: bits[:(yOff2-yOff1+1)*rowBits]})
		if height > maxH {
			maxH = height
		}
		left += width + 10
	}
	return ret
}

func setBorder(dc *gg.Context, left, top, width, height int) {
	dc.SetRGB(0, 1, 1)
	for j := 0; j < width; j++ {
		dc.SetPixel(j+left, top)
		dc.SetPixel(j+left, top+height)
	}
	for j := 0; j < height; j++ {
		dc.SetPixel(left, j)
		dc.SetPixel(left+width, j)
	}
}

// i is the index of the point
// e is the aspect ratio
// step is the angle increment
// a is the distance from the center point
// b is the pitch
func position(i int, e, step, a, b float64) (dx, dy float64) {
	s := float64(i) * step
	dx, dy = e*(a+b*s)*math.Cos(s), (a+b*s)*math.Sin(s)
	return
}

func createBoard(w, h int) []uint64 {
	return make([]uint64, (w*h-1)>>6+1)
}

func (_this *WordCloud) Render(words []*Word, filePath string) error {
	_, err := _this.RenderAll(words, filePath, false)
	return err
}

func (_this *WordCloud) RenderAll(words []*Word, filePath string, renderOthers bool) ([]uint64, error) {
	w, h := _this.width, _this.height
	dc := _this.dc
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	pixels := _this.calPixels(words)
	for len(pixels) < len(words) {
		pixels = append(pixels, _this.calPixels(words[len(pixels):])...)
	}

	board := createBoard(w, h)
	for i := 0; i < len(pixels); i++ {
		_this.place(dc, board, words[i], pixels[i])
	}

	if renderOthers {
		path := ""
		if idx := strings.LastIndex(filePath, "/"); idx != -1 {
			path = filePath[:idx] + "/"
		}
		err := _this.calDc.SavePNG(path + "cal_board.png")
		if err != nil {
			panic(err)
		}
		_this.renderBoard(board, path+"board.png")
	}
	return board, dc.SavePNG(filePath)
}

// renderBoard for test
func (_this *WordCloud) renderBoard(board []uint64, filePath string) {
	dc := gg.NewContext(_this.width, _this.height)
	dc.SetRGB(0, 1, 0)
	for y := 0; y < _this.height; y++ {
		for x := 0; x < _this.width; x++ {
			_p := y*_this.width + x
			_idx, _offset := _p>>6, _p&63
			if board[_idx]&(1<<_offset) != 0 {
				dc.SetPixel(x, y)
			}
		}
	}
	err := dc.SavePNG(filePath)
	if err != nil {
		panic(err)
	}
}

func (_this *WordCloud) place(dc *gg.Context, board []uint64, word *Word, px *pixel) {
	dt := 1
	if rand.Intn(2) == 1 {
		dt = -dt
	}
	x, y, success := 0, 0, false
	for i := 0; i < 10 && !success; i++ {
		x, y, success = _this.placeByDt(dt, board, px)
		dt = -dt
	}
	if !success {
		panic("too many words, please increase the width and height of the canvas.")
	}
	setFontFace(dc, _this.font, float64(word.size))
	if word.color != nil {
		dc.SetColor(word.color)
		defer func() {
			dc.SetRGB(0, 0, 0)
		}()
	}
	if px.rotate != 0 {
		radians := gg.Radians(float64(px.rotate))
		sin := math.Sin(radians)
		cos := math.Cos(radians)
		_x, _y := math.Max(float64(x)+math.Ceil(float64(px.padding)*(math.Abs(sin)+math.Abs(cos))), math.Ceil(float64(x)-sin*float64(px.oh)*.9)), math.Ceil(float64(px.oh*9/10)*cos)+float64(y)
		if sin < 0 {
			_y = math.Ceil(float64(px.height)*9/10) + float64(y)
		}
		rotatedDraw(dc, _x, _y, radians, func(dc *gg.Context) {
			dc.DrawString(word.text, 0, 0)
		})
	} else {
		dc.DrawString(word.text, float64(x+px.padding), float64(y+px.height-px.height/10))
	}
	_this.record(board, x, y, px)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (_this *WordCloud) placeByDt(dt int, board []uint64, px *pixel) (x, y int, success bool) {
	const (
		step = 0.1
		b    = 1
		a    = 0
	)
	w, h := _this.width, _this.height
	maxDelta := math.Sqrt(float64(w*w + h*h))
	e := float64(w) / float64(h)
	t := 0
	// x, y = w*(rand.Intn(1001)+500)/2000, h*(rand.Intn(1001)+500)/2000
	x, y = w>>1, h>>1
	for t += dt; ; t += dt {
		dx, dy := position(t, e, step, a, b)
		if math.Min(math.Abs(dx), math.Abs(dy)) > maxDelta {
			return
		}
		_x, _y := x+int(dx)-px.xOff1-(px.xOff2-px.xOff1)>>1, y+int(dy)-px.yOff1-(px.yOff2-px.yOff1)>>1
		if !_this.checkPoint(_x, _y, px) {
			continue
		}
		if !_this.isCollide(board, _x, _y, px) {
			x, y = _x, _y
			success = true
			return
		}
	}
}

func (_this *WordCloud) checkPoint(x, y int, px *pixel) bool {
	if min(x, x+px.yOff1-px.padding) < 0 || x+px.xOff2+1+px.padding >= _this.width {
		return false
	}

	if y+px.yOff1-px.padding < 0 || y+px.yOff2+1+px.padding >= _this.height {
		return false
	}

	return true
}

func (_this *WordCloud) record(board []uint64, x, y int, px *pixel) {
	rowBits := (px.width-1)>>6 + 1
	for i := px.yOff1; i <= px.yOff2; i++ {
		idx := rowBits * (i - px.yOff1)
		for j := 0; j < px.width; j, idx = j+64, idx+1 {
			bit := px.bits[idx]
			if bit == 0 {
				continue
			}
			_y, _x := i+y, j+x
			_p := _y*_this.width + _x
			_idx, _offset := _p>>6, _p&63
			board[_idx] |= bit << _offset
			left := bit >> (64 - _offset)
			if left != 0 {
				board[_idx+1] |= bit >> (64 - _offset)
			}
		}
	}
}

func (_this *WordCloud) isCollide(board []uint64, x, y int, px *pixel) bool {

	rowBits := (px.width-1)>>6 + 1

	for i := px.yOff1; i <= px.yOff2; i++ {
		idx := rowBits * (i - px.yOff1)
		for j := 0; j <= px.xOff2; j += 64 {
			bit := px.bits[idx]
			_y, _x := i+y, j+x
			_p := _y*_this.width + _x
			_idx, _offset := _p>>6, _p&63
			left := bit >> (64 - _offset)
			if board[_idx]&(bit<<_offset) != 0 || left != 0 && board[_idx+1]&(bit>>(64-_offset)) != 0 {
				return true
			}
			idx++
		}
	}
	return false
}
