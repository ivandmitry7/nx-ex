package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/o-kos/nx-ex/internal/task"
	"os"
)

func main() {
	usage := `NAVTEX alert exercise messages extractor.
Usage:
  nxex [-vqr] [--config=NAME] [--out-dir=DIR] [--out-fmt=FMT] <MASK>
  nxex --version
  nxex -h | --help
Arguments:
  MASK source files name mask
Options:
  -h --help              show this help message and exit
  --version              show version and exit
  -v --verbose           print status messages
  -q --quiet             report only file names
  -r --recursive         recursive process nested directories
  -c NAME --config=NAME  config file name with parsing rules [default: ./nxex.yml]
  --out-dir=DIR          dir for save JSON results [default: ./]
  --out-fmt=FMT          template for out JSON file name [default: *+.json]`

	arguments, _ := docopt.ParseArgs(usage, nil, "0.0.1")
	tsk := task.NewTask()
	if err := tsk.Execute(arguments); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
