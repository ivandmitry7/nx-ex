package task

import (
	"fmt"
	"github.com/o-kos/nx-ex/internal/parser"
	"regexp"
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
	search *regexp.Regexp
}

type reCoords struct {
	kind   string
	search *regexp.Regexp
}

type ReProc struct {
	source *regexp.Regexp
	reason *regexp.Regexp
	times  []reTimes
	coords []reCoords
}

func (p *ReProc) Compile(cfg *ParserCfg) error {
	re, err := regexp.Compile(cfg.Source)
	if err != nil {
		return err
	}
	p.source = re

	re, err = regexp.Compile(cfg.Reason)
	if err != nil {
		return err
	}
	p.reason = re

	for _, ts := range cfg.Times {
		re, err := regexp.Compile(ts.Search)
		if err != nil {
			return err
		}
		p.times = append(p.times, reTimes{re})
	}

	for _, cs := range cfg.Coords {
		re, err := regexp.Compile(cs.Search)
		if err != nil {
			return err
		}
		p.coords = append(p.coords, reCoords{cs.Type, re})
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
		ts, err := parser.ParseDateTime(msg, t.search)
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
		a, err := parser.ParseCoords(msg, c.kind, c.search)
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
