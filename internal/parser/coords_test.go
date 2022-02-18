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

	search := `(?m)^WITHIN (?P<r>\d+) (?P<t>C)ABLES OF (?P<x>\d\d-\d\d.\d\s?[NS])\s+(?P<y>0\d\d-\d\d.\d\d\s?[EW])$`

	re, err := regexp.Compile(search)
	require.NoError(t, err)

	area, err := ParseCoords(msg, "circle", re)
	require.NoError(t, err)
	require.Equal(t, 1, len(area))
	assert.Equal(t, "circle", area[0].Type)
	assert.Equal(t, float32(1852), area[0].Radius)
	assert.Equal(t, 1, len(area[0].Coords))
	assert.Equal(t, float32(45.01889), area[0].Coords[0].Lat)
	assert.Equal(t, float32(33.51472), area[0].Coords[0].Lon)
}

func TestParseCoords_Polygon(t *testing.T) {
	msg :=
		`ZCZC CL46
2607303UTC NOV 21
UKRAINE COASTAL WARNING 584/21 ODESA-NAVTEX
BLACK SEA
TENDRIVSKA KOSA ISLAND

1. NAVAL TRAINING 
010000 UTC DEC 21 THRU 010000 UTC MAR 22
NAVIGATION DANGEROUS IN AREA BOUNDED BY
46-17.6N  031-27.1E
46-19.2N  031-27.1E
46-19.2N  031-32.1E
46-17.6N  031-32.6E
2. CANCEL THIS MSG 010100 UTC MAR 22
NNN`

	search := `AREA BOUNDED BY\r?\n((?:\d\d-\d\d.\d+\s?[NS]\s+0\d\d-\d\d.\d+\s?[EW]+\r?\n)+)2. CANCEL`
	re, err := regexp.Compile(search)
	require.NoError(t, err)

	area, err := ParseCoords(msg, "polygon", re)
	require.NoError(t, err)
	require.Equal(t, 1, len(area))
	assert.Equal(t, "polygon", area[0].Type)
	assert.Equal(t, float32(0), area[0].Radius)
	assert.Equal(t, 4, len(area[0].Coords))
	assert.Equal(t, float32(46.285), area[0].Coords[0].Lat)
	assert.Equal(t, float32(31.45028), area[0].Coords[0].Lon)
}

func TestParseCoords_Polygons(t *testing.T) {
	msg :=
		`ZCZC CL88
281050 UTC DEC 21
UKRAINE COASTAL WARNING 636/21 ODESA-NAVTEX
BLACK SEA

1. NAVAL TRAINING
271800 DEC 21 THRU 271800 UTC JAN 22
NAVIGATION DANGEROUS IN AREAS BOUNDED BY
A. 45-44.5N 031-33.0E
   45-45.3N 032-01.0E
   45-36.3N 032-01.0E
   45-36.0N 031-33.0E
B. 45-24.0N 031-34.0E
   45-24.0N 032-00.0E
   45-11.0N 032-00.0E
   45-11.0N 031-34.0E
2. CANCEL THIS MSG 271900 UTC JAN 22
NNN`

	search := `AREAS BOUNDED BY\r?\n(((?:[A-Z])?.\s+\d\d-\d\d.\d+\s?[NS]\s+0?\d\d-\d\d.\d+\s?[EW]\r?\n)+)2. CANCEL`
	re, err := regexp.Compile(search)
	require.NoError(t, err)

	area, err := ParseCoords(msg, "polygons", re)
	require.NoError(t, err)
	require.Equal(t, 2, len(area))
	assert.Equal(t, "polygon", area[0].Type)
	assert.Equal(t, float32(0), area[0].Radius)
	assert.Equal(t, 4, len(area[0].Coords))
	assert.Equal(t, float32(45.734722), area[0].Coords[0].Lat)
	assert.Equal(t, float32(31.55), area[0].Coords[0].Lon)
}
