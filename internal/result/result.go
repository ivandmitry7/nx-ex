package result

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Date struct {
	Beg int64 `json:"beg"`
	End int64 `json:"end"`
}

type Coords struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
	raw string
}

type Area struct {
	Type   string   `json:"type"`
	Coords []Coords `json:"coords"`
	Radius float32  `json:"radius"`
	raw    string
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
		as = as + fmt.Sprintf("[%s], ", a.raw)
	}
	as = strings.TrimSuffix(as, ", ")

	id := fmt.Sprintf(
		"%s#%s:%s...%s[%s]",
		r.Source, r.Reason,
		time.Unix(r.Date.Beg, 0).UTC().Format(time.RFC3339),
		time.Unix(r.Date.End, 0).UTC().Format(time.RFC3339),
		as,
	)

	return id
}

func (r *Result) PushCoords(radius string, coords []string) {
	a := Area{}
	if radius == "" {
		a.Radius = 0
		a.Type = "polygon"
	} else {
		a.Type = "circle"
		if ml, err := strconv.ParseFloat(radius, 32); err == nil {
			a.Radius = float32(ml * 1.852)
		} else {
			fmt.Printf("unable to convert radius value %q\n", radius)
			a.Radius = 0
		}
	}

	cs := ""
	for _, c := range coords {
		a.Coords = append(a.Coords, Coords{raw: c})
		cs += fmt.Sprintf("[%s], ", c)
	}
	a.raw = fmt.Sprintf("%s/%s[%s]", a.Type, radius, strings.TrimSuffix(cs, ", "))
	r.Area = append(r.Area, a)
}

func (r *Result) Commit() {
	hash := md5.Sum([]byte(r.makeID()))
	r.Hash = hex.EncodeToString(hash[:])
}
