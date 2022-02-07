package parser

import (
	"errors"
	"regexp"
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

func ParseDateTime(msg string, re *regexp.Regexp, template string) (times Date, err error) {
	m := re.FindAllStringSubmatch(msg, -1)
	if m == nil {
		err = errors.New("unable to match coords pattern")
		return
	}
	if len(m[0]) < 2 {
		err = errors.New("invalid results of coords matching")
		return
	}

	timeStr, err := TemplateToStringFunc(template, m[0], monthsToNumber)
	if err != nil {
		return
	}

	ts := strings.Split(timeStr, ";")
	if len(ts) != 2 {
		err = errors.New("invalid results of datetime matching")
		return
	}
	t0, err := time.Parse(time.RFC3339, ts[0])
	if err != nil {
		err = errors.New("invalid results of datetime matching")
		return
	}
	t1, err := time.Parse(time.RFC3339, ts[1])
	if err != nil {
		err = errors.New("invalid results of datetime matching")
		return
	}

	times.Beg = t0.Unix()
	times.End = t1.Unix()
	times.Raw = timeStr
	return
}
