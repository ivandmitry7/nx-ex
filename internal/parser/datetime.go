package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

func templateToString(template string, values []string) (string, error) {
	r := regexp.MustCompile(`\${(\d+)}`)
	str := r.ReplaceAllStringFunc(
		template, func(m string) string {
			//fmt.Printf("m: %q, %v\n", m, parts)
			parts := r.FindStringSubmatch(m)
			idx := parts[1]
			return values[idx]
		},
	)
	return str, nil
}

func ParseTime(msg string, re *regexp.Regexp, template string) ([]time.Time, error) {
	m := re.FindAllStringSubmatch(msg, -1)
	if m == nil {
		return nil, errors.New("unable to match datetime pattern")
	}
	if len(m[0]) < 2 {
		return nil, errors.New("invalid results of datetime matching")
	}

	timeStr, err := templateToString(template, m[0])
	if err != nil {
		return nil, err
	}
	replacer := strings.NewReplacer(
		"JAN", "01", "FEB", "02", "MAR", "03", "APR", "04", "MAY", "05", "JUN", "06",
		"JUL", "07", "AUG", "08", "SEP", "09", "OCT", "10", "NOV", "11", "DEC", "12",
	)
	timeStr = replacer.Replace(timeStr)
	ts := strings.Split(timeStr, "\n")
	if len(ts) != 2 {
		continue
	}
	t0, err := time.Parse(time.RFC3339, ts[0])
	if err != nil {
		continue
	}
	t1, err := time.Parse(time.RFC3339, ts[1])
	if err != nil {
		continue
	}
	return []time.Time{t0, t1}, nil
}
