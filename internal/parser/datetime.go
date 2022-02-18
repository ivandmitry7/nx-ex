package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func monthsToNumber(month string) string {
	months := map[string]string{
		"JAN": "01", "FEB": "02", "MAR": "03", "APR": "04", "MAY": "05", "JUN": "06",
		"JUL": "07", "AUG": "08", "SEP": "09", "OCT": "10", "NOV": "11", "DEC": "12",
	}
	if val, ok := months[month]; ok {
		return val
	}

	return month
}

func matchToTimes(groups map[string]string) (raw string, t0 time.Time, t1 time.Time, err error) {
	var y, m, d, h, n, s int

	k2v := func(key string) string {
		idx := key + "0"
		if groups[idx] == "" {
			idx = key + "1"
		}
		return groups[idx]
	}

	y, err = strconv.Atoi(k2v("y"))
	if err != nil {
		return
	}
	m, err = strconv.Atoi(monthsToNumber(k2v("m")))
	if err != nil {
		return
	}
	d, err = strconv.Atoi(k2v("d"))
	if err != nil {
		return
	}
	h, err = strconv.Atoi(k2v("h"))
	if err != nil {
		return
	}
	n, err = strconv.Atoi(k2v("n"))
	if err != nil {
		return
	}
	s, _ = strconv.Atoi(k2v("s"))

	str0 := fmt.Sprintf("20%.2d-%.2d-%.2dT%.2d:%.2d:%.2dZ", y, m, d, h, n, s)
	t0, err = time.Parse(time.RFC3339, str0)
	if err != nil {
		return
	}

	if groups["y1"] != "" {
		y, _ = strconv.Atoi(groups["y1"])
	}
	if groups["m1"] != "" {
		m, _ = strconv.Atoi(monthsToNumber(groups["m1"]))
	}
	if groups["d1"] != "" {
		d, _ = strconv.Atoi(groups["d1"])
	}
	if groups["h1"] != "" {
		h, _ = strconv.Atoi(groups["h1"])
	}
	if groups["n1"] != "" {
		n, _ = strconv.Atoi(groups["n1"])
	}
	if groups["s1"] != "" {
		s, _ = strconv.Atoi(groups["s1"])
	}
	str1 := fmt.Sprintf("20%.2d-%.2d-%.2dT%.2d:%.2d:%.2dZ", y, m, d, h, n, s)
	t1, err = time.Parse(time.RFC3339, str1)
	raw = str0 + ";" + str1
	return
}

func ParseDateTime(msg string, re *regexp.Regexp) (times Date, err error) {
	m := re.FindAllStringSubmatch(msg, -1)
	if m == nil {
		err = errors.New("unable to match pattern")
		return
	}

	err = errors.New("invalid results of datetime matching")
	if len(m[0]) < 2 {
		return
	}

	var t0 time.Time
	var t1 time.Time
	var timeStr string
	var e error
	groups := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 {
			groups[name] = m[0][i]
		}
	}
	timeStr, t0, t1, e = matchToTimes(groups)
	if e != nil {
		err = e
		return
	}

	err = nil
	times.Beg = t0.Unix()
	times.End = t1.Unix()
	times.Raw = timeStr
	return
}
