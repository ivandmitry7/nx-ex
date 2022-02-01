package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type Task struct {
	cfg        *Config
	processors []Processor
}

func NewTask() *Task {
	return &Task{}
}

func (t *Task) Execute(opts map[string]interface{}) error {
	cfgPath := opts["--config"].(string)
	cfg, err := NewConfig(cfgPath)
	if err != nil {
		fmt.Printf("unable to load config file %q: %v\n", cfgPath, err)
		return nil
	}
	t.cfg = cfg

	for _, pc := range t.cfg.Parsers {
		p, err := NewProcessor(&pc)
		if err != nil {
			fmt.Println(err)
		} else {
			t.processors = append(t.processors, p)
		}
	}

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
		_ = t.processFile(fn)
	}
	return nil
}

func (t *Task) processFile(filename string) error {
	fmt.Printf("Process file %q... ", filename)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("unable to read file %q: %v\n", filename, err)
		return err
	}
	msg := string(b)

	var r *Result
	for i, pc := range t.processors {
		if pc.Check(msg) {
			fmt.Printf("Pattern %q dectected, parsing...", t.cfg.Parsers[i].Name)
			r, err = pc.Parse(msg)
			if err == nil {
				r.Commit()
				file, _ := json.MarshalIndent(r, "", "\t")
				jfn := filename + ".json"
				if err := ioutil.WriteFile(jfn, file, 0644); err != nil {
					fmt.Printf("unable to write result to file %q\n", jfn)
					return err
				}
				fmt.Println("Ok, json file created")
				return nil
			} else {
				fmt.Printf("Parsing error: %v\n", err)
				return err
			}
		}
	}

	fmt.Println("unable to detect pattern")
	return errors.New("unable to detect pattern")

	//fn := filepath.Base(filename)
	//switch {
	//case fn == "37.txt":
	//	r = &Result{
	//		Hash:   "77b790f1a15b6444a6f9c9784aa5d6c9",
	//		Source: "ODESA-NAVTEX",
	//		Reason: "NAVAL-TRAINING",
	//		Date:   Date{Beg: 1638316800, End: 1646092800},
	//		Area: []Area{
	//			{
	//				Type: "polygon",
	//				Coords: []Coords{
	//					{Lat: 46.285, Lon: 31.450277777777778},
	//					{Lat: 46.3172222, Lon: 31.450277777777778},
	//					{Lat: 46.3172222, Lon: 31.533611111111114},
	//					{Lat: 46.285, Lon: 31.535},
	//				},
	//				Radius: 0,
	//			},
	//		},
	//	}
	//case fn == "96.txt":
	//	r = &Result{
	//		Hash:   "9fe985729949ade919a0f0c5027f92c6",
	//		Source: "ODESA-NAVTEX",
	//		Reason: "NAVAL TRAINING",
	//		Date:   Date{Beg: 1640577600, End: 1642874400},
	//		Area: []Area{
	//			{
	//				Type: "polygon",
	//				Coords: []Coords{
	//					{Lat: 44.56805555555556, Lon: 37.06666666666667},
	//					{Lat: 44.5180556, Lon: 37.833333333333336},
	//					{Lat: 44.3347222, Lon: 37.8},
	//					{Lat: 44.4, Lon: 37.65138888888889},
	//				},
	//				Radius: 0,
	//			},
	//			{
	//				Type: "polygon",
	//				Coords: []Coords{
	//					{Lat: 44.8680556, Lon: 37.18555555555555},
	//					{Lat: 44.8844444, Lon: 37.25111111111111},
	//					{Lat: 44.7344444, Lon: 37.31777777777778},
	//					{Lat: 44.7180556, Lon: 37.25222222222222},
	//				},
	//				Radius: 0,
	//			},
	//		},
	//	}
	//case fn == "156.txt":
	//	r = &Result{
	//		Hash:   "f42ae0ba6fb05fd25ea883e7987c00cb",
	//		Source: "KALININGRAD",
	//		Reason: "SHIPS EXERCISES",
	//		Date:   Date{Beg: 1641790800, End: 1642262400},
	//		Area: []Area{
	//			{
	//				Type: "circle",
	//				Coords: []Coords{
	//					{Lat: 55.1, Lon: 19.833333333333332},
	//				},
	//				Radius: 298.172,
	//			},
	//		},
	//	}
	//default:
	//	return nil
	//}
}
