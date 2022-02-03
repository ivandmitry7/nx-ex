package parser

import (
	"errors"
	"fmt"
	"regexp"
)

func ParseText(caption string, msg string, re *regexp.Regexp) (string, error) {
	m := re.FindAllStringSubmatch(msg, -1)
	if m == nil {
		return "", fmt.Errorf("unable to match %s pattern", caption)
	}
	if len(m[0]) != 2 {
		return "", errors.New("invalid count of caption pattern matching")
	}

	return m[0][1], nil
}
