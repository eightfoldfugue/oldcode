package ast

import (
	"io/ioutil"
	"unicode"
	"unicode/utf8"
	"log"
)

func lexFile(path string) tokens {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	if !utf8.Valid(data) {
		log.Fatal(path, " is not valid utf8")
	}
	return lex(data)
}

func lex(b []byte) tokens {
	s := &subject{input:b, line_no:1}
	s.nextToken()
	return s
}

func (s *subject) nextToken() {
	if !s.more() {
		s.commit(eof)
		return
	}
	switch {
	// junk
	case one(isNewline)(s):
		s.line_no++
		s.junk()
	case oneOrMore(unicode.IsSpace)(s):
		s.junk()
	case isComments(s):
		s.junk()

	// delimiters	
	case match("(")(s):
		s.commit(lpar)
	case match(")")(s):
		s.commit(rpar)
	case match("|")(s):
		s.commit(bar)

	// immeadiate data
	case isFloat(s):
		s.commit(float)
	case isInt(s):
		s.commit(integer)
	case isTag(s):
		s.commit(tag)
	case match("true")(s), match("false")(s):
		s.commit(boolean)

	case isOperator(s):
		optype, isop := operators[s.inspect()]
		s.checkBadOp(isop)
		s.commit(optype)
		
	case isIden(s):
		keytype, isKey := keywords[s.inspect()]
		if isKey {
			s.commit(keytype)
		} else {
			s.commit(identifier)
		}

	default:
		msg := "lexer chocked on %v line %v"
		log.Fatalf(msg, s.failed(), s.line_no)
	}
}

func (s *subject) checkBadOp(ok bool) {
	if !ok {
		msg := "unknown operator %v line %v"
		log.Fatalf(msg, s.inspect(), s.line_no)
	}
}

type lexr func(s *subject) bool

var isComments lexr = seq(one(isComment), zeroOrMore(not(isNewline)))

var isFloat lexr = seq(
	zeroOrOne(isDash),
	oneOrMore(unicode.IsDigit),
	match("."),
	oneOrMore(unicode.IsDigit),
	peekNot(isAlphaSymbolic))

var isInt lexr = seq(
	zeroOrOne(isDash),
	oneOrMore(unicode.IsDigit),
	peekNot(isAlphaSymbolic))

var isTag lexr = seq(
	one(isTagChar),
	one(isAlphaSymbolic),
	zeroOrMore(isIdenChar))

var isOperator lexr = oneOrMore(isOp)

var isIden lexr = seq(
	one(isAlphaSymbolic),
	zeroOrMore(isIdenChar))


type pred func(r rune) bool

func isNewline(r rune) bool {
	return string(r) == "\n"
}

func isDash(r rune) bool {
	return string(r) == "-"	
}

func isDelim(r rune) bool {
	s := string(r)
	return s == "(" || s == "|" || s == ")"
}

func isOp(r rune) bool {
	s := string(r)
	return s == "<" || s == "." || s == ">"
}

func isComment(r rune) bool {
	return string(r) == ";"
}

func isAlphaSymbolic(r rune) bool {
	p1 := !isDelim(r) && !isOp(r) && !isComment(r)
	p2 := !unicode.IsSpace(r)
	p3 := !unicode.IsDigit(r)
	return unicode.IsGraphic(r) && p1 && p2 && p3
}

func isAlphaNumeric(r rune) bool {
	p1 := !isDelim(r) && !isOp(r) && !isComment(r)
	p2 := !unicode.IsSpace(r)
	return unicode.IsGraphic(r) && p1 && p2
}

func isTagChar(r rune) bool {
	return string(r) == ":"
}

func isIdenChar(r rune) bool {
	return isAlphaNumeric(r) && !isTagChar(r)
}


// end of defining forms, machinery follows from here on out
func not(p pred) pred {
	f := func(r rune) bool {
		return !p(r)
	}
	return f
}


func zeroOrMore(p pred) lexr {
	f := func(s *subject) bool {
		for one(p)(s) {}
		return true
	}
	return f
}

func oneOrMore(p pred) lexr {
	f := func(s *subject) bool {
		if one(p)(s) {
			for one(p)(s){}
			return true
		}
		return false
	}
	return f
}

func zeroOrOne(p pred) lexr {
	f := func(s *subject) bool {
		one(p)(s)
		return true
	}
	return f
}

func one(p pred) lexr {
	f := func(s *subject) bool {
		if s.more() && p(s.peekChar()) {
			s.step()
			return true
		}
		return false
	}
	return f
}

func peekNot(p pred) lexr {
	f := func(s *subject) bool {
		return !p(s.peekChar())
	}
	return f
}

func seq(lexs ...lexr) lexr {
	f := func(s *subject) bool {
		pass := true
		for _, lex := range lexs {
			pass = pass && lex(s)
			if !pass {
				s.reject()
				break
			}
		}
		return pass	
	}
	return f
}

func match(str string) lexr {
	f := func(s *subject) bool {
		pass := true
		for i, w := 0, 0; i < len(str); i+=w {
			rstr, wstr := utf8.DecodeRuneInString(str[i:])
			rsub := s.step()
			pass = pass && (rstr == rsub)
			if !pass {
				s.reject()
				break
			}
			w = wstr
		}
		return pass
	}
	return f
}


type subject struct {
	current token
	input   []byte
	begin   int
	end     int
	line_no int
}

func (s *subject) peek() token {
	return s.current
}
func (s *subject) take() token {
	c := s.current
	s.nextToken()
	return c
}
func (s *subject) assert(t tokType) token {
	tok := s.take()
	if tok.typeof != t {
		msg := "parse failure, wanted %v got %v"
		log.Fatalf(msg, t, s.current)
	}
	return tok
}

func (s *subject) decode() (rune, int) {
	return utf8.DecodeRune(s.input[s.end:])
}

func (s *subject) more() bool {
	r, _ := s.decode()
	if r == utf8.RuneError {
		return false
	}
	return true
}

func (s *subject) peekChar() rune {
	r, _ := s.decode()
	return r
}

func (s *subject) step() rune {
	r, w := s.decode()
	s.end += w
	return r
}

func (s *subject) junk() {
	s.begin = s.end
	s.nextToken()
}
func (s *subject) reject() {
	s.end = s.begin
}

func (s *subject) commit(t tokType) {
	tok := token{typeof: t, literal: s.inspect(), line_no: s.line_no}
	s.current = tok
	s.begin = s.end
}

func (s *subject) inspect() string {
	return string(s.input[s.begin:s.end])	
}

func (s *subject) failed() string {
	return string(s.input[s.begin:s.begin+12]) + "..."
}
