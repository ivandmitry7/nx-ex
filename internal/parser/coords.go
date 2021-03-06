package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func coordsToFloat(d string, m string, s string, l string) float32 {
	di, _ := strconv.Atoi(d)
	mi, _ := strconv.Atoi(m)
	si, _ := strconv.Atoi(s)

	res := float32(di) + float32(mi)/60.0 + float32(si)/3600.0
	if l == "S" || l == "W" {
		res = -res
	}

	return res
}

func stringToCoord(cs string) (coord float32, e error) {
	re, err := regexp.Compile(`(\d+)-(\d+)(?:[.,-](\d+))?\s?([NSEW])`)
	if err != nil {
		e = fmt.Errorf("unable to parse coordinate string %q", cs)
		return
	}
	m := re.FindAllStringSubmatch(cs, -1)
	if m == nil {
		e = fmt.Errorf("unable to parse coordinate string %q", cs)
		return
	}
	if len(m[0]) != 5 {
		e = fmt.Errorf("unable to parse coordinate string %q", cs)
		return
	}

	coord = coordsToFloat(m[0][1], m[0][2], m[0][3], m[0][4])
	return
}

func parseCircle(names []string, values []string) (area []Area, err error) {
	groups := make(map[string]string)
	for i, name := range names {
		if i != 0 {
			groups[name] = values[i]
		}
	}

	var radius float32
	ri, e := strconv.Atoi(groups["r"])
	if e != nil {
		err = e
		return
	}
	radius = float32(ri) * 1852.0
	if groups["t"] == "C" {
		radius /= 10.0
	}

	re := regexp.MustCompile(`(\d+)-(\d+)(?:[.,](\d+))?\s?([NSEW])`)
	mx := re.FindAllStringSubmatch(groups["x"], -1)
	my := re.FindAllStringSubmatch(groups["y"], -1)
	if mx == nil || my == nil || len(mx[0]) < 2 || len(my[0]) < 2 {
		err = errors.New("unable to match circle coords pattern")
		return
	}

	coords := Coords{
		Lat: coordsToFloat(mx[0][1], mx[0][2], mx[0][3], mx[0][4]),
		Lon: coordsToFloat(my[0][1], my[0][2], my[0][3], my[0][4]),
	}

	area = []Area{{
		Type:   "circle",
		Radius: radius,
		Coords: []Coords{coords},
		Raw:    fmt.Sprintf("%s%s %s %s", groups["r"], groups["t"], groups["x"], groups["y"]),
	}}
	return
}

func parsePolygon(cstr string) (area []Area, e error) {
	re, err := regexp.Compile(`((\d+-\d+(?:[.,-]\d+)?\s?[NS])\s+0?(\d+-\d+(?:[.,-]\d+)?\s?[EW])+)\r?\n?`)
	if err != nil {
		e = err
		return
	}

	m := re.FindAllStringSubmatch(cstr, -1)
	if m == nil {
		e = errors.New("unable to match polygon coords pattern")
		return
	}
	if len(m) < 2 {
		e = errors.New("invalid results of polygon coords matching")
		return
	}

	var coords []Coords
	raw := "["
	for _, cp := range m {
		if len(cp) != 4 {
			e = errors.New("unable to match polygon coords pattern")
			return
		}
		lat, err := stringToCoord(cp[2])
		if err != nil {
			e = err
			return
		}
		lon, err := stringToCoord(cp[3])
		if err != nil {
			e = err
			return
		}
		coords = append(coords, Coords{lat, lon})
		raw += fmt.Sprintf("%s %s, ", cp[2], cp[3])
	}
	raw = strings.TrimSuffix(raw, ", ") + "]"

	area = []Area{{
		Type:   "polygon",
		Radius: 0,
		Coords: coords,
		Raw:    raw,
	}}

	return
}

func parsePolygons(cstr string) (area []Area, e error) {
	re, err := regexp.Compile(`[A-Z].((?:\s+.*\r?\n)+)`)
	if err != nil {
		e = err
		return
	}

	m := re.FindAllStringSubmatch(cstr, -1)
	if m == nil {
		e = errors.New("unable to match polygons coords pattern")
		return
	}
	if len(m) != 2 {
		e = errors.New("invalid results of polygons coords matching")
		return
	}

	for _, cp := range m {
		if len(cp) != 2 {
			e = errors.New("unable to match polygons coords pattern")
			return
		}
		a, err := parsePolygon(cp[1])
		if err != nil {
			e = err
			return
		}
		area = append(area, a[0])
	}
	return
}

func ParseCoords(msg string, kind string, re *regexp.Regexp) (area []Area, err error) {
	m := re.FindAllStringSubmatch(msg, -1)
	if m == nil {
		err = errors.New("unable to match coords pattern")
		return
	}
	if len(m[0]) < 2 {
		err = errors.New("invalid results of coords matching")
		return
	}

	switch kind {
	case "circle":
		area, err = parseCircle(re.SubexpNames(), m[0])
	case "polygon":
		area, err = parsePolygon(m[0][1])
	case "polygons":
		area, err = parsePolygons(m[0][1])
	default:
		err = fmt.Errorf("unknown coordinates type %q", kind)
	}
	if err != nil {
		return
	}

	return
}
