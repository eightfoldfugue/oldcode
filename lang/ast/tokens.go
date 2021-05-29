package ast

import (
	"fmt"
)

// concrete type *subject satisfies the tokens interface
type tokens interface {
	peek() token
	take() token
	assert(tokType) token
}

type token struct {
	typeof  tokType
	literal string
	line_no int
}

func (t token) String() string {
	s := "{%v, %v, %v}"
	return fmt.Sprintf(s, t.typeof.String(), t.literal, t.line_no)
}

type tokType int

const (
	eof tokType = iota

	// immeadiate data
	float
	integer
	tag
	boolean
	
	identifier

	// delimiters
	lpar
	rpar
	bar

	// operators
	dotPipe
	onePipe
	twoPipe
	triPipe
	oneApp
	twoApp
	triApp

	// keywords, top level
	def
	obj
	fun

	// keywords, expression level
	fn
	let
	cond
	elsetok  // else reserved for go
	where
)

var stringer = [...]string{
	"EOF",

	"float",
	"integer",
	"tag",
	"boolean",
	
	"identifier",

	"lpar",
	"rpar",
	"bar",
}

type tokMap map[string]tokType

var operators = tokMap {
	"."  : dotPipe,
	"<"  : onePipe,
	"<<" : twoPipe,
	"<<<": triPipe,
	">"  : oneApp,
	">>" : twoApp,
	">>>": triApp,
}

var keywords = tokMap {
	"def"  : def,
	"obj"  : obj,
	"fun"  : fun,
	"fn"   : fn,
	"="    : let,
	"cond" : cond,
	"else" : elsetok,
	"where": where,
}

func (t tokType) String() string {
	key, isKey := find(t, keywords)
	op, isOp := find(t, operators)
	

	switch {
	case isKey: return key
	case isOp: return op
	default: return stringer[t]
	}
}

func find(t tokType, m tokMap) (string, bool) {
	for k, v := range m {
		if t == v {
			return k, true
		}
	}
	return "", false
}


// information for binding power lookups

const (
	assoc uint = iota
	dot 
	app
	pip

)

// all > associate to the right, everything else to the left
// default indicates application by association

func infixBindPow(t token) (uint, uint, uint) {
	switch t.typeof {

	case triPipe: return pip, 3, 4
	case triApp:  return app, 2, 1
	
	case twoPipe: return pip, 7, 8 
	case twoApp:  return app, 6, 5

	case onePipe: return pip, 11, 12
	case oneApp:  return app, 10, 9
	
	case dotPipe: return dot, 15, 16
	default:      return assoc, 13, 14 
	}
}
