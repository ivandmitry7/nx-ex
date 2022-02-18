package parser

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type Date struct {
	Beg int64  `json:"beg"`
	End int64  `json:"end"`
	Raw string `json:"-"`
}

type Coords struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

type Area struct {
	Type   string   `json:"type"`
	Coords []Coords `json:"coords"`
	Radius float32  `json:"radius"`
	Raw    string   `json:"-"`
}

type Result struct {
	Hash   string `json:"hash"`
	Source string `json:"source"`
	Reason string `json:"reason"`
	Date   Date   `json:"date"`
	Area   []Area `json:"area"`
}

func (r *Result) makeID() string {
	as := ""
	for _, a := range r.Area {
		as = as + fmt.Sprintf("%s, ", a.Raw)
	}
	as = strings.TrimSuffix(as, ", ")

	id := fmt.Sprintf(
		"%s/%s...%s[%s]",
		r.Reason,
		time.Unix(r.Date.Beg, 0).UTC().Format(time.RFC3339),
		time.Unix(r.Date.End, 0).UTC().Format(time.RFC3339),
		as,
	)

	return id
}

func (r *Result) Commit() {
	hash := md5.Sum([]byte(r.makeID()))
	r.Hash = hex.EncodeToString(hash[:])
}
