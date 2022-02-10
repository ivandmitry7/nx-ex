package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

func stringToTime(s string) (t time.Time, err error) {
	re := regexp.MustCompile(`(\d+)-(\d+)-(\d+)T(\d+):(\d+):(\d+)Z`)
	m := re.FindAllStringSubmatch(s, -1)
	if m == nil {
		err = fmt.Errorf("unable to convert datetime string %q", s)
		return
	}
	if len(m[0]) != 7 {
		err = fmt.Errorf("invalid datetime format %q", s)
		return
	}

	yi, _ := strconv.Atoi(m[0][1])
	mi, _ := strconv.Atoi(m[0][2])
	di, _ := strconv.Atoi(m[0][3])
	hi, _ := strconv.Atoi(m[0][4])
	ni, _ := strconv.Atoi(m[0][5])
	si, _ := strconv.Atoi(m[0][6])

	s = fmt.Sprintf("%.4d-%.2d-%.2dT%.2d:%.2d:%.2dZ", yi, mi, di, hi, ni, si)
	t, err = time.Parse(time.RFC3339, s)
	return
}

func ParseDateTime(msg string, re *regexp.Regexp, template string) (times Date, err error) {
	m := re.FindAllStringSubmatch(msg, -1)
	if m == nil {
		err = errors.New("unable to match pattern")
		return
	}

	err = errors.New("invalid results of datetime matching")
	if len(m[0]) < 2 {
		return
	}

	timeStr, e := TemplateToStringFunc(template, m[0], monthsToNumber)
	if e != nil {
		err = e
		return
	}

	ts := strings.Split(timeStr, ";")
	if len(ts) != 2 {
		return
	}
	t0, e := stringToTime(ts[0])
	if e != nil {
		err = e
		return
	}
	t1, e := stringToTime(ts[1])
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
