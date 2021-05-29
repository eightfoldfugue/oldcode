package tok

import "fmt"

type Kind uint

const (
	Eof Kind = iota

	// simple data
	Boolean
	Integer
	Float
	Symbol
	Tag
	Identifier

	// delimiters
	LPar
	RPar
	Bar
	Eq

	// operators
	DotPipe
	OnePipe
	TwoPipe
	TriPipe
	OneApp
	TwoApp
	TriApp

	// keywords
	Def
	Obj
	Fun

	Fn
	Let
	In
	Cond
	Else
	Where
)


type Token struct {
	Type Kind
	Literal string
	Line_no int
	Value interface{}
}

const (
	Seq uint = iota
	Dot 
	App
	Pip
)

// try so that all > applies assoc to the right
func InfixBindPow(t Token) (uint, uint, uint) {
	switch t.Type {
	case TriPipe: return Pip, 3, 4
	case TriApp:  return App, 2, 1
	
	case TwoPipe: return Pip, 7, 8 
	case TwoApp:  return App, 6, 5

	case OnePipe: return Pip, 11, 12
	case OneApp:  return App, 10, 9
	
	case DotPipe: return Dot, 15, 16
	default:      return Seq, 13, 14 
	}
}

var ops = map[string]Kind {
	"<": OnePipe,
	"<<": TwoPipe,
	"<<<": TriPipe,
	">": OneApp,
	">>": TwoApp,
	">>>": TriApp,
}

var keyWords = map[string]Kind {
	"def": Def,
	"obj": Obj,
	"fun": Fun,
	"fn" : Fn,
	"let": Let,
	"in" : In,
	"cond": Cond,
	"else": Else,
	"where": Where,
}


func (t Token) String() string {
	s := "token line %v, %v\n"
	return fmt.Sprintf(s, t.Line_no, t.Literal)
}

/*
instead of bars and arrows for seperating constructs, just bars
prefer making application expressive over making lists expressive

fun add-rat r1 r2 | add r1.num r2.num <rat> add r1.den r2.den

fun reduce expr 
|	let is? = eq? expr.tag-of
	in
	cond
		is? :num | expr
	|	is? :sym | expr
	|	is? :add | reduce-add expr.e1 expr.e1
	|   is? :mul | reduce-mul expr.e1 expr.e2
	else error ""

let unfortunate-commas = [1, 2, 3, 4, 5]

obj item w | weight = w

obj pair l r 
| left = l 
| right = r 
| weight = l.weight <add> r.weight


def main mobile.weight
where mobile = item 4 <pair> item 7 <pair> item 8 <pair> item 6

stream x | x + 1 then x*2 ; then aka fby


now all that's left is to do some more tests,
add binding powers, and this experiment is ready
to be moved into the lang project

on the longer horizon, once i have the language with
objects complete, i should not immeadiately go into
writing up lists or modules into the language.
instead i should take time to write about what I've learned

i could move this right into lang... i don't
see any problem with that
*/
