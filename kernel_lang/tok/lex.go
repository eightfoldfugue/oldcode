package tok

import (
	"unicode"
	"unicode/utf8"
	"strings"
	"strconv"
	"log"
	"io/ioutil"
)

func LexFile(path string) []Token {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	if !utf8.Valid(data) {
		log.Fatal(path, " is not valid utf8")
	}
	return Lex(string(data))
}

func Lex(str string) []Token {
	lxr := lexer{subject: str, line_no: 1}
	lxr.lex()
	return lxr.output
}

type lexer struct {
	subject string
	pos int
	line_no int
	output  []Token
}

func (l *lexer) lex() {
	for !l.eof() {
		r := l.peek()
		switch {
		case l.matchChar("\n"):
			l.line_no++		
		case unicode.IsSpace(r):
			l.advance()
		case l.matchChar(";"):
			for string(l.peek()) != "\n" && !l.eof() {
				l.advance()
			}
		case l.matchChar("("):
			l.emit(LPar, "(")
		case l.matchChar(")"):
			l.emit(RPar, ")")
		case l.matchChar("|"):
			l.emit(Bar, "|")
		case l.matchChar("="):
			l.emit(Eq, "=")
		case l.matchChar("."):
			l.emit(DotPipe, ".")
		case isApp(r):
			s := l.lexMany(isApp)
			l.emitOp(s)
		case isPipe(r):
			s := l.lexMany(isPipe)
			l.emitOp(s)
		case isChar(r):
			s := l.lexMany(isChar)
			l.emitWord(s)
		default:
			log.Fatalf("unrecognized token line %v", l.line_no)
		}
	}
	l.emit(Eof, "")
}

func (l *lexer) eof() bool {
	return l.peek() == utf8.RuneError
}

func (l *lexer) peek() rune {
	r, _ := utf8.DecodeRuneInString(l.subject[l.pos:])
	return r
}

func (l *lexer) advance() {
	_, w := utf8.DecodeRuneInString(l.subject[l.pos:])
	l.pos += w
}

func (l *lexer) emit(k Kind, lit string) {
	tok := Token{k, lit, l.line_no, nil}
	l.output = append(l.output, tok)
}

func isChar(r rune) bool {
	reserved := ";()|.<>="
	unreserved := !strings.ContainsRune(reserved, r)
	notspace := !unicode.IsSpace(r) && r != utf8.RuneError
	return unicode.IsGraphic(r) && unreserved && notspace
}

func isApp(r rune) bool {
	return string(r) == ">"
}

func isPipe(r rune) bool {
	return string(r) == "<"	
}

func (l *lexer) lexMany(f func(rune) bool) string {
	mark := l.pos
	for f(l.peek()) {
		l.advance()
	}
	return l.subject[mark:l.pos]
}

func (l *lexer) matchChar(s string) bool {
	if string(l.peek()) == s {
		l.advance()
		return true
	}
	return false
}

func (l *lexer) emitOp(s string) {
	op, ok := ops[s]
	if !ok {
		msg := "%v is unrecogized operator line %v"
		log.Fatalf(msg, s, l.line_no)
	}
	l.emit(op, s)
}

func (l *lexer) emitWord(s string) {
	key, ok := keyWords[s]
	if !ok {
		l.nonKey(s)
	} else {
		l.emit(key, s)
	}
	
}

func (l *lexer) construct(k Kind, lit string, i interface{}) {
 	tok := Token{k, lit, l.line_no, i}
 	l.output = append(l.output, tok)
}

func (l *lexer) nonKey(str string) {
	b_check := str == "true" || str == "false"
	b, _ := strconv.ParseBool(str) // edge case on t, f
	i, int_err := strconv.ParseInt(str, 10, 64)
	f, flo_err := strconv.ParseFloat(str, 64)
	s_check := string(str[0]) == "'"
	t_check := string(str[0]) == ":"

	switch {
	case s_check:
		l.construct(Symbol, str, nil)
	case t_check:
		l.construct(Tag, str, nil)
	case b_check:
		l.construct(Boolean, str, b)
	case int_err == nil:
		l.construct(Integer, str, i)
	case flo_err == nil:
		l.construct(Float, str, f)
	default:
		l.construct(Identifier, str, nil)
	}
}
