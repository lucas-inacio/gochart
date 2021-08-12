package main

import (
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/lucas-inacio/gochart"
)

const (
	NUM_SAMPLES = 40
	AMPLITUDE   = 10
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Canvas")
	x := []float64{}
	y := []float64{}

	for i := 0; i < NUM_SAMPLES; i++ {
		x = append(x, 2 * math.Pi * float64(i) / float64(NUM_SAMPLES))
		y = append(y, AMPLITUDE * math.Sin(2.0 * math.Pi * float64(i) / float64(NUM_SAMPLES)))
	}

	err := gochart.SetFont("./assets/fonts/Molengo-Regular.ttf", 11)
	if err != nil {
		panic(err)
	}

	chart := gochart.NewBarChart(640, 480)
	chart.SetGrowY(true)
	chart.SetData(x, y)

	// Call refresh on the chart so it can update it's cursor position
	// and show it
	animation := fyne.NewAnimation(33 * time.Millisecond, func(_ float32) {
		chart.Refresh()
	})
	animation.RepeatCount = fyne.AnimationRepeatForever
	animation.Start()

	myWindow.SetContent(chart)
	myWindow.Resize(fyne.NewSize(640, 480))
	myWindow.ShowAndRun()
}