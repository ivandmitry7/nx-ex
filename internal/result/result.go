package result

type Date struct {
	Beg int64 `json:"beg"`
	End int64 `json:"end"`
}

type Coords struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

type Area struct {
	Type   string   `json:"type"`
	Coords []Coords `json:"coords"`
	Radius float32  `json:"radius"`
}

type Result struct {
	Hash   string `json:"hash"`
	Source string `json:"source"`
	Reason string `json:"reason"`
	Date   Date   `json:"date"`
	Area   []Area `json:"area"`
}
