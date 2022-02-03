package parser

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProcessor_Parse(t *testing.T) {
	input := "20${5}-${4}-${1}T${2}:${3}:00Z\n20${10}-${9}-${5}T${7}:${8}:00Z"
	r := regexp.MustCompile(`\${(\d+)}`)
	str := r.ReplaceAllStringFunc(
		input, func(m string) string {
			matches = []string{"", "01", "02", "03", "04", "05", "06", "07", "08", "09", "10"}
			fmt.Printf("m: %q, %v\n", m, parts)
			parts := r.FindStringSubmatch(m)
			idx := part
			return "<" + parts[1] + ">"
		},
	)

	assert.Equal(t, "2022", str)
}
