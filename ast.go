package goon

import (
  "fmt"
  "strings"
)

type ASTNode interface {
  Evaluate(runtime *Runtime) *Value
  Describe(indent int)
}

// VALUE

type ValueNode struct {
  value *Value
}

func (n *ValueNode) Evaluate(runtime *Runtime) *Value {
  return n.value
}

func (n *ValueNode) Describe(indent int) {
  fmt.Printf("# %sVALUE: %s\n", strings.Repeat("  ", indent), n.value)
}

// IDENT

type IdentNode struct {
  ident string
}

func (n *IdentNode) Evaluate(runtime *Runtime) *Value {
  return runtime.ns[n.ident];
}

func (n *IdentNode) Describe(indent int) {
  fmt.Printf("# %sIDENT: `%s`\n", strings.Repeat("  ", indent), n.ident)
}

// EXPRESSION

type Operator string
const (
  AndOp Operator        = "and"
  OrOp Operator         = "or"
  CompareOp Operator    = "=="
  InvCompareOp Operator = "!="
  AddOp Operator        = "+"
  SubtractOp Operator   = "-"
  MultiplyOp Operator   = "*"
  DivideOp Operator     = "/"
  PowerOp Operator      = "^"
)

type ExpressionNode struct {
  operator Operator
  left ASTNode
  right ASTNode
}

func (n *ExpressionNode) Evaluate(runtime *Runtime) *Value {
  left := n.left.Evaluate(runtime)
  right := n.right.Evaluate(runtime)

  switch n.operator {
  case AndOp:
    return left.And(right)
  case OrOp:
    return left.Or(right)
  case CompareOp:
    return left.Equals(right)
  case InvCompareOp:
    return left.NotEquals(right)
  case AddOp:
    return left.Add(right)
  case SubtractOp:
    return left.Subtract(right)
  case MultiplyOp:
    return left.Multiply(right)
  case DivideOp:
    return left.Divide(right)
  }

  return NIL
}

func (n *ExpressionNode) Describe(indent int) {
  var ex string;
  switch n.operator {
  case AndOp:
    ex = "AND"
  case OrOp:
    ex = "OR"
  case CompareOp:
    ex = "EQUALS"
  case InvCompareOp:
    ex = "NOT_EQUALS"
  case AddOp:
    ex = "ADD"
  case SubtractOp:
    ex = "SUBTRACT"
  case MultiplyOp:
    ex = "MULTIPLY"
  case DivideOp:
    ex = "DIVIDE"
  }

  fmt.Printf("# %s%s:\n", strings.Repeat("  ", indent), ex)
  n.left.Describe(indent + 1)
  n.right.Describe(indent + 1)
}

// ASSIGN

type AssignNode struct {
  ident string
  expr ASTNode
}

func (n *AssignNode) Evaluate(runtime *Runtime) *Value {
  value := n.expr.Evaluate(runtime)

  runtime.ns[n.ident] = value
  return value
}

func (n *AssignNode) Describe(indent int) {
  fmt.Printf("# %sASSIGN `%s` to:\n", strings.Repeat("  ", indent), n.ident)
  n.expr.Describe(indent+1)
}

// BLOCK

type BlockNode struct {
  children []ASTNode
}

func (n *BlockNode) Evaluate(runtime *Runtime) *Value {
  var last *Value
  for _, n := range n.children {
    last = n.Evaluate(runtime)
  }

  return last
}

func (n *BlockNode) Describe(indent int) {
  fmt.Printf("# %sBLOCK (%d):\n", strings.Repeat("  ", indent), len(n.children))
  for _, child := range n.children {
    child.Describe(indent+1)
  }
}

// BRANCH

type CondNode struct {
  cond ASTNode
  then ASTNode
}

type BranchNode struct {
  branches []CondNode
  default_branch ASTNode
}

func (n *BranchNode) AddCond(cond ASTNode, then ASTNode) {
  n.branches = append(n.branches, CondNode{cond, then})
}

func (n *BranchNode) Evaluate(runtime *Runtime) *Value {
  for _, branch := range n.branches {
    v := branch.cond.Evaluate(runtime)
    if v.IsTruthy() {
      return branch.then.Evaluate(runtime)
    }
  }

  if n.default_branch != nil {
    return n.default_branch.Evaluate(runtime)
  }

  return NIL
}

func (n *BranchNode) Describe(indent int) {
  first := true

  for _, branch := range n.branches {
    if first {
      fmt.Printf("# %sIF:\n", strings.Repeat("  ", indent))
      first = false
    } else {
      fmt.Printf("# %sELSE IF:\n", strings.Repeat("  ", indent))
    }

    branch.cond.Describe(indent+1)

    fmt.Printf("# %sTHEN:\n", strings.Repeat("  ", indent))
    branch.then.Describe(indent+1)
  }

  if n.default_branch != nil {
    fmt.Printf("# %sELSE THEN:\n", strings.Repeat("  ", indent))
    n.default_branch.Describe(indent+1)
  }
}

