package goon

import (
  "strconv"
  "errors"
  "fmt"
)

var SyntaxError = errors.New("Syntax error!")

func UnexpectedError(l *Lexeme, expected string) error {
  return errors.New(fmt.Sprintf("Unexpected %s, expected %s", l, expected))
}

type Parser struct {
  lexer *Lexer
  stack []ASTNode
  lexemes []*Lexeme
}

func (p *Parser) expand() {
  l := <-p.lexer.stream
  p.lexemes = append(p.lexemes, &l)
}

func (p *Parser) shift() *Lexeme {
  if len(p.lexemes) == 0 {
    p.expand()
  }

  l := p.lexemes[0]
  p.lexemes = p.lexemes[1:]
  return l
}

func (p *Parser) accept(t LexemeType) *Lexeme {
  if len(p.lexemes) == 0 {
    p.expand()
  }

  if p.lexemes[0].lexeme_type == t {
    return p.shift()
  }

  return nil
}

func (p *Parser) peek(i int) LexemeType {
  for {
    if len(p.lexemes) > i {
      return p.lexemes[i].lexeme_type
    }

    p.expand()
  }

  return ErrLexeme
}

func (p *Parser) acceptOneOf(ts ...LexemeType) *Lexeme {
  if len(p.lexemes) == 0 {
    p.expand()
  }

  for _, t := range ts {
    if p.lexemes[0].lexeme_type == t {
      return p.shift()
    }
  }

  return nil
}

func (p *Parser) pushNode(node ASTNode) {
  p.stack = append(p.stack, node)
}

func (p *Parser) popNode() ASTNode {
  end := len(p.stack)-1
  n := p.stack[end]
  p.stack = p.stack[:end]

  return n
}

func (p *Parser) pushExpression(op Operator) {
  node := &ExpressionNode{op, nil, nil}
  node.right = p.popNode()
  node.left = p.popNode()
  p.pushNode(node)
}

func (p *Parser) pushValue(v *Value) {
  p.pushNode(&ValueNode{v})
}

/*
block = ((statement EOL) | EOL)+ EOF
*/

func block(p *Parser) error {
  for {
    if p.peek(0) == EOFLexeme {
      break
    }

    empty_line := p.accept(EOLLexeme)
    if empty_line != nil {
       continue
    }

    err := statement(p)
    if err != nil {
      return err
    }

    eol := p.accept(EOLLexeme)
    if eol == nil {
      return UnexpectedError(p.lexemes[0], "EOL")
    }
  }

  node := &BlockNode{make([]ASTNode, len(p.stack))}
  _ = copy(node.children, p.stack)

  p.stack = p.stack[0:0]
  p.pushNode(node)

  return nil
}

/*
statement = ID ASSIGN expression
          / expression
*/
func statement(p *Parser) error {
  if p.peek(0) == IdentLexeme && p.peek(1) == AssignLexeme {
    ident, _ := p.shift(), p.shift()

    err := expression(p)
    if err != nil {
      return err
    }

    expr := p.popNode()
    node := &AssignNode{ident.value, expr}
    p.pushNode(node)

    return nil
  }

  return expression(p)
}

/*
expression = equality ((AND | OR) equality)*
*/

func expression(p *Parser) error {
  err := equality(p)
  if err != nil {
    return err
  }

  for {
    l := p.acceptOneOf(AndLexeme, OrLexeme)
    if l == nil {
      break
    }

    err := equality(p)
    if err != nil {
      return err
    }

    if l.lexeme_type == AndLexeme {
      p.pushExpression(AndOp)
    } else if l.lexeme_type == OrLexeme {
      p.pushExpression(OrOp)
    }
  }

  return nil
}

/*
equality = sum ((EQUALS | DOES_NOT_EQUAL) sum)*
*/
func equality(p *Parser) error {
  err := sum(p)
  if err != nil {
    return err
  }

  for {
    l := p.acceptOneOf(CompareLexeme, InvCompareLexeme)
    if l == nil {
      break
    }

    err := sum(p)
    if err != nil {
      return err
    }

    if l.lexeme_type == CompareLexeme {
      p.pushExpression(CompareOp)
    } else if l.lexeme_type == InvCompareLexeme {
      p.pushExpression(InvCompareOp)
    } else {
      return UnexpectedError(p.lexemes[0], "EOL")
    }
  }

  return nil
}

/*
sum = product ((PLUS | MINUS) product)*
*/
func sum(p *Parser) error {
  err := product(p)
  if err != nil {
    return err
  }

  for {
    l := p.acceptOneOf(AddLexeme, SubtractLexeme)
    if l == nil {
      break
    }

    err := product(p)
    if err != nil {
      return err
    }

    if l.lexeme_type == AddLexeme {
      p.pushExpression(AddOp)
    } else if l.lexeme_type == SubtractLexeme {
      p.pushExpression(SubtractOp)
    }
  }

  return nil
}

/*
product = value ((TIMES | DIVIDED_BY) value)*
*/
func product(p *Parser) error {
  err := value(p)
  if err != nil {
    return err
  }

  for {
    l := p.acceptOneOf(MultiplyLexeme, DivideLexeme)
    if l == nil {
      break
    }

    err := sum(p)
    if err != nil {
      return err
    }

    if l.lexeme_type == MultiplyLexeme {
      p.pushExpression(MultiplyOp)
    } else if l.lexeme_type == DivideLexeme {
      p.pushExpression(DivideOp)
    }
  }

  return nil
}

/*
value = NIL
      / TRUE
      / FALSE
      / NUMBER
      / ID
      / OPEN expression CLOSE
*/
func value(p *Parser) error {
  l := p.acceptOneOf(NilLexeme, TrueLexeme, FalseLexeme, NumberLexeme,
                     IdentLexeme, LeftParenLexeme)

  if l == nil {
    return UnexpectedError(p.lexemes[0], "value")
  }

  switch l.lexeme_type {
  case NilLexeme:
    p.pushValue(NIL)
  case TrueLexeme:
    p.pushValue(TRUE)
  case FalseLexeme:
    p.pushValue(FALSE)
  case NumberLexeme:
    i, _ := strconv.Atoi(l.value)
    p.pushValue(&Value{i, IntType})
  case IdentLexeme:
    p.pushNode(&IdentNode{l.value})
  case LeftParenLexeme:
    err := expression(p)
    if err != nil {
      return err
    }

    close_paren := p.accept(RightParenLexeme)
    if close_paren == nil {
      return UnexpectedError(p.lexemes[0], "')'")
    }
  }

  return nil
}

func Parse(input string) (ASTNode, error) {
  lexer := Lex(input)
  parser := &Parser{lexer, make([]ASTNode, 0, 1024), make([]*Lexeme, 0)}

  err := block(parser)
  if err != nil {
    return nil, err
  }

  if len(parser.stack) != 1 {
    return nil, SyntaxError
  } else if parser.peek(0) != EOFLexeme {
    return nil, UnexpectedError(parser.lexemes[0], "EOF")
  }

  root := parser.popNode()
  return root, nil
}
