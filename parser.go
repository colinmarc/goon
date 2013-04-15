package goon

import (
  "strconv"
  "errors"
)

var SyntaxError = errors.New("Syntax error!")

type Parser struct {
  lexer *Lexer
  stack []ASTNode
  next *Lexeme
}

func (p *Parser) getNext() {
  l := <-p.lexer.stream
  p.next = &l
}

func (p *Parser) accept(t LexemeType) *Lexeme {
  if p.next.lexeme_type == t {
    l := p.next
    p.getNext()

    return l
  }

  return nil
}

func (p *Parser) acceptOneOf(ts ...LexemeType) *Lexeme {
  for _, t := range ts {
    l := p.accept(t)
    if l != nil {
      return l
    }
  }

  return nil
}

func (p *Parser) push(node ASTNode) {
  p.stack = append(p.stack, node)
}

func (p *Parser) pop() ASTNode {
  end := len(p.stack)-1
  n := p.stack[end]
  p.stack = p.stack[:end]

  return n
}

func (p *Parser) empty() {
  p.stack = p.stack[0:0]
}

func (p *Parser) pushExpression(op Operator) {
  node := &ExpressionNode{op, nil, nil}
  node.right = p.pop()
  node.left = p.pop()
  p.push(node)
}

func (p *Parser) pushValue(v *Value) {
  p.push(&ValueNode{v})
}

/*
block = ((statement EOL) | EOL)+ EOF
*/

func block(p *Parser) error {
  for {
    if p.next.lexeme_type == EOFLexeme {
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
       return SyntaxError
    }
  }

  node := &BlockNode{make([]ASTNode, len(p.stack))}
  _ = copy(node.children, p.stack)

  p.empty()
  p.push(node)

  return nil
}

/*
statement = ID ASSIGN expression
          / expression
*/

func statement(p *Parser) error {
  // TODO
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
    } else {
      return SyntaxError
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
      return SyntaxError
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
    } else {
      return SyntaxError
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
    } else {
      return SyntaxError
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
    return SyntaxError
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
    p.push(&IdentNode{l.value})
  case LeftParenLexeme:
    err := expression(p)
    if err != nil {
      return err
    }

    close_paren := p.accept(RightParenLexeme)
    if close_paren == nil {
      return SyntaxError
    }
  }

  return nil
}

func Parse(input string) (ASTNode, error) {
  lexer := Lex(input)

  parser := &Parser{lexer, make([]ASTNode, 0, 1024), nil}
  parser.getNext()


  err := block(parser)
  if err != nil {
    return nil, err
  }

  if len(parser.stack) != 1 || parser.next.lexeme_type != EOFLexeme {
    return nil, SyntaxError
  }

  root := parser.stack[0]
  parser.empty()
  return root, nil
}
