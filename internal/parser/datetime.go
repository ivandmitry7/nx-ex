package parser

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func templateToString(template string, values []string) (string, error) {
	r := regexp.MustCompile(`\${(\d+)}`)
	var retErr error = nil
	dtErr := errors.New("unable to apply datetime replace pattern")
	str := r.ReplaceAllStringFunc(
		template, func(m string) string {
			if retErr != nil {
				return ""
			}
			parts := r.FindStringSubmatch(m)
			if len(parts) < 2 {
				retErr = dtErr
				return ""
			}
			idx, err := strconv.Atoi(parts[1])
			if err != nil || idx >= len(values) {
				retErr = dtErr
				return ""
			}

			months := map[string]string{
				"JAN": "01", "FEB": "02", "MAR": "03", "APR": "04", "MAY": "05", "JUN": "06",
				"JUL": "07", "AUG": "08", "SEP": "09", "OCT": "10", "NOV": "11", "DEC": "12",
			}
			if val, ok := months[values[idx]]; ok {
				return val
			}

			return values[idx]
		},
	)
	return str, retErr
}

func ParseDateTime(msg string, re *regexp.Regexp, template string) (raw string, times [2]time.Time, err error) {
	m := re.FindAllStringSubmatch(msg, -1)
	if m == nil {
		err = errors.New("unable to match coords pattern")
		return
	}
	if len(m[0]) < 2 {
		err = errors.New("invalid results of coords matching")
		return
	}

	timeStr, err := templateToString(template, m[0])
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

	raw = timeStr
	times = [2]time.Time{t0, t1}
	return
}