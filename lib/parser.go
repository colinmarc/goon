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
  indentation int
}

func (p *Parser) expand() {
  l := <-p.lexer.stream
  if l.lexeme_type == 0 {
    l = Lexeme{EOFLexeme, ""}
  }

  //fmt.Printf("got lexeme: %s\n", l.String())
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
  if p.peek(0) == t {
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

func (p *Parser) popTwoNodes() (ASTNode, ASTNode) {
  second, first := p.popNode(), p.popNode()
  return first, second
}

func (p *Parser) pushExpression(op Operator) {
  left, right := p.popTwoNodes()
  node := &ExpressionNode{op, left, right}
  p.pushNode(node)
}

func (p *Parser) pushValue(v *Value) {
  p.pushNode(&ValueNode{v})
}

/*
block = ((control EOL) | EOL)+ EOF
*/
func block(p *Parser) error {
  mark := len(p.stack)

  for {
    indent := p.accept(IndentLexeme)
    spaces := len(indent.value)
    if (spaces % 2 != 0) || ((spaces / 2) > p.indentation) {
      return errors.New(fmt.Sprintf("Unexpected indent (%d)", spaces))
    } else if (spaces / 2) < p.indentation {
      break
    }

    if p.peek(0) == EOFLexeme {
      break
    }

    empty_line := p.accept(EOLLexeme)
    if empty_line != nil {
       continue
    }

    err := control(p)
    if err != nil {
      return err
    }

    eol := p.acceptOneOf(EOLLexeme, EOFLexeme)
    if eol == nil {
      return UnexpectedError(p.lexemes[0], "EOL/EOF")
    } else if eol.lexeme_type == EOFLexeme {
      break
    }
  }

  node := &BlockNode{make([]ASTNode, len(p.stack) - mark)}
  _ = copy(node.children, p.stack[mark:])

  p.stack = p.stack[0:mark]
  p.pushNode(node)

  return nil
}

/*
control = (IF | UNLESS) expression THEN block
            (ELIF expression THEN block)*
            (ELSE expression THEN block)?
        / definition
*/
func control(p *Parser) error {
  // TODO: unless won't work
  var l *Lexeme
  var err error

  branch := p.acceptOneOf(IfLexeme, UnlessLexeme)
  if branch != nil {
    branch_node := &BranchNode{make([]CondNode, 0), nil}

    err = expression(p)
    if err != nil {
      return err
    }

    l = p.accept(ThenLexeme)
    if l == nil {
      return UnexpectedError(p.lexemes[0], "':'")
    }

    l = p.accept(EOLLexeme)
    if l == nil {
      return UnexpectedError(p.lexemes[0], "EOL")
    }

    p.indentation++

    err = block(p)
    if err != nil {
      return err
    }

    p.indentation--

    branch_node.AddCond(p.popTwoNodes())

    for {
      elif_branch := p.accept(ElifLexeme)
      if elif_branch == nil {
        break
      }

      err = expression(p)
      if err != nil {
        return err
      }

      l = p.accept(ThenLexeme)
      if l == nil {
        return UnexpectedError(p.lexemes[0], "':'")
      }

      l = p.accept(EOLLexeme)
      if l == nil {
        return UnexpectedError(p.lexemes[0], "EOL")
      }

      p.indentation++

      err = block(p)
      if err != nil {
        return err
      }

      p.indentation--

      branch_node.AddCond(p.popTwoNodes())
    }

    else_branch := p.accept(ElseLexeme)
    if else_branch != nil {
      l = p.accept(ThenLexeme)
      if l == nil {
        return UnexpectedError(p.lexemes[0], "':'")
      }

      l = p.accept(EOLLexeme)
      if l == nil {
        return UnexpectedError(p.lexemes[0], "EOL")
      }

      p.indentation++

      err = block(p)
      if err != nil {
        return err
      }

      p.indentation--

      branch_node.default_branch = p.popNode()
    }

    p.pushNode(branch_node)
    return nil
  }

  return definition(p)
}

/*
definition = METHOD_ID (LEFT_P ID (COMMA ID)+ RIGHT_P)? DEF block
           / inline_conditional
*/
func definition(p *Parser) error {
  var l *Lexeme

  method_id := p.accept(MethodIdentLexeme)
  if method_id != nil {
    node := &DefNode{method_id.value, make([]string, 0), nil}

    l = p.accept(LeftParenLexeme)
    if l != nil {
      var ident *Lexeme

      ident = p.accept(IdentLexeme)
      if ident == nil {
        return UnexpectedError(p.lexemes[0], "ID")
      }
      node.AddArgument(ident.value)

      for {
        l = p.acceptOneOf(CommaLexeme, RightParenLexeme)
        if l == nil {
          return UnexpectedError(p.lexemes[0], "')' or ','")
        }

        if l.lexeme_type == RightParenLexeme {
          break
        }

        ident = p.accept(IdentLexeme)
        if ident == nil {
          return UnexpectedError(p.lexemes[0], "ID")
        }
        node.AddArgument(ident.value)
      }
    }

    def := p.accept(DefLexeme)
    if def == nil {
      return UnexpectedError(p.lexemes[0], "'->'")
    }

    eol := p.accept(EOLLexeme)
    if eol == nil {
      return UnexpectedError(p.lexemes[0], "EOL")
    }

    p.indentation++

    err := block(p)
    if err != nil {
      return err
    }

    p.indentation--

    node.block = p.popNode().(*BlockNode)
    p.pushNode(node)

    return nil
  }

  return inline_conditional(p)
}

/*
inline_conditional = statement ((IF | UNLESS) expr)?
*/
func inline_conditional(p *Parser) error {
  err := statement(p)
  if err != nil {
    return err
  }

  cond := p.acceptOneOf(IfLexeme, UnlessLexeme)
  if cond != nil {
    err := expression(p)
    if err != nil {
      return err
    }

    branch_node := &BranchNode{make([]CondNode, 0), nil}
    branch_node.AddCond(p.popNode(), p.popNode())
    p.pushNode(branch_node)
  }

  return nil
}

/*
statement = ID ASSIGN expression
          / KEYWORD expression
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

  kw := p.acceptOneOf(ReturnLexeme, PrintLexeme)
  if kw != nil {
    err := expression(p)
    if err != nil {
      return err
    }

    n := p.popNode()
    var k Keyword
    switch kw.lexeme_type {
    case ReturnLexeme:
      k = ReturnKeyword
    case PrintLexeme:
      k = PrintKeyword
    }

    p.pushNode(&KeywordNode{k, n})
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
value = LEFT_P expression RIGHT_P
      / NIL
      / TRUE
      / FALSE
      / NUMBER
      / id
*/
func value(p *Parser) error {
  l := p.acceptOneOf(NilLexeme, TrueLexeme, FalseLexeme, NumberLexeme, LeftParenLexeme)

  if l == nil {
    return id(p)
  }

  switch l.lexeme_type {
  case LeftParenLexeme:
    err := expression(p)
    if err != nil {
      return err
    }

    close_paren := p.accept(RightParenLexeme)
    if close_paren == nil {
      return UnexpectedError(p.lexemes[0], "')'")
    }
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

  }

  return nil
}

/*
id = ID (LEFT_P (expression (COMMA expression)+)? RIGHT_P)?
*/
func id(p *Parser) error {
  ident := p.accept(IdentLexeme)
  if ident == nil {
    return UnexpectedError(p.lexemes[0], "ID")
  }

  var l *Lexeme

  l = p.accept(LeftParenLexeme)
  if l != nil {

    node := &CallNode{ident.value, make([]ASTNode, 0)}
    first := true
    for {
      l = p.accept(RightParenLexeme)
      if l != nil {
        break
      }

      if first {
        first = false
      } else {
        l = p.accept(CommaLexeme)
        if l == nil && !first {
          return UnexpectedError(p.lexemes[0], "')' or ','")
        }
      }

      err := expression(p)
      if err != nil {
        return err
      }

      node.AddArgument(p.popNode())
    }

    p.pushNode(node)
  } else {
    p.pushNode(&IdentNode{ident.value})
  }

  return nil
}

func Parse(input string) (ASTNode, error) {
  lexer := Lex(input)
  parser := &Parser{lexer, make([]ASTNode, 0, 1024), make([]*Lexeme, 0), 0}

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
