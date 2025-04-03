package numstore

import (

	"github.com/lincketheo/numstore/internal/utils"
)

type parse_tok int

type looseRange struct {
	start        int64
	startPresent bool
	stop         int64
	stopPresent  bool
	step         int32
	stepPresent  bool
}

type Range struct {
	start uint64
	stop  uint64
	step  int32
}

const (
	TOK_LEFT_BRACKET parse_tok = iota
	TOK_RIGHT_BRACKET
	TOK_COLON
	TOK_NUMBER
)

func parseRange(str string) ([]parse_tok, []int, error) {
	i := 0
	end := 0
	ret_toks := make([]parse_tok, 0, 10)
	ret_stack := make([]int, 0, 10)

	for i < len(str) {
		i = end
		if str[i] == ':' {
			ret_toks = append(ret_toks, TOK_COLON)
			end += 1
		} else if str[i] == '[' {
			ret_toks = append(ret_toks, TOK_LEFT_BRACKET)
			end += 1
		} else if str[i] == ']' {
			ret_toks = append(ret_toks, TOK_RIGHT_BRACKET)
			end += 1
		} else if utils.IsDigit(str[i]) {
			for utils.IsDigit(str[end]) {
				end += 1
			}
			// TODO
			//ret_stack = append(ret_stack, int(str[i:end]))
		}
	}

  return ret_toks, ret_stack, nil
}

func interpretRange(toks []parse_tok, stack []int) (looseRange, error) {
  return looseRange{}, nil
}

func looseRangePromote(s looseRange, arrLen uint64) Range {
	var step int32
	if !s.stepPresent {
		step = 1
	}

	ret := Range{
		start: srangeToRange(s.start, s.startPresent, 0, arrLen),
		stop:  srangeToRange(s.stop, s.stopPresent, arrLen, arrLen),
		step:  step,
	}

	return ret
}

func srangeToRange(s int64, isSPresent bool, dflt uint64, arrLen uint64) uint64 {
	if !isSPresent {
		return dflt
	}

	if s < 0 {
		_s := int64(arrLen) + s

		if _s < 0 {
			return 0
		}

		return uint64(_s)
	}

	return uint64(s)
}
