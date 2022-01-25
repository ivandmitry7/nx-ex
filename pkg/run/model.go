package run

import (
	"container/list"
	"math"
	"time"
)

// Options File processing options
type Options struct {
	Freq   int               // Central frequency, Hz
	OutDir string            // Results dir name
	Batch  uint              // Read buffer size, ms
	Params map[string]string // Demodulator parameters
}

// Position Current processing position
type Position struct {
	Start   time.Time
	Stop    time.Time
	Sample  int64
	Tick    time.Duration
	Percent int
}

type ResultStat struct {
	HasData bool
	Text    int
	Bits    int
	Phase   int
	Ale     int
}

type Result struct {
	Position Position
	Log      string
	Freq     int
	Stat     ResultStat
}

type Model struct {
	Path       string
	SrcInfo    string
	Duration   time.Duration
	Opt        Options
	Results    list.List
	Progress   Position
	Total      ResultStat
	SampleRate int
}

func (m *Model) PushResults(count int, res Results) {
	m.Progress.Sample += int64(count)
	m.Progress.Tick += time.Duration(count*1000/m.SampleRate) * time.Millisecond
	m.Progress.Percent =
		int(math.Round(float64(m.Progress.Tick.Milliseconds()*100.0) / float64(m.Duration.Milliseconds())))

	r := Result{Position: m.Progress}
	r.Stat.HasData = res.HasData
	r.Freq = res.Freq
	r.Log = res.Log

	// @formatter:off
	r.Stat.Text = len(res.Text)
	r.Stat.Bits = len(res.Bits)
	r.Stat.Phase = len(res.Phase)
	r.Stat.Ale = len(res.Ale)

	m.Total.Bits += r.Stat.Bits
	m.Total.Phase += r.Stat.Phase
	m.Total.Text += r.Stat.Text
	m.Total.Ale += r.Stat.Ale
	// @formatter:on

	m.Results.PushBack(r)
	if m.Results.Len() > 5 {
		m.Results.Remove(m.Results.Front())
	}
}
