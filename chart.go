package gochart

import (
	"fmt"
	"image"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/fogleman/gg"

	"golang.org/x/image/font"
)

var currentFont font.Face

func SetFont(path string, points float64) error {
	fontFace, err := gg.LoadFontFace(path, points)
	if err != nil {
		return err
	}
	currentFont = fontFace
	return nil
}

type Chart interface {
	fyne.Widget
	SetGrowY(grow bool)
	SetData(x, y []float64)

	BeginDraw(dc *gg.Context, x, y float64)
	DrawPoint(dc *gg.Context, x, y, yOffset float64, pointCount int, rect image.Rectangle)
	EndDraw(dc *gg.Context)

	DrawCursor(dc *gg.Context, ratioX float64)
	DrawFrame(dc *gg.Context, ratioX, ratioY, yOffset, height, width, xSpan, ySpan float64)

	ResizeIfNeeded(w, h int)

	GetData() ([]float64, []float64)
	GetScreen() image.Image
	GetGrowY() bool
	GetMargin() float64
	GetMouse()  (float32, float32)
	GetShowLabel() bool
}

type ChartRenderer struct {
	chart Chart
	render *canvas.Raster
	objects []fyne.CanvasObject
}

func (l *ChartRenderer) draw(w, h int) image.Image {
	l.chart.ResizeIfNeeded(w, h)

	screen := l.chart.GetScreen()
	margin := l.chart.GetMargin()
	xData, yData := l.chart.GetData()
	autoGrowY := l.chart.GetGrowY()

	width := float64(w) - 2 * margin
	height := float64(h) - 2 * margin

	// Scales x span acordingly
	xSpan := float64(xData[len(xData) - 1] - xData[0])
	ratioX := width / xSpan

	// Wether all y values should be normalized in relation to the canvas
	// height and use all height
	// The values below are used to make sure the line won't overflow
	ratioY := 1.0
	yOffset := 0.0
	if autoGrowY {
		offset, modulus := getParameters(yData)
		ratioY = height / modulus
		yOffset = offset
	}

	dc := gg.NewContextForRGBA(screen.(*image.RGBA))
	dc.SetRGB(0.84, 0.72, 0.40)
	dc.Clear()
	dc.SetFontFace(currentFont)

	// Draws the path. Y is subtracted from the canvas height to follow normal
	// convention (bigger number means up)
	_, zeroY := computeCoords(0, 0, ratioX, ratioY, yOffset, height)
	x, y := computeCoords(xData[0], yData[0], ratioX, ratioY, yOffset, height)
	dc.Push()
	l.chart.BeginDraw(dc, x + margin, y + margin)
	for index := 0; index < len(yData); index++ {
		x, y = computeCoords(xData[index], yData[index], ratioX, ratioY, yOffset, height)
		rect := image.Rect(int(margin), int(margin), int(width + margin), int(height + margin))
		l.chart.DrawPoint(dc, x + margin, y + margin, zeroY + margin, len(yData), rect)
	}
	l.chart.EndDraw(dc)
	dc.Pop()
	l.chart.DrawFrame(dc, ratioX, ratioY, yOffset, height, width, xSpan, height / ratioY)
	l.chart.DrawCursor(dc, ratioX)

	return screen
}

func (l *ChartRenderer) Destroy() {

}

func (l *ChartRenderer) Layout(size fyne.Size) {
	l.render.Resize(size)
}

func (l *ChartRenderer) MinSize() fyne.Size {
	// return l.render.MinSize()
	return fyne.NewSize(250, 250)
}

func (l *ChartRenderer) Objects() []fyne.CanvasObject {
	return l.objects
}

func (l *ChartRenderer) Refresh() {
	canvas.Refresh(l.render)
}

type BaseChart struct {
	widget.BaseWidget
	Screen image.Image
	xData, yData []float64
	autoGrowY bool
	margin float64

	// Interactivity
	showLabel bool
	mouseX, mouseY float32
	labelCallback func (x, y float64) string
}

func NewBaseChart(w, h int) *BaseChart {
	l := &BaseChart{}
	l.Screen = image.NewRGBA(image.Rect(0, 0, w, h))
	l.margin = 50
	l.ExtendBaseWidget(l)

	l.labelCallback = func (x, y float64) string {
		return fmt.Sprintf("%.3v, %.3v", x, y)
	}

	return l
}

func (l *BaseChart) CreateRenderer() fyne.WidgetRenderer {
	renderer := &ChartRenderer{chart: l}
	render := canvas.NewRaster(renderer.draw)
	renderer.render = render
	renderer.objects = []fyne.CanvasObject{render}
	return renderer
}

func (c *BaseChart) SetGrowY(grow bool) {
	c.autoGrowY = grow
}

func (c *BaseChart) SetData(x, y []float64) {
	c.xData = x
	c.yData = y
}

func (c *BaseChart) DrawCursor(dc *gg.Context, ratioX float64) {
	if c.showLabel {
		index, x := c.findClosestX(float64(c.mouseX), ratioX)
		dc.SetRGB(0.6, 0, 0)
		dc.DrawLine(x, c.margin, x, float64(c.Screen.Bounds().Dy()) - c.margin)
		dc.Stroke()
		
		text := c.labelCallback(c.xData[index], c.yData[index])
		left := float64(c.mouseX)
		top := float64(c.mouseY)
		w, h := dc.MeasureString(text)
		dc.SetRGB(1, 1, 1)
		dc.DrawRectangle(left - w / 4, top - h - h / 2, w * 1.5, h * 2)
		dc.Fill()
		
		dc.SetRGB(0.0, 0, 0.0)
		dc.DrawString(text, float64(c.mouseX), float64(c.mouseY))
		dc.Stroke()
	}
}

