package goon

import (
  "strings"
  "unicode"
)

const (
  digits string = "0123456789"
  eof rune = -1
)

func in(r rune, s string) bool {
  return (strings.IndexRune(s, r) >= 0)
}

type LexemeType int
type LexFn func(*Lexer) LexFn

const (
  ErrLexeme LexemeType = iota
  NilLexeme
  TrueLexeme
  FalseLexeme
  NumberLexeme
  IdentLexeme
  AssignLexeme
  AndLexeme
  OrLexeme
  AddLexeme
  SubtractLexeme
  MultiplyLexeme
  DivideLexeme
  CompareLexeme
  InvCompareLexeme
  OperatorLexeme
  LeftParenLexeme
  RightParenLexeme
  SpaceLexeme
  EOLLexeme
  EOFLexeme
)

var symbols = map[rune]LexemeType{
  '+': AddLexeme,
  '-': SubtractLexeme,
  '*': MultiplyLexeme,
  '/': DivideLexeme,
  '(': LeftParenLexeme,
  ')': RightParenLexeme,
}

var keywords = map[string]LexemeType{
  "nil":   NilLexeme,
  "true":  TrueLexeme,
  "false": FalseLexeme,
  "and":   AndLexeme,
  "or":    OrLexeme,
}

type Lexeme struct {
  lexeme_type LexemeType
  value string
}

func (l *Lexeme) string() string {
  return l.value
}

type Lexer struct {
  input []rune
  window struct {
    start int
    end int
  }
  stream chan Lexeme
}

func (l *Lexer) peek() rune {
  if l.window.end >= len(l.input) {
    return eof
  }
  return l.input[l.window.end]
}

func (l *Lexer) expand() {
  l.window.end += 1
}

func (l *Lexer) skip() {
  l.expand()
  l.discard()
}

func (l *Lexer) discard() {
  l.window.start = l.window.end
}

func (l *Lexer) current() []rune {
  return l.input[l.window.start:l.window.end]
}

func (l *Lexer) emit(lexeme_type LexemeType) {
  chunk := l.current()
  l.stream <- Lexeme{lexeme_type, string(chunk)}
  l.discard()
}

func lexStart(l *Lexer) LexFn {
  // TODO: be smarter about whitespace

  for {
    r := l.peek()

    if r == eof {
      l.emit(EOFLexeme)
      return nil
    } else if r == '\n' {
      l.expand()
      l.emit(EOLLexeme)
    } else if r == ' ' {
      l.skip()
    } else if in(r, digits) {
      return lexNumber
    } else if t, present := symbols[r]; present {
      l.expand()
      l.emit(t)
    } else if unicode.IsLetter(r) || r == '_' {
      return lexWord
    } else {
      l.expand()
      l.emit(ErrLexeme)
    }
  }

  return lexStart
}

func lexNumber(l *Lexer) LexFn {
  for {
    r := l.peek()
    if strings.IndexRune(digits, r) >= 0 {
      l.expand()
    } else {
      l.emit(NumberLexeme)
      break
    }
  }

  return lexStart
}

func lexCompare(l *Lexer) LexFn {
  var r rune
  r = l.peek()
  var t LexemeType

  if r == '=' {
    t = CompareLexeme
  } else {
    t = InvCompareLexeme
  }

  l.expand()
  r = l.peek()
  if r == '=' {
    l.emit(t)
  } else {
    l.emit(ErrLexeme)
  }

  return lexStart
}

func lexWord(l *Lexer) LexFn {
  first := true

  for {
    r := l.peek()

    if unicode.IsLetter(r) || r == '_' || (unicode.IsNumber(r) && !first)  {
      first = false
      l.expand()
    } else {
      current := string(l.current())
      if t, present := keywords[current]; present {
        l.emit(t)
      } else {
        l.emit(IdentLexeme)
      }

      break
    }
  }

  return lexStart
}

func (l *Lexer) Run() {
  for state := lexStart; state != nil; {
    state = state(l)
  }
  close(l.stream)
}

func Lex(input string) *Lexer {
  l := &Lexer {input: []rune(input), stream: make(chan Lexeme)}

  go l.Run()
  return l
}
