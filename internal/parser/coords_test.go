package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestParseCoords_Circle(t *testing.T) {
	msg :=
		`ZCZC CL96
092230 UTC JAN 22
UKRAINE COASTAL WARNING 03/22 ODESA-NAVTEX
BLACK SEA

1. NAVAL TRAINING
FROM 092100 UTC THRU 312100 UTC JAN 22
NAVIGATION DANGEROUS
WITHIN 10 CABLES OF 45-01.8 N  033-30.53 E
2. CANCEL THIS MSG 312200 UTC JAN 22
NNN`

	search := `(?m)^WITHIN (\d\d) CABLES OF (\d\d-\d\d.\d)\s?([NS])\s+(0\d\d-\d\d.\d\d)\s?([EW])\r?$`
	replace := `${1}Cr${2}${3} ${4}${5}`

	re, err := regexp.Compile(search)
	require.NoError(t, err)

	area, err := ParseCoords(msg, "circle", re, replace)
	require.NoError(t, err)
	require.Equal(t, 1, len(area))
	assert.Equal(t, "circle", area[0].Type)
	assert.Equal(t, float32(1852), area[0].Radius)
	assert.Equal(t, 1, len(area[0].Coords))
	assert.Equal(t, float32(45.01889), area[0].Coords[0].Lat)
	assert.Equal(t, float32(33.51472), area[0].Coords[0].Lon)
}
