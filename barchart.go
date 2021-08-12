package gochart

import (
	"fmt"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	"github.com/fogleman/gg"
)

type BarChart struct {
	BaseChart
}

func NewBarChart(w, h int) *BarChart {
	l := &BarChart{BaseChart: BaseChart{}}
	l.Screen = image.NewRGBA(image.Rect(0, 0, w, h))
	l.margin = 50
	l.ExtendBaseWidget(l)

	l.labelCallback = func (x, y float64) string {
		return fmt.Sprintf("%.3v, %.3v", x, y)
	}

	return l
}

func (l *BarChart) CreateRenderer() fyne.WidgetRenderer {
	renderer := &ChartRenderer{chart: l}
	render := canvas.NewRaster(renderer.draw)
	renderer.render = render
	renderer.objects = []fyne.CanvasObject{render}
	return renderer
}

func (c *BarChart) BeginDraw(dc *gg.Context, x, y float64) {
	dc.SetRGB(0, 0.3, 0.5)
}

func (c *BarChart) DrawPoint(dc *gg.Context, x, y, yOffset float64, pointCount int, rect image.Rectangle) {
	width := float64(rect.Dx()) / float64(pointCount)
	dc.DrawRectangle(x - width / 2, y, width, yOffset - y)
}

func (c *BarChart) EndDraw(dc *gg.Context) {
	dc.Fill()
}