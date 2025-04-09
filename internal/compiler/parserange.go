package compiler

import (
	"strconv"

	"github.com/lincketheo/numstore/internal/numstore"
)

type DimRange struct {
	start int
	stop  int
	step  int
	isInf bool
}

var rangeTokStarts = []tokenType{
	TOK_INTEGER,
	TOK_COLON,
	TOK_COMMA,
	TOK_RIGHT_BRACKET,
}

func (p *parser) parseRange() ([]numstore.DimRange, bool) {
	if _, ok := p.expect(TOK_LEFT_BRACKET); !ok {
		return []numstore.DimRange{}, true
	}

	ret := []numstore.DimRange{}

	for {
		if r, ok := p.parseSingleRange(); !ok {
			return nil, false
		} else {
			ret = append(ret, r)
		}

		if _, ok := p.expect(TOK_RIGHT_BRACKET); ok {
			return ret, true
		} else if _, ok := p.expect(TOK_COMMA); ok {
			continue
		} else {
			return nil, false
		}
	}
}

func (p *parser) parseSingleRange() (numstore.DimRange, bool) {
	builder := dimRangeBuilder{}

	for {
		if tok, ok := p.expect(TOK_INTEGER, TOK_COLON); ok {
			if !builder.next(tok) {
				return numstore.DimRange{}, false
			}
		} else {
			return builder.build(), true
		}
	}
}

type dimRangeBuilder struct {
	nums   [3]int
	nset   [3]bool
	i      int
	single bool
}

func (d *dimRangeBuilder) next(t token) bool {
	if d.i >= 3 {
		return false
	}
	switch t.ttype {

	case TOK_COLON:
		// Advance forward - check that no more than 2 colons
		d.single = false
		if d.i >= 2 {
			return false
		}
		d.i++
		return true

	case TOK_INTEGER:
		// Check that nset[i] isn't already set, then update nums and nset
		if d.nset[d.i] {
			return false
		}
		d.nset[d.i] = true
		var err error
		if d.nums[d.i], err = strconv.Atoi(t.value); err != nil {
			return false
		}
		return true

	default:
		return false
	}
}

func (d dimRangeBuilder) build() numstore.DimRange {
	ret := numstore.DimRange{
		Start:    d.nums[0],
		Stop:     d.nums[1],
		Step:     d.nums[2],
		IsInf:    !d.nset[1],
		IsSingle: d.single,
	}
	if !d.nset[0] {
		ret.Start = 0
	}
	if !d.nset[1] {
		ret.IsInf = true
	}
	if !d.nset[2] {
		ret.Step = 1
	}
	if d.i == 0 {
		ret.IsSingle = true
	}
	return ret
}
