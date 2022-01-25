package run

// Results Demodulation results
type Results struct {
	HasData bool
	Freq    int
	Text    string
	Bits    string
	Ale     string
	Log     string
	Phase   [][2]float32
	Params  []map[string]string
	Ticks   int64
}

// ResultsArray Set of demodulation results
type ResultsArray []Results
