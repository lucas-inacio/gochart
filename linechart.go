package gochart

import (
	"fmt"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	"github.com/fogleman/gg"
)

type LineChart struct {
	BaseChart
}

func NewLineChart(w, h int) *LineChart {
	l := &LineChart{BaseChart: BaseChart{}}
	l.Screen = image.NewRGBA(image.Rect(0, 0, w, h))
	l.margin = 50
	l.ExtendBaseWidget(l)

	l.labelCallback = func (x, y float64) string {
		return fmt.Sprintf("%.3v, %.3v", x, y)
	}

	return l
}

func (l *LineChart) CreateRenderer() fyne.WidgetRenderer {
	renderer := &ChartRenderer{chart: l}
	render := canvas.NewRaster(renderer.draw)
	renderer.render = render
	renderer.objects = []fyne.CanvasObject{render}
	return renderer
}

func (c *LineChart) BeginDraw(dc *gg.Context, x, y float64) {
	dc.SetRGB(0, 0, 0)
	dc.MoveTo(x, y)
}

func (c *LineChart) DrawPoint(dc *gg.Context, x, y, yOffset float64, pointCount int, rect image.Rectangle) {
	dc.LineTo(x, y)
}

func (c *LineChart) EndDraw(dc *gg.Context) {
	dc.Stroke()
}