package run

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/o-kos/dmd-cli.go/pkg/dmdintf"
)

func TestResults_Unmarshal(t *testing.T) {
	err := dmdintf.InitTestSystem()
	require.NoError(t, err)
	defer dmdintf.SystemDone()

	dmd, err := dmdintf.CreateTestDmd()
	require.NoError(t, err)

	err = dmd.Start(8000)
	require.NoError(t, err)

	var samples16s = make([]int16, 800, 800)
	for i := range samples16s {
		samples16s[i] = int16(i)
	}

	res, err := dmd.Process16s(samples16s)
	require.NoError(t, err)
	require.Equal(t, false, res)
	res, err = dmd.Process16s(samples16s)
	require.NoError(t, err)
	require.Equal(t, false, res)
	res, err = dmd.Process16s(samples16s)
	require.NoError(t, err)
	require.Equal(t, true, res)

	rs, err := dmd.ReadResults()
	require.NoError(t, err)
	var ra ResultsArray
	err = json.Unmarshal([]byte(rs), &ra)
	require.NoError(t, err)
	require.Equal(t, 2, len(ra))
	assert.Equal(t, int64(200), ra[0].Ticks)
	assert.Equal(t, "Log:      3 16s:0...799 (800.1)", ra[0].Log)
	assert.Equal(t, "Txt:      3 16s:0...799 (800.1)\nTxt:      3 16s:0...799 (800.2)\n", ra[0].Text)
	assert.Equal(t, int64(200), ra[1].Ticks)
	assert.Equal(t, "Log:      3 16s:0...799 (800.3)", ra[1].Log)
	assert.Equal(t, "Txt:      3 16s:0...799 (800.3)\nTxt:      3 16s:0...799 (800.4)\n", ra[1].Text)
}
