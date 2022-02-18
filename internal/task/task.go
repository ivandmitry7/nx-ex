package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/o-kos/nx-ex/internal/parser"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Task struct {
	cfg *Config
	opt struct {
		verbose bool
		quiet   bool
		outDir  string
	}
	processors []Processor
	stats      map[string]int
}

func NewTask() *Task {
	return &Task{stats: map[string]int{}}
}

var ErrUndetectedPattern error = errors.New("unable to detect pattern")

func printStats(stats map[string]int) {
	type kv struct {
		key   string
		value int
	}

	var ss []kv
	for k, v := range stats {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(
		ss, func(i, j int) bool {
			return ss[i].value > ss[j].value
		},
	)

	for _, kv := range ss {
		fmt.Printf("%-20s %8d\n", kv.key, kv.value)
	}
}

func (t *Task) Execute(opts map[string]interface{}) error {
	cfgPath := opts["--config"].(string)
	cfg, err := NewConfig(cfgPath)
	if err != nil {
		fmt.Printf("unable to load config file %q: %v\n", cfgPath, err)
		return nil
	}
	t.cfg = cfg
	t.opt.outDir = opts["--out-dir"].(string)
	if t.opt.outDir == `""` {
		t.opt.outDir = ""
	}
	t.opt.verbose = opts["--verbose"].(bool)
	t.opt.quiet = opts["--quiet"].(bool)
	if t.opt.quiet {
		t.opt.verbose = false
	}

	for _, pc := range t.cfg.Parsers {
		p, err := NewProcessor(&pc)
		if err != nil {
			fmt.Println(err)
		} else {
			t.processors = append(t.processors, p)
		}
	}

	if t.opt.verbose {
		fmt.Printf("Loaded %d message processors\n", len(t.processors))
	}
	mask := opts["<MASK>"].(string)
	files, err := filepath.Glob(mask)
	if err != nil {
		return err
	}
	totalFiles := len(files)
	if totalFiles == 0 {
		return fmt.Errorf("unable to find files with mask %q\n", mask)
	}
	if totalFiles > 1 && !t.opt.quiet {
		fmt.Printf("Found %d files\n", len(files))
	}

	start := time.Now()
	okFiles := 0
	badFiles := 0
	for i, fn := range files {
		if totalFiles > 1 && !t.opt.quiet {
			fmt.Printf("%d/%d [%d%%] ", i+1, totalFiles, (i+1)*100/totalFiles)
		}
		err := t.processFile(fn)
		if err == nil {
			okFiles++
		} else {
			if !errors.Is(err, ErrUndetectedPattern) {
				badFiles++
			}
		}
	}
	if totalFiles > 1 {
		badMsg := ""
		if badFiles > 0 {
			badMsg = fmt.Sprintf("%d files has errors, ", badFiles)
		}
		fmt.Printf("Total patterns detected in %d files, %selapsed time %s\n", okFiles, badMsg, time.Since(start))
		if t.opt.verbose {
			printStats(t.stats)
		}
	}
	return nil
}

func normalizeCRLF(str string) string {
	for {
		ns := strings.ReplaceAll(str, "\r\n", "\n")
		if len(ns) == len(str) {
			break
		}
		str = ns
	}
	return strings.ReplaceAll(str, "\r", "\n")
}

func makeOutName(filename string, outDir string) string {
	if outDir != "" {
		filename = filepath.Join(outDir, filepath.Base(filename))
	}
	return filename + ".json"
}

func (t *Task) processFile(filename string) error {
	if !t.opt.quiet {
		fmt.Printf("Process file %q... ", filename)
	}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("unable to read file %q: %v\n", filename, err)
		return err
	}
	msg := normalizeCRLF(string(b))

	var r *parser.Result
	for i, pc := range t.processors {
		if pc.Check(msg) {
			if !t.opt.quiet {
				fmt.Printf("pattern %q dectected, ", t.cfg.Parsers[i].Name)
			}
			r, err = pc.Parse(msg)
			if err == nil {
				r.Commit()
				file, _ := json.MarshalIndent(r, "", "\t")
				jfn := makeOutName(filename, t.opt.outDir)
				if err := ioutil.WriteFile(jfn, file, 0644); err != nil {
					fmt.Printf("unable to write result to file %q\n", jfn)
					return err
				}
				if !t.opt.quiet {
					fmt.Printf("%q message extracted\n", r.Source)
				}

				count := t.stats[r.Source]
				t.stats[r.Source] = count + 1
				return nil
			} else {
				fmt.Printf("parsing error: %v\n", err)
				return err
			}
		}
	}

	fmt.Println(ErrUndetectedPattern.Error())
	return ErrUndetectedPattern
}
