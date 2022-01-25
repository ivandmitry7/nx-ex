package run

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestModel_PushResults(t *testing.T) {
	m := Model{
		Duration:   1 * time.Second,
		SampleRate: 8000,
	}

	r := Results{
		HasData: true,
		Freq:    1800,
		Text:    "This is text",
		Bits:    "0000010110010100",
		Ale:     "This is ale",
		Log:     "This is log",
		Phase:   [][2]float32{{1.0, 1.1}, {1.2, 1.3}},
	}
	m.PushResults(800, r)
	assert.Equal(t, m.Progress.Sample, int64(800))
	assert.Equal(t, m.Progress.Tick, 100*time.Millisecond)
	assert.Equal(t, m.Progress.Percent, 10)
	assert.Equal(t, m.Total.Bits, 16)
	assert.Equal(t, m.Total.Text, 12)
	assert.Equal(t, m.Total.Phase, 2)
	assert.Equal(t, m.Total.Ale, 11)
	assert.Equal(t, m.Results.Len(), 1)
}
