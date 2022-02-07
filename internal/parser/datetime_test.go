package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

func TestParseDateTime(t *testing.T) {
	msg := "BLACK SEA\n1. NAVAL TRAINING \n011234 UTC DEC 21 THRU 100713 UTC MAR 22\nNAVIGATION DANGEROUS IN AREA BOUNDED BY"
	search := `(?s)1\. [[A-Z\s\s\r]*\n(\d\d)(\d\d)(\d\d) UTC ([A-Z]{3}) (\d+) THRU (\d\d)(\d\d)(\d\d) UTC ([A-Z]{3}) (\d+)`
	replace := `20${5}-${4}-${1}T${2}:${3}:00Z;20${10}-${9}-${6}T${7}:${8}:00Z`

	re, err := regexp.Compile(search)
	require.NoError(t, err)

	tms, err := ParseDateTime(msg, re, replace)
	require.NoError(t, err)
	assert.Equal(t, "2021-12-01T12:34:00Z;2022-03-10T07:13:00Z", tms.Raw)
	t0, _ := time.Parse(time.RFC3339, "2021-12-01T12:34:00Z")
	tb := time.Unix(tms.Beg, 0).UTC()
	assert.Equal(t, t0, tb)
	t1, _ := time.Parse(time.RFC3339, "2022-03-10T07:13:00Z")
	te := time.Unix(tms.End, 0).UTC()
	assert.Equal(t, t1, te)

	_, err = ParseDateTime(msg, re, "${111}")
	assert.EqualError(t, err, "unable to apply replace template pattern")
}
