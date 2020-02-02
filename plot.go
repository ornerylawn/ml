package ml

import (
	"image"
	"image/color"
	"math"

	"github.com/ornerylawn/linear"
)

type Color struct {
	R, G, B, A float64 // not "pre-multiplied"
}

// RGBA returns the "pre-multiplied" 8-bit per channel color.
func (c Color) RGBA() (r, g, b, a uint32) {
	return uint32(c.R * c.A * 255.0),
		uint32(c.G * c.A * 255.0),
		uint32(c.B * c.A * 255.0),
		uint32(c.A * 255.0)
}

func ColorFromRaster(c color.Color) Color {
	r, g, b, a := c.RGBA() // "pre-multiplied"
	alpha := float64(a) / 255.0
	return Color{
		(float64(r) / alpha) / 255.0,
		(float64(g) / alpha) / 255.0,
		(float64(b) / alpha) / 255.0,
		alpha,
	}
}

type Point struct {
	X, Y float64
}

type Canvas interface {
	SetStrokeColor(color Color)
	SetFillColor(color Color)
	SetStrokeWidth(width float64)
	SetCircleRadius(radius float64)

	MovePenTo(p Point)
	DrawLineTo(p Point)
	DrawCircle()
	DrawRectangleTo(p Point)

	Clear(color Color)
}

type Raster struct {
	image                     *image.RGBA
	strokeColor, fillColor    Color
	strokeWidth, circleRadius float64
	penPosition               Point
}

func NewRaster(width, height int) *Raster {
	r := &Raster{
		image:        image.NewRGBA(image.Rect(0, 0, width, height)),
		strokeColor:  Color{0, 0, 0, 255},
		fillColor:    Color{0, 0, 0, 0},
		strokeWidth:  1.0,
		circleRadius: 3.0,
	}
	r.Clear(Color{255, 255, 255, 255})
	return r
}

func (r *Raster) SetStrokeColor(color Color)     { r.strokeColor = color }
func (r *Raster) SetFillColor(color Color)       { r.fillColor = color }
func (r *Raster) SetStrokeWidth(width float64)   { r.strokeWidth = width }
func (r *Raster) SetCircleRadius(radius float64) { r.circleRadius = radius }

func (r *Raster) MovePenTo(p Point) { r.penPosition = p }
func (r *Raster) DrawLineTo(p Point) {
	// TODO: difficult!
}
func over(paint, canvas Color) Color {
	// TODO: gamma correction.
	return porterDuff(paint, canvas, paint, 1.0, 1.0, 1.0)
}
func porterDuff(paint, canvas, both Color, pRegion, cRegion, bRegion float64) Color {
	areaPaint := paint.A * (1.0 - canvas.A)
	areaCanvas := canvas.A * (1.0 - paint.A)
	areaBoth := paint.A * canvas.A
	r := areaPaint*paint.R + areaCanvas*canvas.R + areaBoth*both.R
	g := areaPaint*paint.G + areaCanvas*canvas.G + areaBoth*both.G
	b := areaPaint*paint.B + areaCanvas*canvas.B + areaBoth*both.B
	a := areaPaint*pRegion + areaCanvas*cRegion + areaBoth*bRegion
	return Color{r, g, b, a}
}

func (r *Raster) DrawCircle() {
	r.fillCircle()
	r.strokeCircle()
}
func computeCircleAreaOnPixel(x, y int, center Point, radius float64) float64 {
	// TODO
	return 0.0
}

func (r *Raster) fillCircle() {
	xlo := int(math.Floor(float64(r.penPosition.X) - r.circleRadius))
	xhi := int(math.Ceil(float64(r.penPosition.X)+r.circleRadius)) + 1
	ylo := int(math.Floor(float64(r.penPosition.Y) - r.circleRadius))
	yhi := int(math.Ceil(float64(r.penPosition.Y)+r.circleRadius)) + 1

	paintColor := r.fillColor

	for y := ylo; y < yhi; y++ {
		for x := xlo; x < xhi; x++ {
			paintColor.A = computeCircleAreaOnPixel(x, y, r.penPosition, r.circleRadius)
			blendedColor := over(paintColor, ColorFromRaster(r.image.At(x, y)))
			r.image.Set(x, y, blendedColor)
		}
	}
}
func (r *Raster) strokeCircle() {
	// TODO: difficult! (one circle subtracted from another?)
}
func (r *Raster) DrawRectangleTo(p Point) {
	r.fillRectangle(p)
	r.strokeRectangle(p)
}
func computeRectangleAreaOnPixel(x, y int, p1, p2 Point) float64 {
	// TODO
	return 0.0
}
func (r *Raster) fillRectangle(p Point) {
	leftX := r.penPosition.X
	rightX := p.X
	if rightX < leftX {
		leftX, rightX = rightX, leftX
	}
	bottomY := r.penPosition.Y
	topY := p.Y
	if topY < bottomY {
		bottomY, topY = topY, bottomY
	}
	xlo := int(math.Floor(leftX))
	xhi := int(math.Ceil(rightX)) + 1
	ylo := int(math.Floor(bottomY))
	yhi := int(math.Ceil(topY)) + 1

	paintColor := r.fillColor

	for y := ylo; y < yhi; y++ {
		for x := xlo; x < xhi; x++ {
			paintColor.A = computeRectangleAreaOnPixel(x, y, r.penPosition, p)
			blendedColor := over(paintColor, ColorFromRaster(r.image.At(x, y)))
			r.image.Set(x, y, blendedColor)
		}
	}
}
func (r *Raster) strokeRectangle(p Point) {
	// TODO: difficult! (one rect subtracted from another?)
}

func (r *Raster) Clear(color Color) {
	b := r.image.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r.image.Set(x, y, color)
		}
	}
}

func (r *Raster) Image() image.Image { return r.image }

func LinePlot(xs, ys linear.Vector, canvas Canvas) {
	canvas.MovePenTo(Point{xs.Get(0), ys.Get(0)})
	for d := 1; d < xs.Dimension(); d++ {
		canvas.DrawLineTo(Point{xs.Get(d), ys.Get(d)})
	}
}

func ScatterPlot(xs, ys linear.Vector, canvas Canvas) {
	for d := 0; d < xs.Dimension(); d++ {
		canvas.MovePenTo(Point{xs.Get(d), ys.Get(d)})
		canvas.DrawCircle()
	}
}

func Histogram(binHeights linear.Vector, binWidth float64, canvas Canvas) {
	for d := 0; d < binHeights.Dimension(); d++ {
		canvas.MovePenTo(Point{float64(d) * binWidth, 0})
		canvas.DrawRectangleTo(Point{float64(d+1) * binWidth, binHeights.Get(d)})
	}
}
