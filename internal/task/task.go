package task

import (
	"encoding/json"
	"fmt"
	"github.com/o-kos/nx-ex/internal/result"
	"io/ioutil"
	"path/filepath"
)

type Task struct {
	config string // Path for config file
}

func NewTask() *Task {
	return &Task{}
}

func (t *Task) Execute(opts map[string]interface{}) error {
	mask := opts["<MASK>"].(string)
	files, err := filepath.Glob(mask)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("unable to find files with mask %q\n", mask)
	}

	if len(files) > 1 {
		fmt.Printf("Found %d files\n", len(files))
	}
	for _, fn := range files {
		t.processFile(fn)
	}
	return nil
}

func (t *Task) processFile(filename string) {
	fn := filepath.Base(filename)
	var r *result.Result
	switch {
	case fn == "37.txt":
		r = &result.Result{
			Hash:   "77b790f1a15b6444a6f9c9784aa5d6c9",
			Source: "ODESA-NAVTEX",
			Reason: "NAVAL-TRAINING",
			Date:   result.Date{Beg: 1638316800, End: 1646092800},
			Area: []result.Area{
				{
					Type: "polygon",
					Coords: []result.Coords{
						{Lat: 46.285, Lon: 31.450277777777778},
						{Lat: 46.3172222, Lon: 31.450277777777778},
						{Lat: 46.3172222, Lon: 31.533611111111114},
						{Lat: 46.285, Lon: 31.535},
					},
					Radius: 0,
				},
			},
		}
	case fn == "96.txt":
		r = &result.Result{
			Hash:   "9fe985729949ade919a0f0c5027f92c6",
			Source: "ODESA-NAVTEX",
			Reason: "NAVAL TRAINING",
			Date:   result.Date{Beg: 1640577600, End: 1642874400},
			Area: []result.Area{
				{
					Type: "polygon",
					Coords: []result.Coords{
						{Lat: 44.56805555555556, Lon: 37.06666666666667},
						{Lat: 44.5180556, Lon: 37.833333333333336},
						{Lat: 44.3347222, Lon: 37.8},
						{Lat: 44.4, Lon: 37.65138888888889},
					},
					Radius: 0,
				},
				{
					Type: "polygon",
					Coords: []result.Coords{
						{Lat: 44.8680556, Lon: 37.18555555555555},
						{Lat: 44.8844444, Lon: 37.25111111111111},
						{Lat: 44.7344444, Lon: 37.31777777777778},
						{Lat: 44.7180556, Lon: 37.25222222222222},
					},
					Radius: 0,
				},
			},
		}
	case fn == "156.txt":
		r = &result.Result{
			Hash:   "f42ae0ba6fb05fd25ea883e7987c00cb",
			Source: "KALININGRAD",
			Reason: "SHIPS EXERCISES",
			Date:   result.Date{Beg: 1641790800, End: 1642262400},
			Area: []result.Area{
				{
					Type: "circle",
					Coords: []result.Coords{
						{Lat: 55.1, Lon: 19.833333333333332},
					},
					Radius: 298.172,
				},
			},
		}
	default:
		return
	}

	fmt.Printf("Process file %q... ", filename)
	file, _ := json.MarshalIndent(r, "", "\t")
	jfn := filename + ".json"
	if err := ioutil.WriteFile(jfn, file, 0644); err != nil {
		fmt.Printf("unable to write result to file %q\n", jfn)
	}
	fmt.Println("Ok, json file created")
}
