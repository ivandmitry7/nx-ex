package task

import (
	"fmt"
	"github.com/o-kos/nx-ex/internal/parser"
	"regexp"
	"strings"
)

type Processor interface {
	Check(msg string) bool
	Parse(msg string) (*parser.Result, error)
}

func NewProcessor(cfg *ParserCfg) (Processor, error) {
	var p ReProc
	if err := p.Compile(cfg); err != nil {
		return nil, fmt.Errorf("unable to compile parser rules for %q: %v", cfg.Name, err)
	}
	return p, nil
}

type reTimes struct {
	search  *regexp.Regexp
	replace string
}

type reCoords struct {
	kind    string
	search  *regexp.Regexp
	replace string
}

type ReProc struct {
	source *regexp.Regexp
	reason *regexp.Regexp
	times  []reTimes
	coords []reCoords
}

func (p *ReProc) Compile(cfg *ParserCfg) error {
	r := strings.NewReplacer(
		`\_`, `(?:\s|\r?\n)`,
		`\M`, `JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC`,
	)

	re, err := regexp.Compile(r.Replace(cfg.Source))
	if err != nil {
		return err
	}
	p.source = re

	re, err = regexp.Compile(r.Replace(cfg.Reason))
	if err != nil {
		return err
	}
	p.reason = re

	for _, ts := range cfg.Times {
		re, err := regexp.Compile(r.Replace(ts.Search))
		if err != nil {
			return err
		}
		p.times = append(p.times, reTimes{re, ts.Replace})
	}

	for _, cs := range cfg.Coords {
		re, err := regexp.Compile(r.Replace(cs.Search))
		if err != nil {
			return err
		}
		p.coords = append(p.coords, reCoords{cs.Type, re, cs.Replace})
	}
	return nil
}

func (p ReProc) Check(msg string) bool {
	str := p.source.FindString(msg)
	if str != "" {
		return true
	} else {
		return false
	}
}

func (p ReProc) Parse(msg string) (*parser.Result, error) {
	res := parser.Result{}
	s, err := parser.ParseText("source", msg, p.source)
	if err != nil {
		return nil, err
	}
	res.Source = s

	s, err = parser.ParseText("reason", msg, p.reason)
	if err != nil {
		return nil, err
	}
	res.Reason = s

	var timeError error = nil
	for _, t := range p.times {
		ts, err := parser.ParseDateTime(msg, t.search, t.replace)
		if err == nil {
			res.Date = ts
			break
		} else {
			timeError = err
		}
	}
	if res.Date.Beg == 0 {
		return nil, fmt.Errorf("%v for datetime field", timeError)
	}

	var coordError error = nil
	for _, c := range p.coords {
		a, err := parser.ParseCoords(msg, c.kind, c.search, c.replace)
		if err == nil {
			res.Area = a
			break
		} else {
			coordError = err
		}
	}
	if res.Area == nil {
		return nil, coordError
	}
	return &res, nil
}
