package main

import (
  "fmt"
  "os"
  "io"
  "bufio"
  "unsafe"
)

/*
#cgo CFLAGS: -Wno-error
#include <stdlib.h>
#include "parser.c"
*/
import "C"

var symbols map[string]int;

func main() {
  // TODO: move this out
  symbols = make(map[string]int, 100);

  reader := bufio.NewReader(os.Stdin)

  for {
    fmt.Printf(">> ")
    raw_line, err := reader.ReadBytes('\n')

    if err != nil {
      if err == io.EOF {
        fmt.Printf("quitting...\n")
        break
      } else {
        fmt.Printf("err: %s", err)
      }
      break
    }

    if len(raw_line) > 1 {
      res, _ := Parse(raw_line)
      fmt.Printf("%d\n", res)
    }
  }
}

func Parse(line []byte) (int, error) {
  cs := C.CString(string(line))
  cs_len := C.int(len(line))
  defer C.free(unsafe.Pointer(cs))

  root := C.gn_parse(C.gn_global_context(), cs, cs_len)

  // if(root == C.NULL) {
  //   return 0, nil
  // }

  return Eval(root), nil
}

func Eval(node *C.gn_ast_node_t) int {
  //fmt.Printf("in eval! %d/%d/%d\n", node.node_type, node.value, node.num_children)

  if node.node_type == C.GN_AST_NUMBER {
    return int(node.value);
  } else if node.node_type == C.GN_AST_SYMBOL {
    symbol := C.GoString(C.gn_get_symbol(C.gn_global_context(), node))
    fmt.Printf("symbol is: %s\n", symbol)
    return symbols[symbol]
  }

  left := C.gn_child_at(node, 0)
  right := C.gn_child_at(node, 1)

  switch node.node_type {
  case C.GN_AST_ASSIGN:
    symbol := C.GoString(C.gn_get_symbol(C.gn_global_context(), left))
    value := Eval(right)

    fmt.Printf("symbol is: %s, value is %d\n", symbol, value)
    symbols[symbol] = value
    return value
  case C.GN_AST_ADD:
    return Eval(left) + Eval(right)
  case C.GN_AST_SUBTRACT:
    return Eval(left) - Eval(right)
  case C.GN_AST_MULTIPLY:
    return Eval(left) * Eval(right)
  case C.GN_AST_DIVIDE:
    return Eval(left) / Eval(right)
  }

  return 0
}
