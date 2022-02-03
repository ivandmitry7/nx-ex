package parser

import (
	"fmt"
	"regexp"
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