func (c *BaseChart) DrawFrame(dc *gg.Context, ratioX, ratioY, yOffset, height, width, xSpan, ySpan float64) {
	screenHeight := float64(c.Screen.Bounds().Dy())
	screenWidth := float64(c.Screen.Bounds().Dx())
	dc.SetRGB(1, 1, 1)
	// Top border
	dc.MoveTo(0, 0)
	dc.LineTo(screenWidth, 0)
	dc.LineTo(screenWidth, c.margin)
	dc.LineTo(0, c.margin)
	dc.Fill()
	// Bottom border
	dc.MoveTo(0, screenHeight)
	dc.LineTo(0, screenHeight - c.margin)
	dc.LineTo(screenWidth, screenHeight - c.margin)
	dc.LineTo(screenWidth, screenHeight)
	dc.Fill()
	// Left border
	dc.MoveTo(0, c.margin)
	dc.LineTo(c.margin, c.margin)
	dc.LineTo(c.margin, screenHeight - c.margin)
	dc.LineTo(0, screenHeight - c.margin)
	dc.Fill()
	// Right border
	dc.MoveTo(screenWidth - c.margin, c.margin)
	dc.LineTo(screenWidth, c.margin)
	dc.LineTo(screenWidth, screenHeight - c.margin)
	dc.LineTo(screenWidth - c.margin, screenHeight - c.margin)
	dc.Fill()

	dc.SetRGB(0, 0, 0)
	count := int(screenHeight * 2 /  (screenHeight - height))
	stepY := ySpan / float64(count)
	start := -count / 2
	for i := start; i < count + start; i++ {
		y := height - (stepY * float64(i) - yOffset) * ratioY
		dc.DrawString(fmt.Sprintf("%.3v", stepY * float64(i)), c.margin / 6, y + c.margin)
	}

	count2 := int(screenWidth /  (screenWidth - width))
	stepX := xSpan / float64(count2)
	for i := 1; i < count2; i++ {
		x := stepX * float64(i) * ratioX
		dc.DrawString(fmt.Sprintf("%.3v", stepX * float64(i)), x + c.margin, screenHeight - c.margin / 2)
	}
}

func (c *BaseChart) ResizeIfNeeded(w, h int) {
	width := c.Screen.Bounds().Dx()
	height := c.Screen.Bounds().Dy()
	if w != width || h != height {
		c.Screen = image.NewRGBA(image.Rect(0, 0, w, h))
	}
}

// Events
// A minimum implementation is provided so the cursor is is tracked when it's over
// the chart area
func (c *BaseChart) Dragged(event *fyne.DragEvent) {
	c.showLabel = c.inside(int(event.Position.X), int(event.Position.Y))
	c.mouseX = event.Position.X
	c.mouseY = event.Position.Y
}

func (c *BaseChart) DragEnd() {
	c.showLabel = false
}


func (c *BaseChart) BeginDraw(dc *gg.Context, x, y float64) {

}

func (c *BaseChart) DrawPoint(dc *gg.Context, x, y, yOffset float64, pointCount int, rect image.Rectangle) {

}

func (c *BaseChart) EndDraw(dc *gg.Context) {
	
}


func (c *BaseChart) inside(x, y int) bool {
	margin := int(c.margin)
	width := c.Screen.Bounds().Dx()
	height := c.Screen.Bounds().Dy()
	if x < width - margin && x > margin && y < height - margin && y > margin {
		return true
	}
	return false
}

func (c *BaseChart) findClosestX(xTarget, ratioX float64) (int, float64) {
	delta := float64(c.Screen.Bounds().Dx())
	xResult := 0.0
	indexResult := 0
	for index, value := range c.xData {
		x := ratioX * value + c.margin
		newDelta := math.Abs(x - xTarget)
		if newDelta < delta {
			delta = newDelta
			xResult = x
			indexResult = index
		}
	}
	
	return indexResult, xResult
}

func (c *BaseChart) GetData() ([]float64, []float64) {
	return c.xData, c.yData
}

func (c *BaseChart) GetScreen() image.Image {
	return c.Screen
}

func (c *BaseChart) GetGrowY() bool {
	return c.autoGrowY
}

func (c *BaseChart) GetMargin() float64 {
	return c.margin
}

func (c *BaseChart) GetMouse()  (float32, float32) {
	return c.mouseX, c.mouseY
}

func (c *BaseChart) GetShowLabel() bool {
	return c.showLabel
}

func getParameters(values []float64) (float64, float64) {
	modulus := 0.0
	offset := 0.0
	for _, value := range values {
		mag := math.Abs(value)
		if mag > modulus {
			modulus = mag
		}

		if value < offset {
			offset = value
		}
	}
	return offset, modulus + math.Abs(offset)
}

func computeCoords(x, y, ratioX, ratioY, yOffset, height float64) (float64, float64) {
	return ratioX * x, height - ratioY * (y - yOffset)
}

