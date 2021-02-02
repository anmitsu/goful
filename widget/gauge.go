package widget

import (
	"fmt"
	"strings"
	"time"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/util"
)

// ProgressGauge is a progress bar changing with time.
type ProgressGauge struct {
	*Window
	maxval    float64
	curval    float64
	starttime time.Time
}

// NewProgressGauge returns a new progress gauge specified coordinates and sizes.
func NewProgressGauge(x, y, width, height int) *ProgressGauge {
	return &ProgressGauge{
		Window: NewWindow(x, y, width, height),
		maxval: 1,
		curval: 1,
	}
}

// Start the timer from now.
func (b *ProgressGauge) Start(maxval float64) {
	b.curval = 0
	b.maxval = maxval
	b.starttime = time.Now()
}

// Update the current value.
func (b *ProgressGauge) Update(value float64) {
	b.curval += value
}

// Finish the progress.
func (b *ProgressGauge) Finish() {
	b.curval = b.maxval
}

// IsFinished reports whether finished.
func (b *ProgressGauge) IsFinished() bool {
	return b.curval == b.maxval
}

func (b *ProgressGauge) elapsedTime() time.Duration {
	return time.Since(b.starttime)
}

func (b *ProgressGauge) eta() string {
	if b.curval == 0 {
		return "unknown"
	}
	elapse := float64(b.elapsedTime())
	remain := time.Duration(elapse*b.maxval/b.curval - elapse)
	return fmtString(remain)
}

func fmtString(d time.Duration) string {
	u := uint64(d)
	switch {
	case u < uint64(time.Minute):
		return fmt.Sprintf("00:%02.0f", d.Seconds())
	case u < uint64(time.Hour):
		s := u % uint64(time.Minute)
		d0 := time.Duration(s)
		d1 := time.Duration(u - s)
		return fmt.Sprintf("%02d:%02.0f", int(d1.Minutes()), d0.Seconds())
	default:
		m := u % uint64(time.Hour)
		s := u % uint64(time.Minute)
		d0 := time.Duration(s)
		d1 := time.Duration(m)
		d2 := time.Duration(u - m - s)
		return fmt.Sprintf("%02.0f:%02d:%02.0f", d2.Hours(), int(d1.Minutes()), d0.Seconds())
	}
}

func (b *ProgressGauge) bps() string {
	bps := b.curval / b.elapsedTime().Seconds()
	return fmt.Sprintf("%sB/s", util.FormatSize(int64(bps)))
}

// Draw the progress gauge from start time to now and estimated time of arrival.
func (b *ProgressGauge) Draw() {
	x, y := b.LeftTop()

	rate := float64(b.curval / b.maxval)
	percent := fmt.Sprintf("%3d%s|", int(rate*100), "%")
	x = SetCells(x, y, percent, look.Default())

	etaBps := fmt.Sprintf("| %s [%s]", b.eta(), b.bps())
	max := b.Width() - len(percent) - len(etaBps)
	current := int(rate * float64(max))
	if current > max {
		current = max
	}
	s := strings.Repeat(">", current)
	x = SetCells(x, y, s, look.Progress())
	s = strings.Repeat("-", max-current)
	x = SetCells(x, y, s, look.Default())
	SetCells(x, y, etaBps, look.Default())
}
