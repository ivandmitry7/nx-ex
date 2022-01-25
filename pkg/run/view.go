package run

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"
)

type Viewer interface {
	PrintHeader(model Model)
	PrintBody(model Model)
	PrintFooter(model Model)
}

type ConsoleView struct {
	lastPercent int
}

func NewConsoleView() Viewer {
	return &ConsoleView{}
}

func (c *ConsoleView) PrintHeader(model Model) {
	fmt.Printf("Process file %q, output dir %q\n", filepath.ToSlash(model.Path), filepath.ToSlash(model.Opt.OutDir))
	fmt.Printf("%v, batch %d ms. Freq = %d Hz", model.SrcInfo, model.Opt.Batch, model.Opt.Freq)
	if len(model.Opt.Params) > 0 {
		s := ""
		for k, v := range model.Opt.Params {
			s = s + fmt.Sprintf("%s=%q, ", k, v)
		}
		s = strings.TrimSuffix(s, ", ")
	}
	fmt.Println()
}

func stat2str(rs *ResultStat) string {
	s := ", "
	if rs.Text > 0 {
		s += fmt.Sprintf("t: %d, ", rs.Text)
	}
	if rs.Bits > 0 {
		s += fmt.Sprintf("b: %d, ", rs.Bits/8)
	}
	if rs.Ale > 0 {
		s += fmt.Sprintf("a: %d, ", rs.Ale)
	}
	if rs.Phase > 0 {
		s += fmt.Sprintf("p: %d, ", rs.Phase)
	}
	return strings.TrimSuffix(s, ", ")
}

func (c *ConsoleView) PrintBody(model Model) {
	r := model.Results.Back().Value.(Result)
	if len(r.Log) == 0 && !r.Stat.HasData && model.Progress.Percent/10 == c.lastPercent {
		return
	}
	c.lastPercent = model.Progress.Percent / 10
	fmt.Printf("[%3d%%] ", int(math.Round(float64(model.Progress.Percent))))
	if len(r.Log) > 0 {
		fmt.Printf("%s", r.Log)
	}
	fmt.Println(stat2str(&r.Stat))
}

func (c *ConsoleView) PrintFooter(model Model) {
	fmt.Printf("Total time %v", model.Progress.Stop.Sub(model.Progress.Start))
	rs := ", "
	if model.Total.Text > 0 {
		rs += fmt.Sprintf("text: %v bytes, ", model.Total.Text)
	}
	if model.Total.Bits > 0 {
		rs += fmt.Sprintf("bits: %v bytes, ", model.Total.Bits)
	}
	if model.Total.Ale > 0 {
		rs += fmt.Sprintf("ALE: %v bytes, ", model.Total.Ale)
	}
	if model.Total.Phase > 0 {
		rs += fmt.Sprintf("phase: %v pts, ", model.Total.Phase)
	}
	rs = strings.TrimRight(rs, ", ")
	if len(rs) > 0 {
		fmt.Print(rs)
	}
	fmt.Println()
}
