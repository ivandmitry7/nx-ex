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
	search := `1.\s.+\n(?:FROM )?(?P<d0>\d\d)(?P<h0>\d\d)(?P<n0>\d\d)(?:\sUTC\s?)?\s?(?P<m0>DEC)?\s(?P<y0>\d+)?\s?\n?THRU (?P<d1>\d\d)(?P<h1>\d\d)(?P<n1>\d\d)\s+UTC (?P<m1>MAR) (?P<y1>\d+)`

	re, err := regexp.Compile(search)
	require.NoError(t, err)

	tms, err := ParseDateTime(msg, re)
	require.NoError(t, err)
	assert.Equal(t, "2021-12-01T12:34:00Z;2022-03-10T07:13:00Z", tms.Raw)
	t0, _ := time.Parse(time.RFC3339, "2021-12-01T12:34:00Z")
	tb := time.Unix(tms.Beg, 0).UTC()
	assert.Equal(t, t0, tb)
	t1, _ := time.Parse(time.RFC3339, "2022-03-10T07:13:00Z")
	te := time.Unix(tms.End, 0).UTC()
	assert.Equal(t, t1, te)
}
