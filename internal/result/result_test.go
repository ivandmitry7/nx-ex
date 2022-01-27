package result

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResult_PushCoords(t *testing.T) {
	r := Result{}
	cs := []string{"44-34.5N 037-42.0E", "44-31.5N 037-50.0E", "44-20.5N 037-48.0E", "44-24.0N 037-39.5E"}
	r.PushCoords("", cs)
	require.Equal(t, 1, len(r.Area))
	assert.Equal(t, "polygon", r.Area[0].Type)
	assert.Equal(t, float32(0), r.Area[0].Radius)
	assert.Equal(t, 4, len(r.Area[0].Coords))
	assert.Equal(
		t, "polygon/[[44-34.5N 037-42.0E], [44-31.5N 037-50.0E], [44-20.5N 037-48.0E], [44-24.0N 037-39.5E]]",
		r.Area[0].raw,
	)

	r = Result{}
	cs = []string{"44-34.5N 037-42.0E", "44-31.5N 037-50.0E"}
	r.PushCoords("1000", cs)
	require.Equal(t, 1, len(r.Area))
	assert.Equal(t, "circle", r.Area[0].Type)
	assert.Equal(t, 1852, int(math.Round(float64(r.Area[0].Radius))))
	assert.Equal(t, 2, len(r.Area[0].Coords))
	assert.Equal(t, "circle/1000[[44-34.5N 037-42.0E], [44-31.5N 037-50.0E]]", r.Area[0].raw)
}

func TestResult_makeID(t *testing.T) {
	r := Result{
		Source: "ODESA-NAVTEX",
		Reason: "NAVAL-TRAINING",
		Date:   Date{Beg: 1638316800, End: 1646092800},
		Area: []Area{
			{
				raw: "polygon/[[44-34.5N 037-42.0E], [44-31.5N 037-50.0E], [44-20.5N 037-48.0E]]",
			},
		},
	}
	assert.Equal(
		t,
		"ODESA-NAVTEX#NAVAL-TRAINING:2021-12-01T00:00:00Z...2022-03-01T00:00:00Z[[polygon/[[44-34.5N 037-42.0E], [44-31.5N 037-50.0E], [44-20.5N 037-48.0E]]]]",
		r.makeID(),
	)
}

func TestResult_Commit(t *testing.T) {
	r := Result{
		Source: "ODESA-NAVTEX",
		Reason: "NAVAL-TRAINING",
		Date:   Date{Beg: 1638316800, End: 1646092800},
		Area: []Area{
			{
				raw: "polygon/[[44-34.5N 037-42.0E], [44-31.5N 037-50.0E], [44-20.5N 037-48.0E]]",
			},
		},
	}
	r.Commit()
	assert.Equal(t, "50677145aef19ed76a6d2b4b7b304455", r.Hash)
}
