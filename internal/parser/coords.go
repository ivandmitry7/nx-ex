package parser

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
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

func parseCircle(template string, values []string) (area []Area, err error) {
	cs, err := TemplateToString(template, values)
	if err != nil {
		return
	}
	re, err := regexp.Compile(`(\d+)([CBK])r(\d+)-(\d+).(\d+)([NS]) (\d+)-(\d+).(\d+)([EW])`)
	if err != nil {
		return
	}

	m := re.FindAllStringSubmatch(cs, -1)
	if m == nil {
		err = errors.New("unable to match circle coords pattern")
		return
	}
	if len(m[0]) < 2 {
		err = errors.New("invalid results of circle coords matching")
		return
	}

	var radius float32
	ri, _ := strconv.Atoi(m[0][1])
	radius = float32(ri) * 1852.0
	if m[0][2] == "C" {
		radius /= 10.0
	}

	coords := Coords{
		Lat: coordsToFloat(m[0][3], m[0][4], m[0][5], m[0][6]),
		Lon: coordsToFloat(m[0][7], m[0][8], m[0][9], m[0][10]),
	}

	area = []Area{{
		Type:   "circle",
		Radius: radius,
		Coords: []Coords{coords},
		Raw:    cs,
	}}

	return
}

func ParseCoords(msg string, kind string, re *regexp.Regexp, template string) (area []Area, err error) {
	_ = ioutil.WriteFile("~msg.txt", []byte(msg), 0644)
	_ = ioutil.WriteFile("~re.txt", []byte(re.String()), 0644)

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
		area, err = parseCircle(template, m[0])
	default:
		err = fmt.Errorf("unknown coordinates type %q", kind)
	}
	if err != nil {
		return
	}

	return
}
