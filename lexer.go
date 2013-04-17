package goon

import (
  "strings"
  "unicode"
  "fmt"
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
  _
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
  IfLexeme
  UnlessLexeme
  ElifLexeme
  ElseLexeme
  ThenLexeme
  SpaceLexeme
  EOLLexeme
  IndentLexeme
  EOFLexeme
)

var symbols = map[rune]LexemeType{
  '+': AddLexeme,
  '-': SubtractLexeme,
  '*': MultiplyLexeme,
  '/': DivideLexeme,
  '(': LeftParenLexeme,
  ')': RightParenLexeme,
  ':': ThenLexeme,
}

var keywords = map[string]LexemeType{
  "nil":    NilLexeme,
  "true":   TrueLexeme,
  "false":  FalseLexeme,
  "and":    AndLexeme,
  "or":     OrLexeme,
  "if":     IfLexeme,
  "unless": UnlessLexeme,
  "elif":   ElifLexeme,
  "else":   ElseLexeme,
}

type Lexeme struct {
  lexeme_type LexemeType
  value string
}

func (l *Lexeme) String() string {
  if l.lexeme_type == EOFLexeme {
    return "EOF"
  } else if l.lexeme_type == EOLLexeme {
    return "EOL"
  } else if l.lexeme_type == IndentLexeme {
    return fmt.Sprintf("INDENT(%d)", len(l.value))
  }

  return fmt.Sprintf("`%s` (%d)", l.value, l.lexeme_type)
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
  return lexIndent
}

func lexCode(l *Lexer) LexFn {
  // TODO: be smarter about whitespace

  for {
    r := l.peek()

    if r == eof {
      l.emit(EOFLexeme)
      return nil
    } else if r == '\n' {
      l.expand()
      l.emit(EOLLexeme)
      return lexIndent
    } else if r == ' ' {
      l.skip()
    } else if in(r, digits) {
      return lexNumber
    } else if t, present := symbols[r]; present {
      l.expand()
      l.emit(t)
    } else if r == '=' || r == '!' {
      return lexCompare
    } else if unicode.IsLetter(r) || r == '_' {
      return lexWord
    } else {
      l.expand()
      l.emit(ErrLexeme)
    }
  }

  return lexStart
}

func lexIndent(l *Lexer) LexFn {
  for {
    if l.peek() != ' ' {
      break
    }

    l.expand()
  }

  l.emit(IndentLexeme)

  return lexCode
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

  return lexCode
}

func lexCompare(l *Lexer) LexFn {
  var r rune

  r = l.peek()
  l.expand()
  var t LexemeType
  if r == '=' {
    t = CompareLexeme
  } else if r == '!' {
    t = InvCompareLexeme
  } else {
    l.emit(ErrLexeme)
  }

  r = l.peek()
  if r == '=' {
    l.expand()
    l.emit(t)
  } else {
    l.emit(AssignLexeme)
  }

  return lexCode
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

  return lexCode
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
