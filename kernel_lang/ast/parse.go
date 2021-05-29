package ast

import (
	"lang/tok"
	"log"
)


func Parse(tks []tok.Token) Defs {
	i := &input{in: tks}
	
	
	return parseDefs(i)
}

type breakset map [tok.Kind] bool

func parseDefs(i *input) Defs {
	d := Defs{}
	for i.peek().Type != tok.Eof {
		i.assert(tok.Def)
		name := i.assert(tok.Identifier).Literal
		sentinals := breakset{tok.Eof: true, tok.Def: true,}
		expr := expr_bp(i, 0, sentinals)

		_, alreadyDefined := d[name]
		if alreadyDefined {
		log.Fatalf("%v is a duplicate definition", name)
	}
		d[name] = expr
	}
	i.assert(tok.Eof)
	return d
}

func expr_bp(i *input, min_bp uint, is_sentinal breakset) Tree {
	var lhs Tree

	lhs = parsePrefix(i, is_sentinal)

	for {
	
		// check for end
		t := i.peek()
		if is_sentinal[t.Type] {
			break
		}
			
		op_type, lbp, rbp := tok.InfixBindPow(t)

		if lbp < min_bp {
			break
		}

		// only bump past explicit tokens, not sequenced application
		if op_type != tok.Seq {
			i.take()
		}

		rhs := expr_bp(i, rbp, is_sentinal)

		// if it's a pipe, reverse it, else keep it in order
		
		if op_type == tok.Pip || op_type == tok.Dot {
			lhs = Apply{rhs, lhs}
		} else {
			lhs = Apply{lhs, rhs}
		}
	}

	return lhs
}

func parsePrefix(i *input, sentinals breakset) Tree {
	var t Tree
	switch i.peek().Type {
	
	case tok.Identifier:
		t = Iden{i.take().Literal}

	case tok.Integer:
		t = Int64{i.take().Value.(int64)}

	case tok.LPar:

		i.take() // eat open
		s := breakset{tok.RPar: true}
		t = expr_bp(i, 0, s)
		i.take() // eat close

	case tok.Fn:
		// take fn tok
		i.take()
		t = parseFn(i, sentinals)
		 
	default:
		log.Fatal("wanted atom, got \n", i.take())
	}

	return t
}


func parseFn(i *input, s breakset) Tree {
	idn := i.assert(tok.Identifier).Literal
	if i.peek().Type == tok.Bar {
		i.take()
		return Lambda{idn, expr_bp(i, 0, s)}
	}
	return Lambda{idn, parseFn(i, s)}
}
s
type input struct {
	in    []tok.Token
	index int
}

func (i *input) peek() tok.Token {
	return i.in[i.index]
}

func (i *input) take() tok.Token {
	t := i.peek()
	i.index++
	return t
}

func (i *input) assert(tokType tok.Kind) tok.Token {
	t := i.take()
	if t.Type != tokType {
		log.Fatalf("parse error\nexpected %v got %v", tokType, t)
	}
	return t
}
