package ast

import (
	"log"
	"strconv"
)


func ParseFile(path string) Ast {
	return parse(lexFile(path))	
}

func ParseBytes(b []byte) Ast {
	return parse(lex(b))
}

func parse(ts tokens) Ast {
	return parseDefs(ts)
}

func parseDefs(ts tokens) Ast {
	if ts.peek().typeof == eof {
		return Nil{}
	}
	ts.assert(def)
	iden := ts.assert(identifier)
	ts.assert(bar)
	expr := parseBinary(ts, 0, breakWith(eof, def))
	return Def{iden.literal, expr, parseDefs(ts)}
}


func parseBinary(ts tokens, min_bp uint, stop breakset) Ast {
	var lhs Ast
	lhs = parsePrefix(ts, stop)

	for {
		tk := ts.peek()

		if stop[tk.typeof] {
			break
		}
 		op_type, lbp, rbp := infixBindPow(tk)

		if lbp < min_bp {
			break
		}
		// only bump past explicit operators, not sequenced application
		if op_type != assoc {
			ts.take()
		}
		rhs := parseBinary(ts, rbp, stop)
		if op_type == pip || op_type == dot {
			lhs = Apply{rhs, lhs}
		} else {
			lhs = Apply{lhs, rhs}
		}
	}
	return lhs
}

func parsePrefix(ts tokens, stop breakset) Ast {
	var ast Ast
	switch ts.peek().typeof {
	case identifier:
		ast = Iden{ts.take().literal}
	case integer:
		ast = parseInt(ts)
	case lpar:
		ts.take()
		ast = parseBinary(ts, 0, breakWith(rpar))
		ts.take()
	case fn:
		ts.take()
		ast = parseFn(ts, stop)
	case cond:
		ts.take()
		ast = parseCond(ts, stop)

	default:
		msg := "expected prefix expression, got %v, line %v"
		log.Fatalf(msg, ts.peek(), ts.peek().line_no)
	}
	return ast
}

func parseInt(ts tokens) Ast {
	tk := ts.take()
	num, err := strconv.ParseInt(tk.literal, 10, 64)
	if err != nil {
		msg := "cannot recognize integer %v, line %v"
		log.Fatalf(msg, tk.literal, tk.line_no)
	}
	return Int64{num}
}

func parseFn(ts tokens, stop breakset) Ast {
	idn := ts.assert(identifier).literal
	if ts.peek().typeof == bar {
		ts.take()
		return Lambda{idn, parseBinary(ts, 0, stop)}
	}
	return Lambda{idn, parseFn(ts, stop)}
}

func parseCond(ts tokens, stop breakset) Ast {
	tests := []Ast{}
	consqs := []Ast{}
	endtest := breakWith(bar)
	endconsq := breakWith(bar, elsetok)
	for {
		tests = append(tests, parseBinary(ts, 0, endtest))
		ts.assert(bar)
		consqs = append(consqs, parseBinary(ts, 0, endconsq))
		nxt := ts.peek().typeof
		if nxt == elsetok {
			ts.take()
			break
		}
		ts.assert(bar)
	}
	alt := parseBinary(ts, 0, stop)
	return Cond{tests, consqs, alt}
}


type breakset map[tokType]bool

func breakWith(ts ...tokType) breakset {
	bs := breakset{}
	for _, t := range ts {
		bs[t] = true
	}
	return bs
}

