package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

func ParseText(caption string, msg string, re *regexp.Regexp) (string, error) {
	m := re.FindAllStringSubmatch(msg, -1)
	if m == nil {
		return "", fmt.Errorf("unable to match %s pattern", caption)
	}
	if len(m[0]) != 2 {
		return "", fmt.Errorf("invalid count of %s pattern matching", caption)
	}

	return m[0][1], nil
}

func TemplateToStringFunc(template string, values []string, repl func(string) string) (string, error) {
	r := regexp.MustCompile(`\${(\d+)}`)
	var retErr error = nil
	dtErr := errors.New("unable to apply replace template pattern")
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

			return repl(values[idx])
		},
	)
	return str, retErr
}

func TemplateToString(template string, values []string) (string, error) {
	return TemplateToStringFunc(
		template, values, func(s string) string {
			return s
		},
	)
}
