package goon

import (
  "unsafe"
  "errors"
)

/*
#cgo CFLAGS: -Wno-error
#include <stdlib.h>
#include "ext/parser.c"
*/
import "C"

var SyntaxError = errors.New("syntax error!")

var symbols map[string]*Value;

func init() {
  symbols = make(map[string]*Value, 1024);
}

func Parse(line []byte) (*Value, error) {
  cs := C.CString(string(line))
  cs_len := C.int(len(line))
  defer C.free(unsafe.Pointer(cs))

  root := C.gn_parse(cs, cs_len)

  if unsafe.Pointer(root) == nil {
    return nil, SyntaxError
  }

  ret := Eval(root)
  if (ret == nil) {
    return nil, SyntaxError
  }

  return Eval(root), nil
}

func Eval(node *C.gn_ast_node_t) *Value {
  //fmt.Printf("in eval! %d/%d/%d\n", node.node_type, node.value, node.num_children)

  if node.node_type == C.GN_AST_NIL {
    return NIL
  } else if node.node_type == C.GN_AST_NUMBER {
    return &Value{int(node.value), TYPE_INT}
  } else if node.node_type == C.GN_AST_BOOLEAN {
    v := int(node.value);
    if v > 0 {
      return &Value{true, TYPE_BOOL}
    } else {
      return &Value{false, TYPE_BOOL}
    }
  } else if node.node_type == C.GN_AST_SYMBOL {
    symbol := C.GoString(C.gn_get_symbol(node))
    return symbols[symbol]
  } else if node.node_type == C.GN_AST_ASSIGN {
    left := C.gn_child_at(node, 0)
    symbol := C.GoString(C.gn_get_symbol(left))
    value := Eval(C.gn_child_at(node, 1))

    symbols[symbol] = value
    return value
  }

  left := Eval(C.gn_child_at(node, 0))
  right := Eval(C.gn_child_at(node, 1))

  switch node.node_type {
  case C.GN_AST_AND:
    return left.And(right)
  case C.GN_AST_OR:
    return left.Or(right)
  case C.GN_AST_COMPARE:
    return left.Equals(right)
  case C.GN_AST_INVERSE_COMPARE:
    return left.NotEquals(right)
  case C.GN_AST_ADD:
    return left.Add(right)
  case C.GN_AST_SUBTRACT:
    return left.Subtract(right)
  case C.GN_AST_MULTIPLY:
    return left.Multiply(right)
  case C.GN_AST_DIVIDE:
    return left.Divide(right)
  }

  return nil
}
