package main

import (
	"time"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/widget"
)

func main() {
	widget.Init()
	defer widget.Fini()

	look.Set("default")
	maxval := 200 * 1024 * 1024

	width, _ := widget.Size()
	gauge := widget.NewProgressGauge(0, 0, width/2, 1)
	gauge.Start(float64(maxval))
	ticker := time.NewTicker(10 * time.Millisecond)

	const n = 50 * 1024 * 1024 / 100 // 50Mb/s
	progress := 0
	for {
		progress += n
		if progress > maxval {
			gauge.Finish()
			break
		}
		gauge.Update(float64(n))
		gauge.Draw()
		widget.Show()
		<-ticker.C
	}
	widget.PollEvent()
	widget.PollEvent()
}
