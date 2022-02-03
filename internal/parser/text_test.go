package parser

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseText(t *testing.T) {
	re, err := regexp.Compile(`(?m)\b(ODESA-NAVTEX)[\s\r]*$`)
	require.NoError(t, err)
	msg := "2607303UTC NOV 21\nUKRAINE COASTAL WARNING 584/21 ODESA-NAVTEX\r\nBLACK SEA"
	str, err := ParseText("source", msg, re)
	require.NoError(t, err)
	assert.Equal(t, "ODESA-NAVTEX", str)
}
