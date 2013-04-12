package goon

import (
  "unsafe"
)

/*
#cgo CFLAGS: -Wno-error
#include <stdlib.h>
#include "ext/parser.c"
*/
import "C"

var symbols map[string]*Value;

func init() {
  symbols = make(map[string]*Value, 1024);
}

func Parse(line []byte) (*Value, error) {
  cs := C.CString(string(line))
  cs_len := C.int(len(line))
  defer C.free(unsafe.Pointer(cs))

  root := C.gn_parse(C.gn_global_context(), cs, cs_len)

  // if(root == C.NULL) {
  //   return 0, nil
  // }

  return Eval(root), nil
}

func Eval(node *C.gn_ast_node_t) *Value {
  //fmt.Printf("in eval! %d/%d/%d\n", node.node_type, node.value, node.num_children)

  if node.node_type == C.GN_AST_NUMBER {
    return &Value{int(node.value), TYPE_INT}
  } else if node.node_type == C.GN_AST_BOOLEAN {
    //todo
    return &Value{int(node.value), TYPE_INT}
  } else if node.node_type == C.GN_AST_SYMBOL {
    symbol := C.GoString(C.gn_get_symbol(C.gn_global_context(), node))
    return symbols[symbol]
  }

  left := C.gn_child_at(node, 0)
  right := C.gn_child_at(node, 1)

  switch node.node_type {
  case C.GN_AST_ASSIGN:
    symbol := C.GoString(C.gn_get_symbol(C.gn_global_context(), left))
    value := Eval(right)

    symbols[symbol] = value
    return value
  case C.GN_AST_ADD:
    return Eval(left).Add(Eval(right))
  case C.GN_AST_SUBTRACT:
    return Eval(left).Subtract(Eval(right))
  case C.GN_AST_MULTIPLY:
    return Eval(left).Multiply(Eval(right))
  case C.GN_AST_DIVIDE:
    return Eval(left).Divide(Eval(right))
  }

  return &Value{nil, TYPE_NIL}
}
